package youandmeandirc

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

// IrcConn represents a connection to an IRC server.
type IrcConn struct {
	conn net.Conn
	reader *bufio.Reader

	host, port string

	Nick string
	Pass string
	Realname string
	Username string
}

// Constants which describe the format of various client-to-server commands.
const (
	JOIN = iota
	NICK
	PASS
	PONG
	PRIVMSG
	USER
)

var CommandMap = []string{
	"%v",
	"%v",
	"%v",
	"%v",
	"%v :%v",
	"%v * * :%v",
}

func (irc *IrcConn) commandForType(t int, args ...string) string {
	format := CommandMap[t]
	return fmt.Sprintf(format, args)
}

/// Network-related.

func (irc *IrcConn) Connect(c net.Conn) {
	irc.conn = c
	irc.reader = bufio.NewReader(irc.conn)
	irc.register()
	return
}

// send transmits a command to the currently connected server.
func (irc IrcConn) send(s string) (int, error) {
	fmt.Println("=>", s)
	return fmt.Fprintln(irc.conn, s)
}

// register sends the PASS, NICK, and USER commands.
func (irc IrcConn) register() {
	messages := []string {
		irc.newPassMsg(irc.Pass),
		irc.newNickMsg(irc.Nick),
		irc.newUserMsg(irc.Username, irc.Realname),
	}

	for _, m := range messages {
		if _, err := irc.send(m) ; err != nil {
			log.Println("Error:", err)
		}
	}
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

func (irc IrcConn) Say(channel, chat string) (int, error) {
	msg := fmt.Sprintf("PRIVMSG %v :%v", channel, chat)
	return irc.send(msg)	
}

func (irc IrcConn) Join(channel string) (int, error) {
	cmd := fmt.Sprintf("JOIN %v", channel)
	return irc.send(cmd)
}

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

func (irc IrcConn) Pong(daemon string) (int, error) {
	// FIXME: shouldn't this be handled automatically?
	cmd := irc.newPongMsg(daemon)
	return irc.send(cmd)
}
