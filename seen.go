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
  seenList := make(map[string]SeenInfo, 0)

  seen = func(msg IrcMessage) (fired, trap bool) {
    if msg.Command != "PRIVMSG" {
      return
    }

    fired = true
    re := regexp.MustCompile(fmt.Sprintf("%s, seen (\\w+)\\?", bot.irc.Nick))
    match := re.FindStringSubmatch(msg.Text)
    if len(match) == 0 {
      info := SeenInfo{msg, time.Now()}
      seenList[msg.Origin] = info
      log.Printf("Storing message from %v: %v\n", msg.Origin, info)
      return
    }

    who := match[1]
    out := fmt.Sprintf("Sorry, haven't seen %v.", who)
    prev, ok := seenList[who]
    if ok {
      out = fmt.Sprintf("I last saw %v at %v, saying \"%v\".",
        who, prev.t, prev.msg.Text)
    }
    bot.irc.Say(msg.Channel, out)
    trap = true
    return
  }

  return seen
}