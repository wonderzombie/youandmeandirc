package youandmeandirc

import (
	"fmt"
	"log"
	"regexp"
	"time"
)

type SeenInfo struct {
	msg IrcMessage
	t   time.Time
}

func (bot *IrcBot) seenListener() (seen Listener) {
	bot.seenList = make(map[string]SeenInfo)

	seen = func(msg IrcMessage) (fired, trap bool) {
		accepted := []ircCmd{
			CmdPrivmsg,
			CmdJoin,
			CmdPart,
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
			out = fmt.Sprintf("I last saw %v at %v, saying \"%v\".", who, prev.t, prev.msg.Text)
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
	seenMap map[string]SeenInfo
}

// Initializes SeenModule.
func (sm *SeenModule) init() {
	sm.seenMap = make(map[string]SeenInfo)
}

func (sm *SeenModule) accepts() []ircCmd {
	return []ircCmd{
		CmdPrivmsg,
		CmdJoin,
		CmdPart,
	}
}

// func SeenListener(bot *IrcBot, msg IrcMessage) (trap bool)

// // Public API.
// func (seenMod *SeenModule) HasSeen(nick string) *SeenInfo
