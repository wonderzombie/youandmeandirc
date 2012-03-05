package youandmeandirc

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// IrcConn represents a connection to an IRC server.
type IrcConn struct {
	Username string
	Pass     string
	Nick     string
	Realname string

	listeners []Listener

	conn   net.Conn
	reader *bufio.Reader

	host, port string
}

/// Network-related.

// send transmits a command to the currently connected server.
func (irc IrcConn) send(s string) (int, error) {
	fmt.Println("=>", s)
	return fmt.Fprintln(irc.conn, s)
}

// register sends the PASS, NICK, and USER commands.
func (irc IrcConn) register() {
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

func (irc IrcConn) pong(daemon string) (int, error) {
	// FIXME: shouldn't this be handled automatically?
	cmd := irc.newPongMsg(daemon)
	return irc.send(cmd)
}

/// Composing various kinds of messages.

func (irc IrcConn) newPassMsg(pass string) string {
	return fmt.Sprintf("PASS %v", pass)
}

func (irc IrcConn) newNickMsg(nick string) string {
	return fmt.Sprintf("NICK %v", nick)
}

func (irc IrcConn) newUserMsg(username, realname string) string {
	return fmt.Sprintf("USER %v * * :%v", username, realname)
}

func (irc IrcConn) newJoinMsg(channel string) string {
	return fmt.Sprintf("JOIN %v", channel)
}

func (irc IrcConn) newPongMsg(daemon string) string {
	return fmt.Sprintf("PONG %v", daemon)
}

/// Public methods.

// Sends a message to a channel.
func (irc IrcConn) Say(channel, chat string) (int, error) {
	msg := fmt.Sprintf("PRIVMSG %v :%v", channel, chat)
	return irc.send(msg)
}

// Joins a given channel.
func (irc IrcConn) Join(channel string) (int, error) {
	cmd := fmt.Sprintf("JOIN %v", channel)
	return irc.send(cmd)
}

// Reads a single message from the server's output.
func (irc IrcConn) Read() (m *IrcMessage, err error) {
	s, err := irc.reader.ReadString('\n')
	if err != nil {
		log.Println("Error:", err)
	}

	m, err = ParseMessage(s)
	if err != nil {
		log.Println("Error:", err)
	}

	return
}

// Connects to IRC with the given connection.
func (irc *IrcConn) Connect(c net.Conn) {
	irc.conn = c
	irc.reader = bufio.NewReader(irc.conn)
	irc.register()
	return
}

// A Listener is called every time a message occurs.
type Listener func(msg string) (fired, trap bool)

// Listen enters a listener into the list of them.
func (irc IrcConn) Listen(l Listener) {
	irc.listeners = append(irc.listeners, l)
}
