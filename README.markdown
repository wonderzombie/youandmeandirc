
### youandmeandirc

An extremely rough first attempt at an IRC library implemented in Go.

### install

    $ go install youandmeandirc/gobot
    $ gobot -help

### TODO

* actually implement event listeners/observers/whatever -- mostly done
  * however, consider adding some notion of events. the way the bot does names is kind of bogus.

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
* be quiet
	* based on time rather than # of messages?
* seen 
  * seen enumerated by user -- DONE
	* seen anybody/everybody? -- DONE
	* seen date/time and/or chat -- DONE
* nick++ 
  * give, dock points -- DONE
	* everyone's score -- DONE
	* reasons, incl. adder/demoter -- in progress
* rumor db
	* add rumor
	* promote/demote rumor
	* query, incl. promote/demote
* internet search
  * basic search
	* "continue" (if feasible)
	* last URL (extended history?)
	* recent searches?

* http interface for debugging state?

### example IRC message formats

#### join
:nick!~user@host JOIN :#channel

#### part
:nick!~user@host PART #channel :message content

#### emote
:nick!~user@host PRIVMSG #channel :ACTION message content

