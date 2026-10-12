package main

import (
	"context"
	"testing"
	"time"
)

func emotionScores(label string) map[string]float64 {
	scores := make(map[string]float64, len(emotionLabels))
	for _, emotion := range emotionLabels {
		scores[emotion] = 0
	}
	scores[label] = 1
	return scores
}

func TestOnChatQueuesMessageText(t *testing.T) {
	hub := NewHub()
	hub.OnChat(ChatMessage{User: "Viewer", Text: "¡Qué alegría!"})

	select {
	case got := <-hub.messageQueue:
		if got.text != "¡Qué alegría!" {
			t.Fatalf("queued message = %+v", got)
		}
	default:
		t.Fatal("chat message was not queued for emotion analysis")
	}
}

func TestOnChatIgnoresEmptyMessages(t *testing.T) {
	hub := NewHub()
	hub.OnChat(ChatMessage{Text: "  \t"})
	if got := len(hub.messageQueue); got != 0 {
		t.Fatalf("queued empty message; queue length = %d", got)
	}
}

func TestApplyEmotionResultsAggregatesOnlyCurrentBatch(t *testing.T) {
	hub := NewHub()
	hub.applyEmotionResults([]map[string]float64{
		emotionScores("joy"),
		emotionScores("joy"),
		emotionScores("anger"),
	})
	if got := hub.currentEmotions["joy"]; got < 0.66 || got > 0.67 {
		t.Fatalf("group joy score = %f, want about 0.667", got)
	}
	if got := hub.currentEmotions["anger"]; got < 0.33 || got > 0.34 {
		t.Fatalf("group anger score = %f, want about 0.333", got)
	}

	hub.applyEmotionResults([]map[string]float64{emotionScores("anger")})
	if got := hub.currentEmotions["anger"]; got != 1 {
		t.Fatalf("new batch retained older emotions; anger score = %f, want 1", got)
	}
	if got := hub.currentEmotions["joy"]; got != 0 {
		t.Fatalf("new batch retained older joy score = %f, want 0", got)
	}
}

func TestInvalidEmotionResultsDoNotChangeGroupMood(t *testing.T) {
	hub := NewHub()
	hub.currentEmotions = emotionScores("others")
	hub.applyEmotionResults([]map[string]float64{{"joy": 1}})

	if got := hub.currentEmotions["others"]; got != 1 {
		t.Fatalf("invalid response changed group mood; neutral score = %f", got)
	}
}

func TestOnChatDropsMessagesWhileEmotionModelIsProcessing(t *testing.T) {
	hub := NewHub()
	hub.beginEmotionProcessing()
	hub.OnChat(ChatMessage{Text: "se descarta"})
	if got := len(hub.messageQueue); got != 0 {
		t.Fatalf("message queued while model was processing; queue length = %d", got)
	}
	if got := hub.droppedEmotionMessage; got != 1 {
		t.Fatalf("dropped message count = %d, want 1", got)
	}
	hub.finishEmotionProcessing()
}

func TestOnChatQueueIsLimitedToOneBatch(t *testing.T) {
	hub := NewHub()
	for i := 0; i < maxEmotionBatchSize+3; i++ {
		hub.OnChat(ChatMessage{Text: "message"})
	}
	if got := len(hub.messageQueue); got != maxEmotionBatchSize {
		t.Fatalf("queued messages = %d, want cap %d", got, maxEmotionBatchSize)
	}
	if got := hub.droppedEmotionMessage; got != 3 {
		t.Fatalf("dropped message count = %d, want 3", got)
	}
}

func TestBeginEmotionProcessingDiscardsMessagesQueuedAfterBatch(t *testing.T) {
	hub := NewHub()
	hub.messageQueue <- queuedChatMessage{text: "late"}
	hub.beginEmotionProcessing()
	if got := len(hub.messageQueue); got != 0 {
		t.Fatalf("queued message survived processing start; queue length = %d", got)
	}
	if got := hub.droppedEmotionMessage; got != 1 {
		t.Fatalf("dropped message count = %d, want 1", got)
	}
	hub.finishEmotionProcessing()
}

func TestNextEmotionBatchFlushesAtSixteenMessages(t *testing.T) {
	hub := NewHub()
	for i := 0; i < maxEmotionBatchSize; i++ {
		hub.messageQueue <- queuedChatMessage{text: "message"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	batch, ok := hub.nextEmotionBatch(ctx)
	if !ok || len(batch) != maxEmotionBatchSize {
		t.Fatalf("batch size = %d, ok = %v, want %d", len(batch), ok, maxEmotionBatchSize)
	}
}

func TestNextEmotionBatchFlushesPartialBatchAfterFixedWait(t *testing.T) {
	hub := NewHub()
	hub.messageQueue <- queuedChatMessage{text: "message"}
	ctx, cancel := context.WithTimeout(context.Background(), 2*emotionBatchWait)
	defer cancel()

	start := time.Now()
	batch, ok := hub.nextEmotionBatch(ctx)
	elapsed := time.Since(start)
	if !ok || len(batch) != 1 {
		t.Fatalf("batch size = %d, ok = %v, want partial batch of 1", len(batch), ok)
	}
	if elapsed < emotionBatchWait || elapsed > 2*emotionBatchWait {
		t.Fatalf("partial batch waited %s, want about %s", elapsed, emotionBatchWait)
	}
}

func TestSnapshotIncludesCurrentGroupMood(t *testing.T) {
	hub := NewHub()
	hub.currentEmotions = emotionScores("joy")
	snapshot := hub.snapshotLocked()

	if snapshot.Type != "snapshot" || snapshot.Scores["joy"] != 1 {
		t.Fatalf("snapshot = %+v", snapshot)
	}
	snapshot.Scores["joy"] = 0
	if got := hub.currentEmotions["joy"]; got != 1 {
		t.Fatalf("snapshot shares mood storage with hub: joy score = %f", got)
	}
}

func TestNewHubSnapshotStartsNeutral(t *testing.T) {
	snapshot := NewHub().snapshotLocked()
	if snapshot.Scores["others"] != 1 {
		t.Fatalf("initial snapshot neutral score = %f, want 1", snapshot.Scores["others"])
	}
}
