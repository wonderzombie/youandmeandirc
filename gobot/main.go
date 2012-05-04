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
