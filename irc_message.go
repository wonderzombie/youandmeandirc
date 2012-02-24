package youandmeandirc

import (
	"fmt"
	"strings"
)

type IrcMessage struct {
	Raw string
	Prefix string
	Command string
	Params string
	Channel string
}

const (
	NOTICE = iota
	PING
	PRIVMSG
	MODE
	JOIN
)

type IrcMessageError struct {
	RawMsg string // the offending message
	Reason string // reason for the error
}

func (e *IrcMessageError) Error() string {
	return e.Reason + ": " + e.RawMsg
}

func ParseMessage(msg string) (*IrcMessage, error) {
	m := new(IrcMessage)		
	if err := m.init(msg) ; err != nil {
		return nil, err
	}
	return m, nil
}

func (m *IrcMessage) init(rawMsg string) error {
	m.Raw = rawMsg

	msgParts := strings.Split(m.Raw, " ")
	if len(msgParts) > 2 {
		return &IrcMessageError{m.Raw, "Unable to parse."}
	}

	m.Command = msgParts[1]

	switch m.Command {
		case "PRIVMSG":
			fmt.Println(m.Command)
	}

	// We can infer server name based on absence of ! and @.

	// :server NOTICE AUTH 
	// :nick!~username@host PRIVMSG #channel :chat
	// :nick!~username@host PRIVMSG #channel :ACTION emote

	return nil
}

func (m *IrcMessage) parsePrefix(prefix string) []string {
	return []string{}
}
