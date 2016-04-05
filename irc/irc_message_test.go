package irc

import (
	"fmt"
	"testing"
)

type MessageTest struct {
	in      string
	command Command
	origin  string
	channel string
	text    string
}

var tests = []MessageTest{
	{
		in:      ":server NOTICE AUTH",
		command: Notice,
		origin:  "server",
		channel: "AUTH",
	},
	{
		in:      ":nick!~username@host PRIVMSG #channel :chat chat chat",
		command: Privmsg,
		origin:  "nick",
		channel: "#channel",
		text:    "chat chat chat",
	},
	{
		in:      ":server PING",
		command: Ping,
		origin:  "server",
	},
	{
		in:      ":nick!~username@host PRIVMSG #channel :ACTION emote",
		command: Privmsg,
		origin:  "nick",
		channel: "#channel",
		text:    "ACTION emote",
	},
	{
		in:      ":trapro!~trahari@75-145-17-54-Washington.hfc.comcastbusiness.net PRIVMSG gobot :HELLO",
		command: Privmsg,
		origin:  "trapro",
		channel: "gobot",
		text:    "HELLO",
	},
}

func TestBasicMessageParsing(t *testing.T) {
	for _, test := range tests {
		m := NewMessage(test.in)
		errors := verify(&test, m)
		if len(errors) > 0 {
			t.Errorf("NewMessage(%q) has the following errors: %+v", test.in, errors)
		}
	}
}

func verify(tt *MessageTest, mm *Message) []string {
	var errors []string
	if tt.command != mm.Command {
		errors = append(errors, fmt.Sprintf("command: got %q, want %q", mm.Command, tt.command))
	}

	if tt.origin != mm.Source {
		errors = append(errors, fmt.Sprintf("origin: got %q, want %q", mm.Source, tt.origin))
	}

	if tt.channel != mm.Channel {
		errors = append(errors, fmt.Sprintf("channel: got %q, want %q", mm.Channel, tt.channel))
	}

	if tt.text != mm.Text {
		errors = append(errors, fmt.Sprintf("text: got %q, want %q", mm.Text, tt.text))
	}

	return errors
}

var whoIsTests = []struct {
	Source string // input
	User   string
	Nick   string
}{
	{
		Source: ":foonick!~foouser@127-0-0-buh.foo.baz",
		User:   "foouser",
		Nick:   "foonick",
	},
	{
		Source: ":server",
		User:   "",
		Nick:   "server",
	},
}

func TestWhoIs(t *testing.T) {
	for _, test := range whoIsTests {
		user, nick := whoIs(test.Source)
		if user != test.User || nick != test.Nick {
			t.Errorf("whoIs(%q) => %q, %q, wanted %q, %q", test.Source, user, nick, test.User, test.Nick)
		}
	}

}

func TestSplitMsg(t *testing.T) {
	var tests = []struct {
		Message string
		Prefix  string
		Postfix string
	}{
		{
			Message: ":server NOTICE AUTH",
			Prefix:  "server NOTICE AUTH",
			Postfix: "",
		},
		{
			Message: ":nick!~username@host PRIVMSG #channel :chat chat chat",
			Prefix:  "nick!~username@host PRIVMSG #channel ",
			Postfix: "chat chat chat",
		},
	}

	for _, test := range tests {
		actualPre, actualPost := splitMsg(test.Message)
		if actualPre != test.Prefix || actualPost != test.Postfix {
			t.Errorf("splitMsg(%q) => %q, %q, wanted %q, %q", test.Message, actualPre, actualPost, test.Prefix, test.Postfix)
		}
	}
}
