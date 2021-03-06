package youandmeandirc

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/wonderzombie/youandmeandirc/irc"
)

type SeenInfo struct {
	Message   irc.Message
	Timestamp time.Time
}

type SeenTrigger struct {
	SeenInfo map[string]SeenInfo
}

func (t SeenTrigger) Id() TriggerId {
	return TriggerId("seen")
}

func (t SeenTrigger) Fire(msg irc.Message, bot *IrcBot, ids []TriggerId) ResultCode {
	fn := bot.seenListener()
	fired, _ := fn(msg)
	if fired {
		return Fired
	}
	return Pass
}

func (bot *IrcBot) seenListener() (seen Listener) {
	bot.seenList = make(map[string]SeenInfo)

	seen = func(msg irc.Message) (fired, trap bool) {
		accepted := []irc.Command{
			irc.Privmsg,
			irc.Join,
			irc.Part,
		}
		if !msg.Matches(accepted) {
			return
		}

		re := regexp.MustCompile(fmt.Sprintf("%v, seen (\\w+)\\?", bot.irc.Nick))

		fired = true
		match := re.FindStringSubmatch(msg.Text)
		if len(match) == 0 {
			info := SeenInfo{msg, time.Now()}
			bot.seenList[msg.Nick] = info
			log.Printf("Storing message from %v: %v\n", msg.Nick, info)
			return
		}

		who := match[1]
		out := fmt.Sprintf("Sorry, haven't seen %v.", who)
		prev, ok := bot.seenList[who]
		if ok {
			out = fmt.Sprintf("I last saw %v at %v, saying \"%v\".", who, prev.Timestamp, prev.Message.Text)
		}
		bot.Say(msg.Channel, out)
		trap = true
		return
	}

	return
}

/// New implementation.

// SeenModule encompasses the Seen lookup, a table containing when IRC nicks were last seen and what they were saying.
type SeenModule struct {
	SeenMap map[string]SeenInfo
}

// Initializes SeenModule.
func (sm *SeenModule) init() {
	sm.SeenMap = make(map[string]SeenInfo)
}

func (sm *SeenModule) accepts() []irc.Command {
	return []irc.Command{
		irc.Privmsg,
		irc.Join,
		irc.Part,
	}
}

// func SeenListener(bot *IrcBot, msg IrcMessage) (trap bool)

// // Public API.
// func (seenMod *SeenModule) HasSeen(nick string) *SeenInfo
