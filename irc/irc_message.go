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
	Num
)

// Lookup table for commands against IDs.
var CommandIndex = map[string]Command{
	"PING":    Ping,
	"PRIVMSG": Privmsg,
	"MODE":    Mode,
	"PART":    Part,
	"JOIN":    Join,
	"NOTICE":  Notice,
	"###":     Num,
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
	Origin  string   // Nick or server which originated the message.
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
	m := new(Message)
	if err := m.init(msg); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Message) newErr(reason string) (e *MessageError) {
	return &MessageError{
		m.Raw,
		reason,
	}
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

func whoIs(prefix string) (nick, user string) {
	prefix = strings.Trim(prefix, ":")
	if fromServer(prefix) {
		return "", prefix
	}

	pp := strings.Split(prefix, "@")
	ss := strings.Split(pp[0], "!")
	nick = ss[0]
	user = strings.Trim(ss[1], "~")
	return
}

func (m *Message) init(msg string) (e *MessageError) {
	// Trim trailing nonprinting character.
	// This may need more delicate treatment if it turns out multiline IRC
	// messages are a problem.
	msg = strings.Trim(strings.TrimSpace(msg), "\x01")
	command, content := splitMsg(msg)
	commandTokens := strings.Fields(command)

	fst := commandTokens[0]
	if cmd, ok := CommandIndex[fst]; ok && cmd == Ping {
		m.Origin = content
		m.Command = cmd
		return
	}

	origin, cmd := commandTokens[0], commandTokens[1]
	m.Nick, m.User = whoIs(origin)
	m.Origin = m.User

	cmdId, ok := CommandIndex[cmd]
	if !ok {
		m.Command = Num
		m.Code = cmd
		if len(commandTokens) > 2 {
			m.Args = commandTokens[2:]
		}
		if content != "" {
			m.Text = content
		}
		return
	}

	m.Command = cmdId
	switch cmdId {
	case Join:
		m.Channel = content
	case Privmsg, Mode, Notice:
		m.Channel = commandTokens[2]
		m.Text = content
	case Part:
		m.Channel = commandTokens[2]
	}

	return
}

func (m *Message) Matches(cmds []Command) bool {
	for _, c := range cmds {
		if c == m.Command {
			return true
		}
	}
	return false
}

func (m *Message) HasText(s string) bool {
	return strings.Index(m.Text, s) >= 0
}

func (m *Message) TextHasAny(ss []string) bool {
	for _, s := range ss {
		if m.HasText(s) {
			return true
		}
	}
	return false
}
