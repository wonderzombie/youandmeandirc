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
	username string
	pass     string
	nick     string
	realname string

	conn   net.Conn
	reader *bufio.Reader

	host, port string
}

/// Network-related.

// send transmits a command to the currently connected server.
func (irc Conn) send(s string) error {
	log.Println("=>", s)
	_, err := fmt.Fprintln(irc.conn, s)
	return err
}

func (irc Conn) sendfln(format string, a ...interface{}) error {
	format += "\n"
	log.Printf("=> %s", fmt.Sprintf(format, a...))
	_, err := fmt.Fprintf(irc.conn, format, a...)
	return err
}

// register sends the PASS, NICK, and USER commands.
func (irc Conn) register() {
	messages := []string{
		// Pass, Nick, then User. This order matters!
		newPassMsg(irc.pass),
		newNickMsg(irc.nick),
		newUserMsg(irc.username, irc.realname),
	}

	for _, m := range messages {
		if err := irc.send(m); err != nil {
			log.Println("Error:", err)
		}
	}
}

/// Composing various kinds of messages.

func newPassMsg(pass string) string {
	return fmt.Sprintf("PASS %v", pass)
}

func newNickMsg(nick string) string {
	return fmt.Sprintf("NICK %v", nick)
}

func newUserMsg(username, realname string) string {
	return fmt.Sprintf("USER %v * * :%v", username, realname)
}

/// Public methods.

// Sends a message to a channel.
func (irc Conn) Say(channel, chat string) error {
	return irc.sendfln("PRIVMSG %v :%v", channel, chat)
}

// Joins a given channel.
func (irc Conn) Join(channel string) error {
	return irc.sendfln("JOIN %v", channel)
}

func (irc Conn) Names(channel string) error {
	return irc.sendfln("NAMES %v", channel)
}

func (irc Conn) SetNick(nick string) error {
	if err := irc.sendfln(newNickMsg(nick)); err != nil {
		return err
	}
	irc.nick = nick
	return nil
}

// Reads a single message from the server's output.
func (irc Conn) Read() (*Message, error) {
	s, err := irc.reader.ReadString('\n')
	log.Println("<=", strings.TrimRight(s, "\r\n"))
	if err != nil {
		log.Println("Error reading from server:", err)
		return nil, err
	}
	return NewMessage(s), nil
}

func (irc Conn) Pong(daemon string) error {
	// FIXME: shouldn't this be handled automatically?
	// Specifically, this is protocol-level stuff. We could (should?) hide this from the user.
	return irc.sendfln("PONG %v", daemon)
}

func (irc Conn) Nick() string {
	return irc.nick
}

// Should probably DTRT IRC protocol-wise, like sending a quit message.
func (irc Conn) Disconnect() error {
	return irc.conn.Close()
}

// Connect initiates the IRC protocol with the given credentails.
func Connect(n net.Conn, nick, realname, username, pass string) *Conn {
	c := &Conn{
		conn:     n,
		reader:   bufio.NewReader(n),
		nick:     nick,
		realname: realname,
		username: username,
		pass:     pass,
	}
	c.register()
	return c
}
