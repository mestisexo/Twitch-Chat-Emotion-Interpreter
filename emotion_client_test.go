package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEmotionHTTPClientSendsBatchAndReadsScores(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("request method = %q, want POST", r.Method)
		}
		var request emotionRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if len(request.Messages) != 2 || request.Messages[1] != "mensaje dos" {
			t.Fatalf("request messages = %v", request.Messages)
		}
		_ = json.NewEncoder(w).Encode(emotionResponse{
			Predictions: []map[string]float64{emotionScores("joy"), emotionScores("sadness")},
		})
	}))
	defer server.Close()

	client := newEmotionHTTPClient(server.URL)
	predictions, err := client.Classify(context.Background(), []string{"mensaje uno", "mensaje dos"})
	if err != nil {
		t.Fatalf("Classify returned error: %v", err)
	}
	if len(predictions) != 2 || predictions[1]["sadness"] != 1 {
		t.Fatalf("predictions = %v", predictions)
	}
}

func TestEmotionHTTPClientReportsNonSuccessStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "model unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := newEmotionHTTPClient(server.URL)
	if _, err := client.Classify(context.Background(), []string{"mensaje"}); err == nil {
		t.Fatal("Classify succeeded despite service returning HTTP 503")
	}
}
