package youandmeandirc

import (
	"testing"
)

type IrcMessageTest struct {
	in string
	command string
	origin string
	channel string
	text string
}

var tests = []IrcMessageTest { 
  {":server NOTICE AUTH", "NOTICE", "server", "", ""},
  {":nick!~username@host PRIVMSG #channel :chat chat chat", "PRIVMSG", "nick", "#channel", "chat chat chat"},
  {":server PING", "PING", "server", "", ""},
  {":nick!~username@host PRIVMSG #channel :ACTION emote", "PRIVMSG", "nick", "#channel", "ACTION emote"},
	{":trapro!~trahari@75-145-17-54-Washington.hfc.comcastbusiness.net PRIVMSG gobot :HELLO", "PRIVMSG", "trapro", "", "HELLO"},
}

func TestBasicMessageParsing(t *testing.T) {
	for i, tt := range tests {
		ircMsg, err := ParseMessage(tt.in)
		if err != nil {
			t.Errorf("%d. ParseMessage(%q) => Error: ", i, tt.in, err)
			continue
		}

		if !verify(&tt, ircMsg) {
			t.Errorf("%d. ParseMessage(%q) => %q, want %q ", i, tt.in, ircMsg, tt)
		}
	}
}

func verify(tt *IrcMessageTest, mm *IrcMessage) bool {
	expected := mm.Command + mm.Origin + mm.Channel + mm.Text
	actual := tt.command + tt.origin + tt.channel + tt.text
	if expected != actual {
		return false
	}
	return true
}
