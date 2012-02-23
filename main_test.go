package main

import (
	"bytes"
	"strings"
	"testing"
)

const (
	FakePrivMsg = ":fakenick!~fakeuser@12-34-56-78.foo.bar.baz.net PRIVMSG #trago :foo bar baz"
)

func TestSendMsg(t *testing.T) {
	fakeConn := new(bytes.Buffer)

	msg := "foo bar baz"
	writeResp(fakeConn, "#fake", msg)
	actual := fakeConn.String()

	if !strings.Contains(actual, msg) {
		t.Log("Messages did not match.")
		t.Log("Expected: %v, Actual: %v")
		t.Fail()
	}
}

func TestChatFromMsg(t *testing.T) {
	chat := chatFromMsg(FakePrivMsg)
	if !strings.Contains(chat, ":foo bar baz") {
		t.Log("Failed to get chat message from PRIVMSG")
		t.Log("Chat: %v")
		t.Fail()
	}
}
