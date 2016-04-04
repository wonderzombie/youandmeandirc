package youandmeandirc

import (
	// "fmt"
	"log"
	"time"

	"github.com/wonderzombie/youandmeandirc/irc"
)

func (bot *IrcBot) sleepListener() (sleep Listener) {
	sleepMinutes := time.Duration(5) * time.Minute
	var sleptAt time.Time

	sleep = func(msg irc.Message) (fired, trap bool) {
		if msg.Command != irc.Privmsg {
			return
		}

		if bot.asleep {
			if msg.TextHas("wake up") && msg.TextHas(bot.irc.Nick()) {
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

		if msg.Command != irc.Privmsg || !msg.TextHas(bot.irc.Nick()) {
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
