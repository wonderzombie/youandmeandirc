package irc

import (
	"strings"
)

type Command int

const (
	Ping Command = iota
	Privmsg
	Mode
	Part
	Join
	Notice
	Num // numeric commands
)

// Lookup table for commands against IDs.
// There's no key for Num on purpose.
var CommandIndex = map[string]Command{
	"PING":    Ping,
	"PRIVMSG": Privmsg,
	"MODE":    Mode, // recognized but ignored
	"PART":    Part, // recognized but ignored
	"JOIN":    Join,
	"NOTICE":  Notice, // recognized but ignored
}

func (c Command) String() string {
	for k, v := range CommandIndex {
		if c == v {
			return k
		}
	}
	return ""
}

// Message is a structured representation of an IRC message.
type Message struct {
	Raw     string
	Command Command
	Channel string   // Channel which the message belongs to, if any.
	Source  string   // Nick or server which originated the message.
	Text    string   // Text of the chat.
	Code    string   // Command code.
	Args    []string // Misc params.
	User    string
	Nick    string
}

// MessageError is returned when an IRC message cannot be parsed.
type MessageError struct {
	Raw    string // the offending message
	Reason string // reason for the error
}

func (e *MessageError) Error() string {
	return e.Reason + ": " + e.Raw
}

// ParseMessage returns a new Message populated with the data from a raw IRC message.
func ParseMessage(msg string) (*Message, error) {
	m, err := NewMessage(msg)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func splitMsg(msg string) (string, string) {
	msg = strings.Trim(msg, ":")
	if strings.Index(msg, ":") == -1 {
		return msg, ""
	}

	pieces := strings.SplitN(msg, ":", 2)
	return pieces[0], pieces[1]
}

func fromServer(prefix string) bool {
	// If there's no ! in the prefix, that means it's a server name.
	return strings.Index(prefix, "!") == -1
}

func whoIs(prefix string) (string, string) {
	prefix = strings.Trim(prefix, ":")
	if fromServer(prefix) {
		return "", prefix
	}

	// Technically, these are not all required, such as if it's just a nickname, without username or host.
	// We'll delineate this way:
	//     somenick!~someuser@hostname
	//     0       ^^        ^
	// In the absence of ! ~ or @ we'll clip it to bangPos, which itself defaults to the length of the string.
	bangPos := strings.Index(prefix, "!")
	if bangPos == -1 {
		bangPos = len(prefix) + 1
	}
	tildePos := strings.Index(prefix, "~")
	if tildePos == -1 {
		tildePos = bangPos
	}
	atPos := strings.Index(prefix, "@")
	if atPos == -1 {
		atPos = bangPos
	}

	nick := prefix[0 : bangPos-1]
	user := prefix[tildePos+1 : atPos-1]

	return nick, user
}

func NewMessage(msg string) (*Message, *MessageError) {
	m := &Message{}
	// Trim trailing nonprinting character.
	msg = strings.Trim(strings.TrimSpace(msg), "\x01")
	command, content := splitMsg(msg)
	commandTokens := strings.Fields(command)

	// Ping is easy to rule out since it's two tokens.
	tok := commandTokens[0]
	if cmd := CommandIndex[tok]; cmd == Ping {
		m.Source = content
		m.Command = cmd
		return m, nil
	}

	source, cmd := commandTokens[0], commandTokens[1]
	m.Nick, m.User = whoIs(source)
	m.Source = m.User

	id, ok := CommandIndex[cmd]
	// A miss means this is probably a numeric code.
	// https://tools.ietf.org/html/rfc2812#section-5.1
	if !ok {
		m.Command = Num
		m.Code = cmd
		// Sometimes there's extra stuff, so don't drop it.
		if len(commandTokens) > 2 {
			m.Args = commandTokens[2:]
		}
		if content != "" {
			m.Text = content
		}
		return m, nil
	}

	// This is a user-relevant event.
	m.Command = id
	switch id {
	case Join:
		m.Channel = content
	case Privmsg, Mode, Notice:
		m.Channel = commandTokens[2]
		m.Text = content
	case Part:
		m.Channel = commandTokens[2]
	}

	return m, nil
}

func (m *Message) MatchesAny(cmds []Command) bool {
	for _, c := range cmds {
		if c == m.Command {
			return true
		}
	}
	return false
}

func (m *Message) TextHas(s string) bool {
	return strings.Contains(m.Text, s)
}

func (m *Message) TextHasAny(ss []string) bool {
	for _, s := range ss {
		if m.TextHas(s) {
			return true
		}
	}
	return false
}

// TextSortaHas does a case-insensitive comparison.
func (m *Message) TextSortaHas(s string) bool {
	l, r := strings.ToLower(m.Text), strings.ToLower(s)
	return strings.Contains(l, r)
}

func (m *Message) TextSortaHasAny(ss []string) bool {
	for _, s := range ss {
		if m.TextSortaHas(s) {
			return true
		}
	}
	return false
}
