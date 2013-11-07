package goirc

import (
	"fmt"
)

// Currently completely ignores channel modes. Probably for the best.
func (c *Bot) ParseMode(modes []string) {
	channel, _ := c.GetChan(string(modes[0]))
	modeString := modes[1]
	modeArgs := modes[2:]
	var op bool
	for i := 0; i < len(modeString); i++ {
		switch m := modeString[i]; m {
		case '+':
			op = true
		case '-':
			op = false
		case 'q', 'a', 'o', 'h':
			nick := modeArgs[0]
			user := channel.GetUser(nick)
			user.Op = op
			modeArgs = modeArgs[1:]
		case 'v':
			nick := modeArgs[0]
			user := channel.GetUser(nick)
			user.Voice = op
			modeArgs = modeArgs[1:]
		}

	}
}
