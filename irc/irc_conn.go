package irc

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

// Conn represents a connection to an IRC server.
type Conn struct {
	Username string
	Pass     string
	Nick     string
	Realname string

	conn   net.Conn
	reader *bufio.Reader

	host, port string
}

/// Network-related.

// send transmits a command to the currently connected server.
func (irc Conn) send(s string) (int, error) {
	log.Println("=>", s)
	return fmt.Fprintln(irc.conn, s)
}

// register sends the PASS, NICK, and USER commands.
func (irc Conn) register() {
	messages := []string{
		irc.newPassMsg(irc.Pass),
		irc.newNickMsg(irc.Nick),
		irc.newUserMsg(irc.Username, irc.Realname),
	}

	for _, m := range messages {
		if _, err := irc.send(m); err != nil {
			log.Println("Error:", err)
		}
	}
}

/// Composing various kinds of messages.

func (irc Conn) newPassMsg(pass string) string {
	return fmt.Sprintf("PASS %v", pass)
}

func (irc Conn) newNickMsg(nick string) string {
	return fmt.Sprintf("NICK %v", nick)
}

func (irc Conn) newUserMsg(username, realname string) string {
	return fmt.Sprintf("USER %v * * :%v", username, realname)
}

func (irc Conn) newJoinMsg(channel string) string {
	return fmt.Sprintf("JOIN %v", channel)
}

func (irc Conn) newPongMsg(daemon string) string {
	return fmt.Sprintf("PONG %v", daemon)
}

/// Public methods.

// Sends a message to a channel.
func (irc Conn) Say(channel, chat string) (int, error) {
	msg := fmt.Sprintf("PRIVMSG %v :%v", channel, chat)
	return irc.send(msg)
}

// Joins a given channel.
func (irc Conn) Join(channel string) (int, error) {
	cmd := fmt.Sprintf("JOIN %v", channel)
	return irc.send(cmd)
}

func (irc Conn) Names(channel string) (int, error) {
	cmd := fmt.Sprintf("NAMES %v", channel)
	return irc.send(cmd)
}

// Reads a single message from the server's output.
func (irc Conn) Read() (m *Message, err error) {
	s, err := irc.reader.ReadString('\n')
	log.Println("<=", strings.TrimRight(s, "\r\n"))
	if err != nil {
		log.Println("Error reading from server:", err)
		return nil, err
	}

	m, err = ParseMessage(s)
	if err != nil {
		log.Println("Error parsing server message:", err)
		return nil, err
	}

	return
}

func (irc Conn) Pong(daemon string) (int, error) {
	// FIXME: shouldn't this be handled automatically?
	cmd := irc.newPongMsg(daemon)
	return irc.send(cmd)
}

// Connects to IRC with the given connection.
func (irc *Conn) Connect(c net.Conn) {
	irc.conn = c
	irc.reader = bufio.NewReader(irc.conn)
	irc.register()
	return
}
