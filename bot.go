package youandmeandirc

import (
	"fmt"
	"log"
	"math/rand"
	// "regexp"
	"strings"
	"time"

	"github.com/wonderzombie/youandmeandirc/irc"
)

// ConnectFn is used to generate connections.
type ConnectFn func() (*irc.Conn, error)

// Listeners are called when a message arrives. The first return value
// indicates whether the message caused the listener to fire. The second return
// value indicates whether this listener requires no other listeners to fire.
type Listener func(irc.Message) (bool, bool)

// TODO: replace Listener with BotListener. This will allow just to just enumerate Listeners
// instead of the rigmarole right now, where methods on IrcBot return Listeners.
type BotListener func(*IrcBot, irc.Message) (bool, bool)

type IrcBot struct {
	irc       *irc.Conn
	listeners []Listener
	connectFn ConnectFn

	channels []string

	namesSet map[string]bool

	uptime time.Time

	seenList map[string]SeenInfo

	asleep bool

	healthList map[string]int

	rng *rand.Rand

	triggers map[TriggerId]Trigger
}

type TriggerId string

// API is like this, roughly:
// Your module implements some function that returns some item that satisfies this interface.
// Id() is used to identify your module to other modules. Therefore anything you export
// on your struct is the API for other modules to interact with yours.
type Trigger interface {
	Id() TriggerId
	Fire(msg irc.Message, bot *IrcBot, ids []TriggerId) bool
}

func (bot *IrcBot) init() (e error) {
	bot.listeners = []Listener{
		bot.sleepListener(), // This must come first.
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

func (bot *IrcBot) Register(t Trigger) {
	id := t.Id()
	bot.triggers[id] = t
}

func (bot *IrcBot) RegisterAll(ts ...Trigger) {
	for _, t := range ts {
		bot.Register(t)
	}
}

func (bot *IrcBot) Trigger(id TriggerId) Trigger {
	t, ok := bot.triggers[id]
	if !ok {
		return nil
	}
	return t
}

func (bot *IrcBot) Random(max int) int {
	return bot.rng.Int() % max
}

type PingTrigger struct{}

func (p *PingTrigger) Id() TriggerId {
	return TriggerId("ping")
}

func (p *PingTrigger) Fire(msg irc.Message, bot *IrcBot, ids []TriggerId) bool {
	if len(ids) > 0 {
		log.Println("Warning: other triggers called before: ", ids)
	}

	if msg.Command == irc.Ping {
		bot.irc.Pong(msg.Source)
		return false
	}

	return true
}

func (bot *IrcBot) pingListener() (pong Listener) {
	pong = func(msg irc.Message) (fired, trap bool) {
		if msg.Command == irc.Ping {
			bot.irc.Pong(msg.Source)
			fired, trap = true, true
		}
		return
	}
	return
}

type NamesTrigger struct {
	NamesSet map[string]bool
}

func (t *NamesTrigger) Id() TriggerId {
	return TriggerId("namelist")
}

func (t *NamesTrigger) Fire(bot *IrcBot, msg irc.Message, ids []TriggerId) bool {
	if msg.Command != irc.Join && msg.Command != irc.Part {
		return false
	}

	if msg.Nick == bot.irc.Nick() {
		return false
	}

	_, ok := t.NamesSet[msg.Source]
	if ok && msg.Command == irc.Part || msg.Command == irc.Quit {
		log.Println("Received part or quit, so removing name:", msg.Nick)
		t.NamesSet[msg.Nick] = false
	}

	//
	if !ok {
		log.Println("Adding nick to list of names:", msg.Nick)
		t.NamesSet[msg.Nick] = true
	}

	return true
}

func (bot *IrcBot) joinListener() (join Listener) {
	join = func(msg irc.Message) (fired, trap bool) {
		if msg.Command != irc.Join {
			return
		}

		_, ok := bot.namesSet[msg.Nick]
		if !ok && msg.Nick != bot.irc.Nick() {
			log.Println("Adding nick to list of names:", msg.Nick)
			// bot.names = append(bot.names, msg.Nick)
			bot.namesSet[msg.Nick] = true
		}

		// This fired, but don't trap it.
		return true, false
	}
	return
}

type MentionMeTrigger struct{}

func (t *MentionMeTrigger) Id() TriggerId {
	return TriggerId("mentionme")
}

func (t *MentionMeTrigger) Fire(msg irc.Message, bot *IrcBot, ids []TriggerId) bool {
	if msg.Command != irc.Privmsg || msg.Source == bot.irc.Nick() {
		return false
	}

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

	mentioned := sortaContains(msg.Text, bot.irc.Nick())
	if !mentioned {
		return false

	}

	choice := bot.Random(len(sayings))
	bot.Say(msg.Channel, sayings[choice])
	return true
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

	name = func(msg irc.Message) (fired, trap bool) {
		if msg.Command != irc.Privmsg || !strings.Contains(msg.Text, bot.irc.Nick()) {
			return
		}

		choice := bot.Random(len(sayings))
		bot.Say(msg.Channel, sayings[choice])
		return true, true
	}
	return
}

func (bot *IrcBot) runTriggers(msg irc.Message) {
	var triggered []TriggerId
	for id, t := range bot.triggers {
		if fired := t.Fire(msg, bot, triggered); fired {
			triggered = append(triggered, id)
		}
	}
}

func (bot *IrcBot) runListeners(msg irc.Message) {
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
	uptime = func(msg irc.Message) (fired, trap bool) {
		if msg.Command != irc.Privmsg {
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

	names = func(msg irc.Message) (fired, trap bool) {
		if code != msg.Code {
			return
		}

		names := strings.Fields(msg.Text)
		ops := "@+"
		for _, name := range names {
			name = strings.Trim(name, ops)
			bot.namesSet[name] = true
		}
		bot.namesSet[bot.irc.Nick()] = true

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
func (bot *IrcBot) Start(c *irc.Conn) {
	bot.irc = c
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
		if !joined && m.Nick == bot.irc.Nick() && m.Command == irc.Mode {
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
