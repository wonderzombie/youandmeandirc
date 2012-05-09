package youandmeandirc

import (
  "fmt"
  "strings"
)

func (bot *IrcBot) combatListener() (combat Listener) {
  combat = func(msg IrcMessage) (fired, trap bool) {
    if msg.Command != CmdPrivmsg {
      return
    }

    if strings.HasPrefix(msg.Text, "\x01ACTION") {
      bot.irc.Say(msg.Channel, "wtf")
      return true, true
    }

    attacks := []string{
      "beat",
      "gouges",
      "hits",
      "kicks",
      "pummels",
      "punches",
      "slap",
      "smack",
      "stabs",
    }

    if !msg.TextHasAny(attacks) {
      return
    }

    fields := strings.Fields(msg.Text)
    // e.g. ":ACTION kicks" isn't a valid attack.
    if len(fields) < 3 {
      return
    }

    // You cannot attack if you're dead.
    attackerHp, ok := bot.healthList[msg.Nick]
    if ok && attackerHp == 0 {
      say := fmt.Sprintf("You can't attack when you're dead, %v!", msg.Nick)
      bot.irc.Say(msg.Channel, say)
      return
    }

    // Is the target present?
    target := fields[len(fields) - 1]
    _, ok = bot.namesSet[target]
    if !ok {
      bot.irc.Say(msg.Channel, fmt.Sprintf("%v flails around.", msg.Nick))
    }

    fired, trap = true, true

    health, ok := bot.healthList[target]
    if !ok {
      health = 10
    }

    var out string
    toHit := bot.rng.Int() % 6
    damage := bot.rng.Int()

    switch toHit {
    case 1:
      out = fmt.Sprintf("%v misses %v!", msg.Nick, target)
    case 6:
      damage *= 2
      out = fmt.Sprintf("%v crits %v for %v damage!", msg.Nick, target, damage)
    default:
      out = fmt.Sprintf("%v hits %v for %v damage!", msg.Nick, target, damage)
    }

    health -= damage
    bot.irc.Say(msg.Channel, out)

    if health <= 0 {
      out = fmt.Sprintf("%v has died!", target)
      bot.irc.Say(msg.Channel, out)
      health = 0
    }

    bot.healthList[target] = health

    return
  }
  return
}