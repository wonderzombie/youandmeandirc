package main

import (
	"flag"
	"log"
	"net"
	"strings"
	"time"

	irclib "github.com/wonderzombie/youandmeandirc"
	irc "github.com/wonderzombie/youandmeandirc/irc"
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
		log.Fatalln("Unable to create bot:", err)
	}

	addr := net.JoinHostPort(*host, *port)
	timeout, _ := time.ParseDuration("1m")

	// TODO(wonderzombie): Fix this so that IrcConn takes a closure, or something
	// which can generate net.Conn items for it.
	n, _ := net.DialTimeout("tcp", addr, timeout)
	// defer conn.Close()
	c, err := irc.Connect(n, *nick, "...", *username, *pass)
	if err != nil {
		log.Fatalln("Unable to connect:", err)
	}
	bot.Start(c)
}
