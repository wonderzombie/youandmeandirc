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

	host string
	port string
}

// send transmits a command to the currently connected server.
func (irc Conn) send(s string) error {
	log.Println("=>", s)
	_, err := fmt.Fprintln(irc.conn, s)
	return err
}

// sendfln is a thin wrapper around *printf.
func (irc Conn) sendfln(format string, a ...interface{}) error {
	msg := fmt.Sprintf(format+"\n", a...)
	log.Printf("=> %s", msg)
	_, err := fmt.Print(irc.conn, msg)
	return err
}

// register sends the PASS, NICK, and USER commands.
func (irc Conn) register() error {
	messages := []string{
		// Pass, Nick, then User. This order matters!
		newPassMsg(irc.pass),
		newNickMsg(irc.nick),
		newUserMsg(irc.username, irc.realname),
	}

	for _, m := range messages {
		if err := irc.send(m); err != nil {
			log.Println("Error:", err)
			return err
		}
	}

	return nil
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

	m := NewMessage(s)
	if m.Command == Ping {
		if err := irc.Pong(m.Source); err != nil {
			log.Printf("Error trying to pong: %v", err)
		}
	}

	return m, nil
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
	irc.sendfln("QUIT :why do you hate me")
	return irc.conn.Close()
}

// Connect initiates the IRC protocol with the given credentails.
func Connect(n net.Conn, nick, realname, username, pass string) (*Conn, error) {
	c := &Conn{
		conn:     n,
		reader:   bufio.NewReader(n),
		nick:     nick,
		realname: realname,
		username: username,
		pass:     pass,
	}
	if err := c.register(); err != nil {
		return nil, err
	}
	return c, nil
}
