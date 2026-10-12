package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultEmotionServiceURL = "http://127.0.0.1:8090/classify"

type emotionHTTPClient struct {
	url    string
	client *http.Client
}

type emotionRequest struct {
	Messages []string `json:"messages"`
}

type emotionResponse struct {
	Predictions []map[string]float64 `json:"predictions"`
}

func newEmotionHTTPClient(url string) *emotionHTTPClient {
	if strings.TrimSpace(url) == "" {
		url = defaultEmotionServiceURL
	}
	return &emotionHTTPClient{
		url: url,
		client: &http.Client{
			Timeout: 45 * time.Second,
		},
	}
}

func (c *emotionHTTPClient) Classify(ctx context.Context, messages []string) ([]map[string]float64, error) {
	payload, err := json.Marshal(emotionRequest{Messages: messages})
	if err != nil {
		return nil, fmt.Errorf("codificar lote: %w", err)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("crear solicitud al clasificador: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("conectar con el clasificador local: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(io.LimitReader(response.Body, 2048))
		if readErr != nil {
			return nil, fmt.Errorf("clasificador respondió HTTP %d y no se pudo leer el error: %w", response.StatusCode, readErr)
		}
		return nil, fmt.Errorf("clasificador respondió HTTP %d: %s", response.StatusCode, strings.TrimSpace(string(body)))
	}

	var result emotionResponse
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&result); err != nil {
		return nil, fmt.Errorf("decodificar respuesta del clasificador: %w", err)
	}
	if len(result.Predictions) != len(messages) {
		return nil, fmt.Errorf("el clasificador devolvió %d resultados para %d mensajes", len(result.Predictions), len(messages))
	}
	return result.Predictions, nil
}
