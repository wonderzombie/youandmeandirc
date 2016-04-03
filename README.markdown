### youandmeandirc

An extremely rough first attempt at an IRC library implemented in Go.

### install

    $ go install youandmeandirc/gobot
    $ gobot -help

### TODO

* actually implement event listeners/observers/whatever -- mostly done
  * however, consider adding some notion of events. the way the bot does names is kind of bogus.

* some notion of listeners for a specific set of messages would also be helpful. there's a lot of boilerplate for each type of listener, where we cancel out for PRIVMSG. A listener could register for certain types of messages and save themselves the trouble of checking that crap.
  * it may even be worthwhile to provide a "shouldFire" function for each listener, so that we have one (optional) set of code which checks whether or not to call it and then the "real" code which can operate under the assumption that our message is valid.
  * alternatively just use this as a design pattern or part of the interface for a module.

* this implies an order of initialization. once we have a healthy connection, *then* initialize stuff. this is because some listeners -may- want to know the bot's nick.
  * different stages of initialization would be overkill. just init all the listeners/modules after we know we've connected to a server, or possibly even as late as channel.

* use channels for reading/writing

* score.go wants to use information from seen.go. this is impossible right now, as all the modules' state is siloed.

* proof of concept: canned responses
	* copy botty's responses -- DONE
	* "no"
* delay in typing rather than instant -- DONE
	* 0.01 - 0.02 per character
* emotes
	* could copy botty's
* go away
	* quit process most likely
* be quiet -- DONE
	* based on time rather than # of messages? -- not done. async messaging not supported (yet?), unfortunately. :(
* seen
  * seen enumerated by user -- DONE
	* seen anybody/everybody? -- DONE
	* seen date/time and/or chat -- DONE
  * track people based on nick changes (e.g. index people by user and/or nick)
* nick++
  * give, dock points -- DONE
	* everyone's score -- DONE
	* reasons, incl. adder/demoter -- DONE
* rumor db
	* add rumor
	* promote/demote rumor
	* query, incl. promote/demote
* internet search
  * basic search
	* "continue" (if feasible)
	* last URL (extended history?)
	* recent searches?
* combat
	* hit points
 	* attacking, incl. crits
  	* healing
   	* death/resurrecting
* uptime -- DONE

### farther afield

* http interface for debugging state?

* text adventure
	* botty already has some stuff like this, such as the fighting functionality.
	* preface commands with >
 	* this could be a bit broader than combat, just some simple verbs like "use" or "examine"
  	* user-defined behavior when it comes to setting what happens when someone uses a command?
   	* or perhaps just a random phrase picked, associated with each command
    	* "what you expected hasn't happened"; "I don't see that here"; "I don't know how to %v"
     	* reaction to swears?

* is calling /whois on anybody interesting in any way?

### example IRC message formats

#### join
`:nick!~user@host JOIN :#channel`

#### part
`:nick!~user@host PART #channel :message content`

#### emote
`:nick!~user@host PRIVMSG #channel :ACTION message content`

### ramblings of questionable value

You can partition the protocol into two major parts, at least when it comes to messages. You have the "chat as application" level and the "chat as protocol" level. Awkward terms, but here's what I mean: PRIVMSG, NOTICE, MODE, and so on are all text. They are also events that a user-agent should display to the user. Numeric codes are "spammy" but only because they're closer to the protocol layer. There's stuff users might look at, sure, but they're not central to the intended use case of IRC: chatting.
