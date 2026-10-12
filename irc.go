package main

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const twitchIRCWebSocket = "wss://irc-ws.chat.twitch.tv:443"

// Reconecta automáticamente si se cae la conexión.
func runIRC(channel string, hub *Hub) {
	for {
		err := connectIRC(channel, hub)
		log.Printf("IRC WebSocket desconectado (%v), reintentando en 5s...", err)
		time.Sleep(5 * time.Second)
	}
}

func connectIRC(channel string, hub *Hub) error {
	conn, _, err := websocket.DefaultDialer.Dial(twitchIRCWebSocket, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	nick := "justinfan" + string(rune('0'+rand.Intn(10))) + string(rune('0'+rand.Intn(10))) + string(rune('0'+rand.Intn(10))) + string(rune('0'+rand.Intn(10))) + string(rune('0'+rand.Intn(10)))
	for _, command := range []string{
		"CAP REQ :twitch.tv/tags twitch.tv/commands",
		"PASS SCHMOOPIIE",
		"NICK " + nick,
		"JOIN #" + strings.ToLower(channel),
	} {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(command+"\r\n")); err != nil {
			return err
		}
	}
	log.Printf("Conectando por IRC WebSocket anónimo al chat de #%s", channel)

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		if messageType != websocket.TextMessage {
			continue
		}
		for _, line := range strings.Split(string(data), "\r\n") {
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, "PING") {
				if err := conn.WriteMessage(websocket.TextMessage, []byte(strings.Replace(line, "PING", "PONG", 1)+"\r\n")); err != nil {
					return err
				}
				continue
			}
			if strings.Contains(line, " NOTICE ") {
				log.Printf("Aviso de Twitch IRC: %s", line)
			}
			if msg, ok := parsePrivmsg(line); ok {
				hub.OnChat(msg)
			}
		}
	}
}

// Formato: @tags :nick!nick@nick.tmi.twitch.tv PRIVMSG #canal :texto
func parsePrivmsg(line string) (ChatMessage, bool) {
	tags := map[string]string{}

	if strings.HasPrefix(line, "@") {
		i := strings.Index(line, " ")
		if i < 0 {
			return ChatMessage{}, false
		}
		for _, kv := range strings.Split(line[1:i], ";") {
			p := strings.SplitN(kv, "=", 2)
			if len(p) == 2 {
				tags[p[0]] = p[1]
			}
		}
		line = line[i+1:]
	}

	if !strings.HasPrefix(line, ":") {
		return ChatMessage{}, false
	}
	parts := strings.SplitN(line[1:], " ", 4)
	if len(parts) < 4 || parts[1] != "PRIVMSG" {
		return ChatMessage{}, false
	}

	nick := parts[0]
	if i := strings.Index(nick, "!"); i > 0 {
		nick = nick[:i]
	}
	user := tags["display-name"]
	if user == "" {
		user = nick
	} else {
		user = unescapeIRCValue(user)
	}
	color := tags["color"]
	if color == "" {
		color = "#afafaf"
	}
	return ChatMessage{
		User:   user,
		UserID: tags["user-id"],
		Color:  color,
		Text:   strings.TrimPrefix(parts[3], ":"),
	}, true
}

func unescapeIRCValue(value string) string {
	return strings.NewReplacer(
		`\\`, `\`,
		`\s`, " ",
		`\:`, ";",
		`\r`, "\r",
		`\n`, "\n",
	).Replace(value)
}
