package youandmeandirc

import (
	"fmt"
	"log"
	"math/rand"
	// "regexp"
	"strings"
	"time"
)

// ConnectFn is used to generate connections.
type ConnectFn func() (*IrcConn, error)

// Listeners are called when a message arrives. The first return value
// indicates whether the message caused the listener to fire. The second return
// value indicates whether this listener requires no other listeners to fire.
type Listener func(IrcMessage) (bool, bool)

// TODO: replace Listener with BotListener. This will allow just to just enumerate Listeners
// instead of the rigmarole right now, where methods on IrcBot return Listeners.
type BotListener func(*IrcBot, IrcMessage) (bool, bool)

type IrcBot struct {
	irc       *IrcConn
	listeners []Listener
	connectFn ConnectFn

	channels []string
	names    []string

	namesSet map[string]bool

	uptime time.Time

	seenList map[string]SeenInfo

	asleep   bool

	healthList map[string]int

	rng *rand.Rand
}

func (bot *IrcBot) init() (e error) {
	bot.listeners = []Listener{
		bot.pingListener(),  // This must come first.
		bot.sleepListener(), // This must come second.
		bot.joinListener(),
		bot.namesListener(),
		bot.regexListener(),
		bot.scoreListener(),
		bot.seenListener(),
		bot.combatListener(),
		bot.uptimeListener(),
		bot.onNameListener(), // This should go last.
	}
	bot.namesSet = make(map[string]bool)
	bot.healthList = make(map[string]int)
	return nil
}

func (bot *IrcBot) pingListener() (pong Listener) {
	pong = func(msg IrcMessage) (fired, trap bool) {
		if msg.Command == CmdPing {
			bot.irc.Pong(msg.Origin)
			fired, trap = true, true
		}
		return
	}
	return
}

func (bot *IrcBot) joinListener() (join Listener) {
	join = func(msg IrcMessage) (fired, trap bool) {
		if msg.Command != CmdJoin {
			return
		}

		_, ok := bot.namesSet[msg.Nick]
		if !ok && msg.Nick != bot.irc.Nick {
			log.Println("Adding nick to list of names:", msg.Nick)
			// bot.names = append(bot.names, msg.Nick)
			bot.namesSet[msg.Nick] = true
		}

		// This fired, but don't trap it.
		return true, false
	}
	return
}

func (bot *IrcBot) onNameListener() (name Listener) {
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

	// src := rand.NewSource(time.Now().UnixNano())
	// rng := rand.New(src)
	max := len(sayings)

	name = func(msg IrcMessage) (fired, trap bool) {
		if msg.Command != CmdPrivmsg || !strings.Contains(msg.Text, bot.irc.Nick) {
			return
		}

		choice := bot.rng.Int() % max
		bot.Say(msg.Channel, sayings[choice])
		return true, true
	}
	return
}

func (bot *IrcBot) runListeners(msg IrcMessage) {
	for _, l := range bot.listeners {
		// TODO: simplify this. We probably only need one and that'd be trap.
		if _, trap := l(msg); trap {
			return
		}
	}
}

func (bot *IrcBot) askForNames() {
	for _, channel := range bot.channels {
		bot.irc.Names(channel)
	}
}

func (bot *IrcBot) uptimeListener() (uptime Listener) {
	uptime = func(msg IrcMessage) (fired, trap bool) {
		if msg.Command != CmdPrivmsg {
			return
		}

		lower := strings.ToLower(msg.Text)
		expected := fmt.Sprintf("%v, uptime?", bot.irc.Nick)
		if lower != expected {
			return
		}

		bot.irc.Say(msg.Channel, fmt.Sprintf("Uptime is %v", time.Since(bot.uptime)))
		return true, true
	}
	return
}

func (bot *IrcBot) namesListener() (names Listener) {
	// This can actually be multiple lines. The termination line that you want is 366.
	code := "353"

	names = func(msg IrcMessage) (fired, trap bool) {
		if code != msg.Code {
			return
		}

		names := strings.Fields(msg.Text)
		ops := "@+"
		for _, name := range names {
			name = strings.Trim(name, ops)
			bot.namesSet[name] = true
		}
		bot.namesSet[bot.irc.Nick] = true

		log.Println("Names are now:", bot.namesSet)

		return true, false
	}
	return
}

// Wrapper around IrcConn.Say which simulates typing.
func (bot *IrcBot) Say(channel, out string) {
	ms := 10 * len(out)
	// Pretend we're typing.
	time.Sleep(time.Duration(ms) * time.Millisecond)
	bot.irc.Say(channel, out)
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

	// Initialize RNG.
	src := rand.NewSource(time.Now().UnixNano())
	bot.rng = rand.New(src)

	for {
		m, err := bot.irc.Read()
		// fmt.Printf("%+v\n", *m)
		if err != nil {
			log.Fatalf("Unable to parse message from server: %v", err)
		}

		// TODO: uh, look at the actual codes so we know when we've joined. This is a bit hacky.
		if !joined && m.Nick == bot.irc.Nick && m.Command == CmdMode {
			bot.irc.Join("#testbot")
			// Sorta dumb, but basically don't count uptime until we've joined a channel.
			bot.uptime = time.Now()
			joined = true
			// Collect a list of names.
			bot.askForNames()
		} else {
			bot.runListeners(*m)
		}
	}
}
