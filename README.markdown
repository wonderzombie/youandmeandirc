
### youandmeandirc

An extremely rough first attempt at an IRC library implemented in Go.

### install

    $ go install youandmeandirc/gobot
    $ gobot -help

### TODO

* actually implement event listeners/observers/whatever

* proof of concept: canned responses
	* copy botty's responses
	* "no"
* delay in typing rather than instant
	* 0.01 - 0.02 per character
* emotes
	* could copy botty's
* go away
	* quit process most likely
* be quiet
	* based on time rather than # of messages?
* seen
  * seen enumerated by user
	* seen anybody/everybody?
	* seen date/time and/or chat
* nick++
  * ++ --
	* everyone's score
	* reasons, incl. adder/demoter
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

