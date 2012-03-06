package youandmeandirc

import (
	"log"
)

type IrcBot struct {
	conn      *IrcConn
	listeners []Listener
	connectFn ConnectFn

	channels []string
}

// ConnectFn is used to generate connections.
type ConnectFn func() *IrcConn

// Listeners are called when a message arrives. The first return value
// indicates whether the message caused the listener to fire. The second return
// value indicates whether this listener requires no other listeners to fire.
type Listener func(msg IrcMessage) (fired, trap bool)

func (bot *IrcBot) Start(fn ConnectFn) {
	bot.connectFn = fn
	bot.conn = bot.connectFn()

	// Join channels.
	for _, ch := range bot.channels {
		bot.conn.Join(ch)
	}

	for {
		m, err := bot.conn.Read()

		if err != nil {
			log.Println("Skipping message: ", m.Raw)
			continue
		}

		log.Println("<=", m.Raw)

		bot.runListeners(*m)
	}
}

func (bot *IrcBot) registerListener(l Listener) {
	bot.listeners = append(bot.listeners, l)
}

func (bot *IrcBot) runListeners(msg IrcMessage) {
	for _, l := range bot.listeners {
		fired, trap := l(msg)
		if trap && fired {
			return
		}
	}
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
	bot.channels = []string{"#testbot"}
	return
}
