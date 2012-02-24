
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
	passwd string
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
func (irc IrcConn) send(cmd string) (int, error) {
	fmt.Println("=>", cmd)
	return fmt.Fprintln(irc.conn, cmd)
}

// register sends the PASS, NICK, and USER commands.
func (irc IrcConn) register() {
	messages := []string {
		irc.newPassMsg(irc.passwd),
		irc.newNickMsg(irc.nick),
		irc.newUserMsg(irc.username, irc.realname),
	}
	fmt.Println(messages)

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

func (irc IrcConn) SendChat(channel, chat string) (int, error) {
	msg := fmt.Sprintf("PRIVMSG %v :%v", channel, chat)
	return irc.send(msg)	
}

func (irc IrcConn) Join(channel string) (int, error) {
	cmd := fmt.Sprintf("JOIN %v", channel)
	return irc.send(cmd)	
}

func (irc IrcConn) Read() (s string) {
	s, err := irc.reader.ReadString('\n')
	if err != nil {
		log.Println("Error:", err)
	}
	return
}

// FIXME: shouldn't this be handled automatically?
func (irc IrcConn) Pong(daemon string) (int, error) {
	/*bits := strings.Split(msg, " ")*/
	/*if len(bits) < 2 {*/
		/*log.Println("Invalid PONG message:", msg)*/
		/*return 0, *new(error)*/
	/*}*/

	cmd := irc.newPongMsg(daemon)
	return irc.send(cmd)
}

/// Misc methods.
func isChannelMsg(channel, msg string) bool {
	expected := fmt.Sprintf("PRIVMSG %v", channel)
	if !strings.Contains(msg, expected) { 
		return false 
	} 
	return true
}

func hasMyName(msg string) bool {
	return strings.Contains(msg, *nick)
}

func readLine(r *bufio.Reader) (s string, err error) {
	s, err = r.ReadString('\n')
	return
}

func chatFromMsg(msg string) string {
	i := strings.LastIndex(msg, ":")
	return msg[i:]
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
		passwd: *pass,
		username: *username,
		realname: "...",
	}
	fmt.Println(irc)
	conn, _ := net.DialTimeout("tcp", addr, timeout)
	defer conn.Close() // FIXME: should the IrcConn take ownership?
	irc.Connect(conn)

	joined := false
	for {
		out := irc.Read()

		m, _ := irclib.ParseMessage(out)
		fmt.Print("<=", m.Raw)

		nickResp := fmt.Sprintf("MODE %v", *nick)
		if strings.Contains(out, nickResp) && !joined {
			irc.Join(*channel)
			joined = true
		} else if strings.Contains(out, "PING :") {
			daemon := strings.Split(out, " ")[1]
			irc.Pong(daemon)
		} else if isChannelMsg(*channel, out) {
			chat := chatFromMsg(out)
			if hasMyName(chat) {
				resp := "zzz"
				if strings.Contains(out, "bethday") {
					resp = "bethday sux"
				} else if strings.Contains(out, ":ACTION") {
					resp = "what are you doing"
				}
				irc.SendChat(*channel, resp)
			}
		}
	}
}
