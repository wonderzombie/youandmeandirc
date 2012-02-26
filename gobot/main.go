
package main

import (
	"bufio"
	"fmt"
	"flag"
	// "io"
	"log"
	"net"
	"strings"
	"time"
	irclib "youandmeandirc"
)


// IrcConn represents a connection to an IRC server.
type IrcConn struct {
	conn net.Conn
	reader *bufio.Reader

	host, port string

	nick string
	pass string
	realname string
	username string
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
		irc.newPassMsg(irc.pass),
		irc.newNickMsg(irc.nick),
		irc.newUserMsg(irc.username, irc.realname),
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

func (irc IrcConn) Read() (m *irclib.IrcMessage, err error) {
	s, err := irc.reader.ReadString('\n')
	if err != nil {
		log.Println("Error:", err)
	}

	m, err = irclib.ParseMessage(s)
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

// Bot layer.
func hasMyName(msg string) bool {
	return strings.Contains(msg, *nick)
}

// Flags.
var (
	channel = flag.String("channel", "#testbot", "Channel to join automatically.")
	nick = flag.String("nick", "gobot", "Nick to use.")
	pass = flag.String("pass", "", "Password for the server, if any.")
	username = flag.String("user", "", "Username for identification.")
	host = flag.String("host", "home.zole.org", "Name of IRC host.")
	port = flag.String("port", "6667", "Port to connect to on host.")
)

func init() {
	flag.Parse()
}

func main() {
	fmt.Println("hello youandmeandirc")

	addr := net.JoinHostPort(*host, *port)
	timeout, _ := time.ParseDuration("1m")

	irc := &IrcConn{
		nick: *nick,
		pass: *pass,
		username: *username,
		realname: "...",
	}
	conn, _ := net.DialTimeout("tcp", addr, timeout)
	defer conn.Close() // FIXME: should the IrcConn take ownership?
	irc.Connect(conn)

	joined := false
	for {
		m, err := irc.Read()
		if err != nil {
			fmt.Println("Skipping message: ", m.Raw)
			continue
		}

		fmt.Print("<=", m.Raw)

		if m.Command == "PING" {
			daemon := m.Params[0]
			irc.Pong(daemon)
			continue
		}

		if !joined && m.Origin == *nick && m.Command == "MODE" {
			// We've finished connecting. Join the channel.
			irc.Join(*channel)
			joined = true
		} else if m.Command == "PRIVMSG" && m.Channel == *channel {
			if hasMyName(m.Text) {
				resp := "zzz"
				if strings.Contains(m.Text, "ACTION") {
					resp = "what are you doing"
				} else if m.Origin == "trapro" {
					resp = "trapro sux"
				}
				irc.Say(*channel, resp)
			}
		}
	}
}
