package youandmeandirc

import (
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

    return
  }
  return
}