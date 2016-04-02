package youandmeandirc

import (
	"strings"
)

type ircCmd int

const (
	CmdPing ircCmd = iota
	CmdPrivmsg
	CmdMode
	CmdPart
	CmdJoin
	CmdNotice
	CmdNum
)

// Lookup table for commands against IDs.
var cmdMap = map[string]ircCmd{
	"PING":    CmdPing,
	"PRIVMSG": CmdPrivmsg,
	"MODE":    CmdMode,
	"PART":    CmdPart,
	"JOIN":    CmdJoin,
	"NOTICE":  CmdNotice,
	"###":     CmdNum,
}

func (c ircCmd) String() string {
	for k, v := range cmdMap {
		if c == v {
			return k
		}
	}
	return ""
}

// IrcMessage is a structured representation of an IRC message.
type IrcMessage struct {
	Raw     string
	Command ircCmd
	Channel string   // Channel which the message belongs to, if any.
	Origin  string   // Nick or server which originated the message.
	Text    string   // Text of the chat.
	Code    string   // Command code.
	Args    []string // Misc params.
	User    string
	Nick    string
}

// IrcMessageError is returned when an IRC message cannot be parsed.
type IrcMessageError struct {
	Raw    string // the offending message
	Reason string // reason for the error
}

func (e *IrcMessageError) Error() string {
	return e.Reason + ": " + e.Raw
}

// ParseMessage returns a new IrcMessage populated with the data from a raw IRC message.
func ParseMessage(msg string) (*IrcMessage, error) {
	m := new(IrcMessage)
	if err := m.init(msg); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *IrcMessage) newErr(reason string) (e *IrcMessageError) {
	return &IrcMessageError{
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

func (m *IrcMessage) init(msg string) (e *IrcMessageError) {
	// Trim trailing nonprinting character.
	// This may need more delicate treatment if it turns out multiline IRC
	// messages are a problem.
	msg = strings.Trim(strings.TrimSpace(msg), "\x01")
	command, content := splitMsg(msg)
	commandTokens := strings.Fields(command)

	fst := commandTokens[0]
	if cmd, ok := cmdMap[fst]; ok && cmd == CmdPing {
		m.Origin = content
		m.Command = cmd
		return
	}

	origin, cmd := commandTokens[0], commandTokens[1]
	m.Nick, m.User = whoIs(origin)
	m.Origin = m.User

	cmdId, ok := cmdMap[cmd]
	if !ok {
		m.Command = CmdNum
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
	case CmdJoin:
		m.Channel = content
	case CmdPrivmsg, CmdMode, CmdNotice:
		m.Channel = commandTokens[2]
		m.Text = content
	case CmdPart:
		m.Channel = commandTokens[2]
	}

	return
}

func (m *IrcMessage) Matches(cmds []ircCmd) bool {
	for _, c := range cmds {
		if c == m.Command {
			return true
		}
	}
	return false
}

func (m *IrcMessage) HasText(s string) bool {
	return strings.Index(m.Text, s) >= 0
}

func (m *IrcMessage) TextHasAny(ss []string) bool {
	for _, s := range ss {
		if m.HasText(s) {
			return true
		}
	}
	return false
}
