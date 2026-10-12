package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"strings"
)

func main() {
	channel := flag.String("channel", "", "canal de Twitch (sin #)")
	addr := flag.String("addr", ":8080", "dirección del servidor")
	flag.Parse()

	if *channel == "" {
		log.Fatal("uso: go run . -channel nombredelcanal")
	}

	hub := NewHub()
	go runIRC(strings.ToLower(*channel), hub)
	go hub.RunEmotionLoop(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", hub.ServeWS)

	log.Printf("WebSocket en ws://localhost%s/ws (canal: %s)", *addr, *channel)
	log.Fatal(http.ListenAndServe(*addr, mux))
}
