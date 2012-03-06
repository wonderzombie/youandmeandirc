package main

import (
	"flag"
	irclib "github.com/wonderzombie/youandmeandirc"
	"log"
	"net"
	"strings"
	"time"
)

// Bot layer.
func hasMyName(msg string) bool {
	return strings.Contains(msg, *nick)
}

// Flags.
var (
	channel  = flag.String("channel", "#testbot", "Channel to join automatically.")
	nick     = flag.String("nick", "gobot", "Nick to use.")
	pass     = flag.String("pass", "", "Password for the server, if any.")
	username = flag.String("user", "", "Username for identification.")
	host     = flag.String("host", "home.zole.org", "Name of IRC host.")
	port     = flag.String("port", "6667", "Port to connect to on host.")
)

func init() {
	flag.Parse()
}

func defaultMain() {
	addr := net.JoinHostPort(*host, *port)
	timeout, _ := time.ParseDuration("1m")

	irc := &irclib.IrcConn{
		Username: *username,
		Pass:     *pass,
		Nick:     *nick,
		Realname: "...",
	}
	// TODO(wonderzombie): Fix this so that IrcConn takes a closure, or something
	// which can generate net.Conn items for it.
	conn, _ := net.DialTimeout("tcp", addr, timeout)
	defer conn.Close()
	irc.Connect(conn)

	joined := false
	for {
		m, err := irc.Read()
		if err != nil {
			log.Println("Skipping message: ", m.Raw)
			continue
		}

		log.Print("<=", m.Raw)

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

func botMain() {

}

func main() {
	log.Println("hello youandmeandirc")

	bot, err := irclib.NewBot()
	if err != nil {
		log.Fatalf("Unable to create bot:", err)
	}

	addr := net.JoinHostPort(*host, *port)
	timeout, _ := time.ParseDuration("1m")

	irc := &irclib.IrcConn{
		Username: *username,
		Pass:     *pass,
		Nick:     *nick,
		Realname: "...",
	}

	// TODO(wonderzombie): Fix this so that IrcConn takes a closure, or something
	// which can generate net.Conn items for it.
	conn, _ := net.DialTimeout("tcp", addr, timeout)
	defer conn.Close()
	irc.Connect(conn)

	bot.Start(irc)
}
