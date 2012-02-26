package youandmeandirc

import (
	"strings"
)

// IrcMessage is a structured representation of an IRC message.
type IrcMessage struct {
	Raw string
	Command string
	Channel string // Channel which the message belongs to, if any.
	Origin string // Nick or server which originated the message.
	Text string // Text of the chat.
	Params []string // Misc params.
	Target []string
}

// IrcMessageError is returned when an IRC message cannot be parsed.
type IrcMessageError struct {
	Raw string // the offending message
	Reason string // reason for the error
}

func (e *IrcMessageError) Error() string {
	return e.Reason + ": " + e.Raw
}

// ParseMessage returns a new IrcMessage populated with the data from a raw IRC message.
func ParseMessage(msg string) (*IrcMessage, error) {
	m := new(IrcMessage)		
	if err := m.init(msg) ; err != nil {
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

func (m *IrcMessage) init(msg string) (e *IrcMessageError) {
	m.Raw = msg

	prefix, rest := firstToken(msg)

	if prefix == "PING" {
		if rest == "" {
			return m.newErr("Received PING without daemon")
		}
		m.Command = prefix
		// Trim the : before the daemon name.
		m.Params = append(m.Params, rest[1:])
		return
	}

	// Take out the "foo" in "foo!~username@host"
	i := strings.Index(prefix, "!")
	if i > -1 {
		m.Origin = prefix[1:i]
	} else {
		m.Origin = prefix[1:]
	}

	m.Command, rest = firstToken(rest)

	// Bail if there's no more to parse.
	if rest == "" {
		return
	}

	f, r := firstToken(rest)
	// For now, we just care about PRIVMSG.
	if m.Command == "PRIVMSG" {
		if f[0] == '#' && r[0] == ':' {
			// PRIVMSG directed at a channel.
			m.Channel = f
			m.Text = r[1:]
		} else {
			// PRIVMSG directly to a user (i.e. me).
			m.Params = append(m.Params, f)
			m.Text = r[1:]
		}
	}

	return
}

func firstToken(s string) (first, rest string) {
	x := strings.SplitN(s, " ", 2)
	first = x[0]
	if len(x) == 2 {
		rest = x[1]
	}
	return
}
