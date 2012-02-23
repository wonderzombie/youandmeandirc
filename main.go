
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
)

const (
	host = "home.zole.org"
	port = "6667"
	user = "trahari"
	nick = "trago"
)

type Message struct {
	Raw string
	Command string // e.g. PRIVMSG
	Server string
	Id string
}

func NewMessage(msg string) (m *Message) {
	m = new(Message)
	m.Raw = msg
	return
}

type UserInfo struct {
	username, nick, pass, realname string
}

//
type IrcConn struct {
	conn net.Conn
	host, port string
	user UserInfo
	reader *bufio.Reader
}

func Connect(c net.Conn, u UserInfo) (irc *IrcConn) {
	irc = new(IrcConn)
	irc.conn = c
	irc.user = u
	irc.reader = bufio.NewReader(irc.conn)
	irc.register()
	return
}

func (irc IrcConn) register() {
	messages := []string {
		irc.newPassMsg(irc.user),
		irc.newNickMsg(irc.user),
		irc.newUserMsg(irc.user),
	}

	for _, m := range messages {
		if _, err := irc.send(m) ; err != nil {
			log.Println("Error:", err)
		}
	}
}

func (irc IrcConn) newPassMsg(user UserInfo) string {
	return fmt.Sprintf("PASS %v", user.pass)
}

func (irc IrcConn) newNickMsg(user UserInfo) string {
	return fmt.Sprintf("NICK %v", user.pass)
}

func (irc IrcConn) newUserMsg(user UserInfo) string {
	return fmt.Sprintf("USER %v fakehost fakeserver :%v", user.username, user.realname)
}

func (irc IrcConn) newJoinMsg(channel string) string {
	return fmt.Sprintf("JOIN %v", channel)
}

func (irc IrcConn) send(cmd string) (int, error) {
	fmt.Println("->", cmd)
	return fmt.Fprintln(irc.conn, cmd)
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
	return strings.Contains(msg, nick)
}

func readLine(r *bufio.Reader) (s string, err error) {
	s, err = r.ReadString('\n')
	return
}

func chatFromMsg(msg string) string {
	i := strings.LastIndex(msg, ":")
	return msg[i:]
}

func main() {
	fmt.Println("hello youandmeandirc")

	var (
		pass string
		channel string
	)
	flag.StringVar(&pass, "passwd", "", "Password for the server.")
	flag.StringVar(&channel, "channel", "", "Channel to join automatically.")
	flag.Parse()

	addr := net.JoinHostPort(host, port)
	timeout, _ := time.ParseDuration("1m")

	conn, _ := net.DialTimeout("tcp", addr, timeout)
	defer conn.Close() // FIXME

	u := UserInfo{username: user, pass: pass, nick: nick, realname: "youandmeandirc"}

	irc := Connect(conn, u)

	/*reader := bufio.NewReader(conn)*/
	/*writePassMsg(conn, pass)*/
	/*writeNickMsg(conn, nick)*/
	/*writeUserMsg(conn, user, "youandmeandirc.go")*/
	joined := false
	for {
		out := irc.Read()
		fmt.Print("<-", out)

		nickResp := fmt.Sprintf("MODE %v", nick)
		if strings.Contains(out, nickResp) && !joined {
			irc.Join(channel)
			joined = true
		} else if strings.Contains(out, "PING :") {
			daemon := strings.Split(out, " ")[1]
			irc.Pong(daemon)
		} else if isChannelMsg(channel, out) {
			chat := chatFromMsg(out)
			if hasMyName(chat) {
				resp := "zzz"
				if strings.Contains(out, "bethday") {
					resp = "bethday sux"
				} else if strings.Contains(out, ":ACTION") {
					resp = "what are you doing"
				}
				irc.SendChat(channel, resp)
			}
		}
	}
}
