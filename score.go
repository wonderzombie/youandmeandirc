package youandmeandirc

import (
  "fmt"
  "log"
  "regexp"
  "time"
)

type Point struct {
  Granter string
  When    time.Time
  Reason  string
}

type Score struct {
  Total  int
  Points []Point
}

var scoreChangeRe = regexp.MustCompile("(\\w+)(\\+\\+|\\-\\-)")
var scoreReqRe = regexp.MustCompile("(\\w+), score?")
var scoreMap = make(map[string]Score, 0)

// TODO: a way to ask for your score.
func (bot *IrcBot) scoreListener() (scorer Listener) {
  scorer = func(msg IrcMessage) (fired, trap bool) {
    if msg.Command != "PRIVMSG" {
      return
    }

    fired, trap = bot.handleScoreChange(msg)
    if fired {
      return
    }

    fired, trap = bot.handleScoreRequest(msg)
    if fired {
      return
    }

    return false, false
  }
  return
}

func (bot *IrcBot) handleScoreChange(msg IrcMessage) (fired, trap bool) {
  scoreChangeMatch := scoreChangeRe.FindStringSubmatch(msg.Text)
  if len(scoreChangeMatch) == 0 {
    return
  }

  nick := scoreChangeMatch[1]
  delta := -1
  if scoreChangeMatch[2] == "++" {
    delta = 1
  }
  granter := msg.Origin
  newPoint := Point{Granter: granter, When: time.Now(), Reason: msg.Text}
  log.Printf("Looks like a score message for %v from %v.\n", nick, granter)

  // TODO: user a pointer instead.
  score := scoreMap[nick]
  score.Total += delta
  score.Points = append(score.Points, newPoint)
  scoreMap[nick] = score

  out := fmt.Sprintf("%v's score is now %d", nick, score.Total)
  bot.irc.Say(msg.Channel, out)

  return true, true
}

func (bot *IrcBot) handleScoreRequest(msg IrcMessage) (fired, trap bool) {
  scoreReqMatch := scoreReqRe.FindStringSubmatch(msg.Text)
  if len(scoreReqMatch) == 0 {
    return
  }

  return
}
