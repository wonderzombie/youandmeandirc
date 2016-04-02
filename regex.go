package youandmeandirc

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

func searchReplForRegex(rgx string) (search, repl string, ok bool) {
	// Need three slashes for s/foo/bar/. Don't want to try parsing escape sequences and whatnot.
	if !strings.HasPrefix(rgx, "s/") || !strings.HasSuffix(rgx, "/") {
		return
	}
	reParts := strings.Split(rgx, "/")
	if len(reParts) != 4 {
		return
	}
	search, repl = reParts[1], reParts[2]
	ok = true
	return
}

func (bot *IrcBot) regexListener() (l Listener) {
	l = func(msg IrcMessage) (fired, trap bool) {
		if msg.Command != CmdPrivmsg {
			return
		}

		// Two cases.
		// 1. If the whole thing consists of a simple regex, you can correct what you said.
		// 2. If the first part is a nick, you can "correct" someone else.
		// We'll do #1 for now.

		parts := strings.SplitN(msg.Text, " ", 2)
		head := parts[0]

		search, repl, ok := searchReplForRegex(head)
		if !ok {
			log.Printf("Doesn't look like a regex: %q\n", head)
			return
		}

		re, err := regexp.Compile(search)
		if err != nil {
			log.Printf("Invalid regex %q: %v", head, err)
			return
		}

		// Retrieve the last message we saw from this user and apply it.
		seen, ok := bot.seenList[msg.Nick]
		if !ok {
			log.Printf("User supplied regex but they haven't been seen until now: %v", msg.Nick)
			return
		}

		replaced := re.ReplaceAllString(seen.msg.Text, repl)
		chat := fmt.Sprintf("%v actually meant: %v", msg.Nick, replaced)
		bot.irc.Say(msg.Channel, chat)

		return true, true
	}

	return
}
