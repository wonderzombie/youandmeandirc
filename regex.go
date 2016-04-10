package youandmeandirc

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/wonderzombie/youandmeandirc/irc"
)

type RegexTrigger struct{}

func (t RegexTrigger) Id() TriggerId {
	return TriggerId("regex")
}

func (t RegexTrigger) Fire(msg irc.Message, bot *IrcBot, ids []TriggerId) ResultCode {
	fn := bot.combatListener()
	fired, _ := fn(msg)
	if fired {
		return Fired
	}
	return Pass
}

type Replacement struct {
	search  string
	replace string
}

func regex(msg string) *Replacement {
	// This means we don't support spaces in regexen.
	words := strings.Fields(msg)
	for _, word := range words {
		if strings.HasPrefix(word, "s/") && strings.HasSuffix(word, "/") {
			// The format:
			//   s/foo/bar/
			//    ^   ^   ^
			parts := strings.Split(word, "/")
			return &Replacement{
				search:  parts[1],
				replace: parts[2],
			}
		}
	}
	return nil
}
}

func (bot *IrcBot) regexListener() (l Listener) {
	l = func(msg irc.Message) (fired, trap bool) {
		if msg.Command != irc.Privmsg {
			return
		}

		// Two cases.
		// 1. If the whole thing consists of a simple regex, you can correct what you said.
		// 2. If the first part is a nick, you can "correct" someone else.
		// We'll do #1 for now.

		parts := strings.SplitN(msg.Text, " ", 2)
		head := parts[0]

		res := regex(head)
		if res == nil {
			log.Printf("Doesn't look like a regex: %q\n", head)
			return
		}

		// Retrieve the last message we saw from this user and apply it.
		seen, ok := bot.seenList[msg.Nick]
		if !ok {
			log.Printf("User supplied regex but they haven't been seen until now: %v", msg.Nick)
			return
		}

		re, err := regexp.Compile(res.search)
		if err != nil {
			log.Printf("Invalid regex %q: %v", head, err)
			return
		}

		replaced := re.ReplaceAllString(seen.Message.Text, res.replace)
		chat := fmt.Sprintf("%v actually meant: %v", msg.Nick, replaced)
		bot.irc.Say(msg.Channel, chat)

		return true, true
	}

	return
}
