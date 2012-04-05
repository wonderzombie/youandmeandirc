package youandmeandirc

import (
  "fmt"
  "log"
  "regexp"
  "time"
)

// TODO: a way to ask for your score.
func (bot *IrcBot) scoreListener() (scorer Listener) {
  type Point struct {
    Granter string
    When    time.Time
    Reason  string
  }

  type Score struct {
    Total  int
    Points []Point
  }

  scoreMap := make(map[string]Score, 0)

  scorer = func(msg IrcMessage) (fired, trap bool) {
    if msg.Command != "PRIVMSG" {
      return
    }

    re := regexp.MustCompile("(\\w+)(\\+\\+|\\-\\-)")
    match := re.FindStringSubmatch(msg.Text)
    if len(match) == 0 {
      log.Printf("Is not a score message, skipping.")
      return
    }
    log.Printf("Is possibly a score message.")

    nick := match[1]
    delta := -1
    if match[2] == "++" {
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
  return
}
