package youandmeandirc

import (
  // "fmt"
  "log"
  "time"
)

func (bot *IrcBot) sleepListener() (sleep Listener) {
  sleepMinutes := time.Duration(5) * time.Minute
  var sleptAt time.Time

  sleep = func(msg IrcMessage) (fired, trap bool) {
    if msg.Command != CmdPrivmsg {
      return
    }

    if bot.asleep {
      if msg.HasText("wake up") && msg.HasText(bot.irc.Nick) {
        // wake up
        bot.asleep = false
        bot.irc.Say(msg.Channel, "I'm awake! I'm awake!")
      } else {
        since := time.Since(sleptAt)
        if since.Minutes() > sleepMinutes.Minutes() {
          // unsleep!
          bot.irc.Say(msg.Channel, "Zzz— what? How long was I out?")
          bot.asleep = false
        } else {
          log.Printf("Zzzz. Still sleeping. %v minutes to go.\n", sleepMinutes.Minutes()-since.Minutes())
        }
      }

      return true, true
    }

    if msg.Command != CmdPrivmsg || !msg.HasText(bot.irc.Nick) {
      return
    }

    valid := []string{
      "shut up",
      "hush",
      "pipe down",
      "be quiet",
      "silence",
    }

    if !msg.TextHasAny(valid) {
      return
    }

    // We've been told to sleep.
    bot.asleep = true
    bot.irc.Say(msg.Channel, "OK, I'll go to sleep. Good night.")
    sleptAt = time.Now()
    return true, true
  }
  return
}
