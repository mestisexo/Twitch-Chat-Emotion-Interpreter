package main

import "testing"

func TestParsePrivmsgIncludesDisplayName(t *testing.T) {
	message, ok := parsePrivmsg("@display-name=Viewer\\sOne;color=#3366cc;user-id=123 :viewerone!viewerone@viewerone.tmi.twitch.tv PRIVMSG #channel :hello")
	if !ok {
		t.Fatal("parsePrivmsg did not recognize a valid PRIVMSG")
	}
	if message.User != "Viewer One" {
		t.Fatalf("display name = %q, want %q", message.User, "Viewer One")
	}
	if message.UserID != "123" {
		t.Fatalf("user ID = %q, want %q", message.UserID, "123")
	}
	if message.Text != "hello" {
		t.Fatalf("message text = %q, want %q", message.Text, "hello")
	}
}
