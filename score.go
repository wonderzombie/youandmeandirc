package youandmeandirc

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/wonderzombie/youandmeandirc/irc"
)

type ScoreTrigger struct{}

func (t ScoreTrigger) Id() TriggerId {
	return TriggerId("score")
}

func (t ScoreTrigger) Fire(msg irc.Message, bot *IrcBot, ids []TriggerId) ResultCode {
	fn := bot.scoreListener()
	fired, _ := fn(msg)
	if fired {
		return Fired
	}
	return Pass
}

type Point struct {
	Granter string
	When    time.Time
	// TODO: Reason doesn't really do anything now. Need access to stuff in seen.go to fix this.
	Reason   string
	Increase bool
}

type Score struct {
	Total  int
	Points []Point
}

var scoreChangeRe = regexp.MustCompile("(\\w+)(\\+\\+|\\-\\-)")
var myScoreRe = regexp.MustCompile("(\\w+), my score\\?")
var scoreListRe = regexp.MustCompile("(\\w+), scores?\\?")
var scoreMap = make(map[string]Score, 0)

func (bot *IrcBot) scoreListener() (scorer Listener) {
	scorer = func(msg irc.Message) (fired, trap bool) {
		if msg.Command != irc.Privmsg {
			return
		}

		// TODO: clean this up
		fired, trap = bot.handleScoreChange(msg)
		if fired {
			return
		}

		fired, trap = bot.handleScoreRequest(msg)
		if fired {
			return
		}

		fired, trap = bot.handleMyScoreRequest(msg)
		if fired {
			return
		}

		return false, false
	}
	return
}

func (bot *IrcBot) handleScoreChange(msg irc.Message) (bool, bool) {
	scoreChangeMatch := scoreChangeRe.FindStringSubmatch(msg.Text)
	if len(scoreChangeMatch) == 0 {
		return false, false
	}

	nick := scoreChangeMatch[1]
	_, ok := bot.namesSet[nick]

	if !ok {
		log.Println("Skipping because this isn't a nick for someone present:", nick)
		return false, false
	}

	delta := -1
	if scoreChangeMatch[2] == "++" {
		delta = 1
	}
	granter := msg.Nick
	// TODO: make Reason interesting. If the originator only said foo++ then this is silly.
	// Instead, we should use the SeenList in seen.go to gather what the person last said. If
	// that person has no entry in the list, then omit a reason and just grant them the point.
	reason := msg.Text
	when := time.Now()
	if seenInfo, ok := bot.seenList[granter]; ok {
		reason = seenInfo.Message.Text
		when = seenInfo.Timestamp
	}
	newPoint := Point{Granter: granter, When: when, Reason: reason}
	newPoint.Increase = delta == 1
	log.Printf("Looks like a score message for %v from %v.\n", nick, granter)

	// TODO: user a pointer instead. Specifically, we should be able to get score out, modify it,
	// and not have to reassign it at the end.
	score := scoreMap[nick]
	score.Total += delta
	score.Points = append(score.Points, newPoint)
	scoreMap[nick] = score

	out := fmt.Sprintf("%v's score is now %d", nick, score.Total)
	bot.Say(msg.Channel, out)

	return true, true
}

func (bot *IrcBot) handleScoreRequest(msg irc.Message) (fired, trap bool) {
	scoreReqMatch := scoreListRe.FindStringSubmatch(msg.Text)
	if len(scoreReqMatch) == 0 {
		return
	}

	if len(scoreMap) == 0 {
		bot.Say(msg.Channel, "Nobody has a score yet!")
		return true, true
	}

	for nick := range bot.namesSet {
		out := fmt.Sprintf("%v has no score.", nick)
		if score, ok := scoreMap[nick]; ok {
			if nick == bot.irc.Nick() {
				out = fmt.Sprintf("My score is %v.", score.Total)
			} else {
				out = fmt.Sprintf("%v's score is %v.", nick, score.Total)
			}
		}
		bot.Say(msg.Channel, out)
	}

	return true, true
}

func (bot *IrcBot) handleMyScoreRequest(msg irc.Message) (fired, trap bool) {
	myScoreMatch := myScoreRe.FindStringSubmatch(msg.Text)
	if len(myScoreMatch) == 0 {
		return
	}

	out := []string{fmt.Sprintf("%v, you don't have a score yet.", msg.Nick)}
	if score, ok := scoreMap[msg.Nick]; ok {
		out = []string{fmt.Sprintf("%v, your score is %v.", msg.Nick, score.Total)}
		for _, point := range score.Points {
			verb := "docked"
			if point.Increase {
				verb = "gave"
			}
			s := fmt.Sprintf("%v %v you a point at %v for saying \"%v\"", point.Granter, verb, point.When, point.Reason)
			out = append(out, s)
		}
	}

	for _, chat := range out {
		bot.Say(msg.Channel, chat)
	}

	return true, true
}
