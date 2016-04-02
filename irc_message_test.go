package youandmeandirc

import (
	"testing"
)

type IrcMessageTest struct {
	in      string
	command string
	origin  string
	channel string
	text    string
}

var tests = []IrcMessageTest{
	{":server NOTICE AUTH", "NOTICE", "server", "", ""},
	{":nick!~username@host PRIVMSG #channel :chat chat chat", "PRIVMSG", "nick", "#channel", "chat chat chat"},
	{":server PING", "PING", "server", "", ""},
	{":nick!~username@host PRIVMSG #channel :ACTION emote", "PRIVMSG", "nick", "#channel", "ACTION emote"},
	{":trapro!~trahari@75-145-17-54-Washington.hfc.comcastbusiness.net PRIVMSG gobot :HELLO", "PRIVMSG", "trapro", "", "HELLO"},
}

func TestBasicMessageParsing(t *testing.T) {
	for i, test := range tests {
		_, err := ParseMessage(test.in)
		if err != nil {
			t.Errorf("%d. ParseMessage(%q) => Error: %s", i, test.in, err)
			continue
		}

		// 	if !verify(&tt, actual) {
		// 		t.Errorf("%d. ParseMessage(%q) => %q, want %q ", i, test.in, actual, test)
		// 	}
		// }
	}
}

func verify(tt *IrcMessageTest, mm *IrcMessage) bool {
	expected := string(mm.Command) + mm.Origin + mm.Channel + mm.Text
	actual := tt.command + tt.origin + tt.channel + tt.text
	if expected != actual {
		return false
	}
	return true
}
