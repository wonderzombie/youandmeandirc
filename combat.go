package youandmeandirc

import (
  "fmt"
  "log"
  "strings"
)

func (bot *IrcBot) combatListener() (combat Listener) {
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

  combat = func(msg IrcMessage) (fired, trap bool) {
    if msg.Command != CmdPrivmsg {
      return
    }

    if !strings.HasPrefix(msg.Text, "\x01ACTION") ||
       !msg.TextHasAny(attacks) {
      return
    }

    fields := strings.Fields(msg.Text)
    // e.g. ":ACTION kicks" isn't a valid attack.
    if len(fields) < 3 {
      return
    }
    log.Printf("Attack received: %q\n", fields)

    // You cannot attack if you're dead.
    attackerHp, ok := bot.healthList[msg.Nick]
    if ok && attackerHp == 0 {
      say := fmt.Sprintf("You can't attack when you're dead, %v!", msg.Nick)
      bot.irc.Say(msg.Channel, say)
      return false, true
    }

    // Is the target present?
    log.Printf("Targets: %q\n", bot.namesSet)
    target := strings.TrimSpace(last(fields))
    ok, _ = bot.namesSet[target]
    if !ok {
      log.Printf("Target is not present: %q\n", target)
      bot.irc.Say(msg.Channel, fmt.Sprintf("%v flails around.", msg.Nick))
      return false, true
    }

    fired, trap = true, true

    health, ok := bot.healthList[target]
    if !ok {
      health = 10
    } else if health == 0 {
      bot.irc.Say(msg.Channel, fmt.Sprintf("%v is already dead!", target))
      return true, true
    }

    var out string
    toHit := 1 + bot.rng.Int() % 6
    damage := 1 + bot.rng.Int() % 10

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