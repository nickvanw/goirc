package goirc

import (
	"errors"
	"strings"
)

func (c *Bot) Join(name string) {
	newchan := &Channel{Name: name}
	c.Channels = append(c.Channels, newchan)
}

func (c *Bot) Part(name string) {
	for num, pchan := range c.Channels {
		if strings.ToLower(pchan.Name) == strings.ToLower(name) {
			c.Channels = append(c.Channels[:num], c.Channels[num+1:]...)
		}
	}
}

func (c *Bot) Names(args []string) {
	addChannel, err := c.GetChan(args[2])
	if err != nil {
		return
	}
	names := strings.Split(args[3], " ")
	for _, name := range names {
		addChannel.UserJoin(name)
	}
}

func (c *Bot) GetChan(channel string) (*Channel, error) {
	for _, chan_try := range c.Channels {
		if strings.ToLower(chan_try.Name) == strings.ToLower(channel) {
			return chan_try, nil
		}
	}
	return &Channel{}, errors.New("Not a channel!")
}

func (c *Bot) GetChanList() string {
	list := make([]string, 0)
	for _, IRCChan := range c.Channels {
		list = append(list, IRCChan.Name)
	}
	return strings.Join(list, ",")
}
