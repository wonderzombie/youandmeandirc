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
type Listener func(IrcMessage) (bool, bool)

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

func (bot *IrcBot) onNameListener() (onName Listener) {
	sayings := []string{
		"I'd love to help, but I need to finish my post on LJ.",
		"Hold that thought, BRB",
		"That's an astute observation.  I would have never thought that!",
		"Sorry, lag.",
		"You know, I try and try, and I'm just never good enough.  Do you ever feel that way?",
		"Can you show me?  Give me a PM.",
		"Are you saying you'll go out with me?",
		"I made some icons of that once and used them in my LJ.",
		"I disagree, but I respect your opinion.",
		"I didn't know you felt that way about me.",
		"Sorry, still catching up with scrollback.",
		"I was thinking the same thing.",
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

func (bot *IrcBot) runListeners(msg IrcMessage) {
	for _, l := range bot.listeners {
		fired, trap := l(msg)
		if trap && fired {
			return
		}
	}
}

func (bot *IrcBot) init() (e error) {
	bot.listeners = []Listener{
		bot.pingListener(),
		bot.scoreListener(),
		bot.seenListener(),
		bot.onNameListener(), // This should go last.
	}
	return nil
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
