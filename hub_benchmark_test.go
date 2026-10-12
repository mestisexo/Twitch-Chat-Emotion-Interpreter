package main

import (
	"testing"
)

func BenchmarkHubOnChat(b *testing.B) {
	hub := NewHub()
	hub.messageQueue = make(chan queuedChatMessage, 1)
	message := ChatMessage{Text: "Estoy muy feliz, ganamos la partida"}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hub.OnChat(message)
		<-hub.messageQueue
	}
}

func BenchmarkApplyEmotionResults(b *testing.B) {
	hub := NewHub()
	results := make([]map[string]float64, maxEmotionBatchSize)
	for i := range results {
		results[i] = emotionScores("joy")
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hub.applyEmotionResults(results)
	}
}
