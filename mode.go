package goirc

import (
	"fmt"
)

func (c *Bot) ParseMode(modes []string) {
	modeString := modes[1]
	modeArgs := modes[2:]
	var op bool
	for i := 0; i < len(modeString); i++ {
		fmt.Println(i, modeArgs, modeString[i], op)
		switch m := modeString[i]; m {
		case '+':
			op = true
		case '-':
			op = false
		case 'q', 'a', 'o', 'h', 'v':
			mode := string(m)
			nick := modeArgs[0]
			fmt.Println(mode, nick, op)
			modeArgs = modeArgs[1:]
		}

	}
}
