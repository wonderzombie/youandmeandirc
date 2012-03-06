package youandmeandirc

import (
	"log"
	"math/rand"
	"strings"
	"time"
)

type IrcBot struct {
	irc       *IrcConn
	listeners []Listener
	connectFn ConnectFn

	channels []string
}

// ConnectFn is used to generate connections.
type ConnectFn func() (*IrcConn, error)

// Listeners are called when a message arrives. The first return value
// indicates whether the message caused the listener to fire. The second return
// value indicates whether this listener requires no other listeners to fire.
type Listener func(msg IrcMessage) (fired, trap bool)

func (bot *IrcBot) runListeners(msg IrcMessage) {
	for _, l := range bot.listeners {
		fired, trap := l(msg)
		if trap && fired {
			return
		}
	}
}

func (bot *IrcBot) pingListener() (pong Listener) {
	pong = func(msg IrcMessage) (fired, trap bool) {
		if msg.Command == "PING" {
			bot.irc.Pong(msg.Origin)
			fired, trap = true, true
		}
		return
	}
	return
}

func (bot *IrcBot) seenListener() (seen Listener) {
	seenList := make(map[string]IrcMessage, 0)

	seen = func(msg IrcMessage) (fired, trap bool) {
		trap = false
		if msg.Command != "PRIVMSG" {
			return
		}
		seenList[msg.Origin] = msg
		fired = true

		return
	}

	return seen
}

func (bot *IrcBot) init() (e error) {
	bot.listeners = []Listener{
		bot.pingListener(),
		bot.onNameListener(),
		bot.seenListener(),
	}
	return nil
}

func (bot *IrcBot) onNameListener() (onName Listener) {
	sayings := []string{
		"What? No",
		"I kissed a boy today.",
	}

	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)
	max := len(sayings)

	onName = func(msg IrcMessage) (fired, trap bool) {
		if !strings.Contains(msg.Text, bot.irc.Nick) {
			return
		}

		fired = true
		trap = true
		choice := rng.Int() % max

		bot.irc.Say(msg.Channel, sayings[choice])
		return
	}
	return
}

// Creates a new bot.
func NewBot() (*IrcBot, error) {
	bot := new(IrcBot)
	if err := bot.init(); err != nil {
		return nil, err
	}
	return bot, nil
}

// Starts a bot running.
func (bot *IrcBot) Start(ircConn *IrcConn) {
	bot.irc = ircConn
	joined := false

	for {
		m, err := bot.irc.Read()
		if err != nil {
			log.Fatalf("Unable to parse message from server: %v", err)
		}

		if !joined && m.Origin == bot.irc.Nick && m.Command == "MODE" {
			bot.irc.Join("#testbot")
			joined = true
		} else {
			bot.runListeners(*m)
		}
	}
}
