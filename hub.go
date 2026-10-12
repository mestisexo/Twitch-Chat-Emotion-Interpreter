package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

const (
	emotionLabelCount   = 7
	maxEmotionBatchSize = 16
	emotionBatchWait    = time.Second
	emotionRetryWait    = 5 * time.Second
	clientSendBuffer    = 128
)

var emotionLabels = []string{"others", "joy", "sadness", "anger", "surprise", "disgust", "fear"}

type queuedChatMessage struct {
	text string
}

type emotionClassifier interface {
	Classify(context.Context, []string) ([]map[string]float64, error)
}

type Hub struct {
	mu                    sync.Mutex
	queueMu               sync.Mutex
	clients               map[*wsClient]bool
	messageQueue          chan queuedChatMessage
	processingEmotion     bool
	droppedEmotionMessage uint64
	classifier            emotionClassifier
	currentEmotions       map[string]float64
}

type wsClient struct {
	conn *websocket.Conn
	send chan []byte
}

type ChatMessage struct {
	User   string
	Color  string
	UserID string
	Text   string
}

type overlayEvent struct {
	Type   string             `json:"type"`
	Scores map[string]float64 `json:"scores,omitempty"`
}

func NewHub() *Hub {
	return &Hub{
		clients:         make(map[*wsClient]bool),
		messageQueue:    make(chan queuedChatMessage, maxEmotionBatchSize),
		classifier:      newEmotionHTTPClient(os.Getenv("EMOTION_SERVICE_URL")),
		currentEmotions: map[string]float64{"others": 1},
	}
}

// OnChat queues message text for group-level emotion analysis.
func (h *Hub) OnChat(message ChatMessage) {
	text := strings.TrimSpace(message.Text)
	if text == "" {
		return
	}
	h.queueMu.Lock()
	defer h.queueMu.Unlock()
	if h.processingEmotion {
		h.droppedEmotionMessage++
		return
	}
	select {
	case h.messageQueue <- queuedChatMessage{text: text}:
	default:
		h.droppedEmotionMessage++
	}
}

func (h *Hub) snapshotLocked() overlayEvent {
	return overlayEvent{Type: "snapshot", Scores: cloneEmotionScores(h.currentEmotions)}
}

func cloneEmotionScores(scores map[string]float64) map[string]float64 {
	if scores == nil {
		return nil
	}
	copy := make(map[string]float64, len(scores))
	for label, score := range scores {
		copy[label] = score
	}
	return copy
}

func (h *Hub) RunEmotionLoop(ctx context.Context) {
	for {
		batch, ok := h.nextEmotionBatch(ctx)
		if !ok {
			return
		}
		h.beginEmotionProcessing()
		for {
			scores, err := h.classifier.Classify(ctx, batch)
			if err == nil {
				h.applyEmotionResults(scores)
				break
			}
			log.Printf("No se pudo analizar el tono del chat: %v; se reintentará en %s", err, emotionRetryWait)
			timer := time.NewTimer(emotionRetryWait)
			select {
			case <-ctx.Done():
				timer.Stop()
				h.finishEmotionProcessing()
				return
			case <-timer.C:
			}
		}
		h.finishEmotionProcessing()
	}
}

func (h *Hub) nextEmotionBatch(ctx context.Context) ([]string, bool) {
	select {
	case <-ctx.Done():
		return nil, false
	case first := <-h.messageQueue:
		batch := make([]string, 0, maxEmotionBatchSize)
		batch = append(batch, first.text)
		if len(batch) == maxEmotionBatchSize {
			return batch, true
		}
		timer := time.NewTimer(emotionBatchWait)
		defer timer.Stop()
		for len(batch) < maxEmotionBatchSize {
			select {
			case <-ctx.Done():
				return nil, false
			case message := <-h.messageQueue:
				batch = append(batch, message.text)
			case <-timer.C:
				return batch, true
			}
		}
		return batch, true
	}
}

func (h *Hub) beginEmotionProcessing() {
	h.queueMu.Lock()
	defer h.queueMu.Unlock()
	h.processingEmotion = true
	for {
		select {
		case <-h.messageQueue:
			h.droppedEmotionMessage++
		default:
			return
		}
	}
}

func (h *Hub) finishEmotionProcessing() {
	h.queueMu.Lock()
	h.processingEmotion = false
	dropped := h.droppedEmotionMessage
	h.droppedEmotionMessage = 0
	h.queueMu.Unlock()
	if dropped > 0 {
		log.Printf("Se descartaron %d mensajes mientras el análisis emocional estaba ocupado o la cola llena", dropped)
	}
}

func (h *Hub) applyEmotionResults(results []map[string]float64) {
	if len(results) == 0 {
		return
	}
	for _, scores := range results {
		if err := validateEmotionScores(scores); err != nil {
			log.Printf("Respuesta emocional inválida: %v", err)
			return
		}
	}

	var totals [emotionLabelCount]float64
	for _, scores := range results {
		for index, label := range emotionLabels {
			totals[index] += scores[label]
		}
	}
	averages := make(map[string]float64, len(emotionLabels))
	for index, label := range emotionLabels {
		averages[label] = totals[index] / float64(len(results))
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	h.currentEmotions = averages
	h.broadcastEventLocked(overlayEvent{Type: "mood", Scores: averages})
}

func validateEmotionScores(scores map[string]float64) error {
	total := 0.0
	for _, label := range emotionLabels {
		score, ok := scores[label]
		if !ok || math.IsNaN(score) || math.IsInf(score, 0) || score < 0 || score > 1 {
			return fmt.Errorf("puntaje inválido para %q", label)
		}
		total += score
	}
	if total < 0.99 || total > 1.01 {
		return fmt.Errorf("la suma de probabilidades es %.4f, se esperaba 1", total)
	}
	return nil
}

// Events are queued while holding the hub lock; network writes happen in each client's writer goroutine.
func (h *Hub) broadcastEventLocked(event overlayEvent) {
	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("No se pudo preparar el evento WebSocket: %v", err)
		return
	}
	for client := range h.clients {
		select {
		case client.send <- payload:
		default:
			h.removeClientLocked(client)
		}
	}
}

func (h *Hub) removeClientLocked(client *wsClient) {
	if !h.clients[client] {
		return
	}
	delete(h.clients, client)
	close(client.send)
	_ = client.conn.Close()
}

func (h *Hub) writePump(client *wsClient) {
	for payload := range client.send {
		_ = client.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		if err := client.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			h.mu.Lock()
			h.removeClientLocked(client)
			h.mu.Unlock()
			return
		}
	}
}

func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	client := &wsClient{conn: conn, send: make(chan []byte, clientSendBuffer)}

	h.mu.Lock()
	h.clients[client] = true
	payload, err := json.Marshal(h.snapshotLocked())
	if err != nil {
		h.removeClientLocked(client)
		h.mu.Unlock()
		return
	}
	client.send <- payload
	h.mu.Unlock()
	go h.writePump(client)

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}

	h.mu.Lock()
	h.removeClientLocked(client)
	h.mu.Unlock()
}
