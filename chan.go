package goirc

import (
	"errors"
	"strings"
	"time"
)

func (c *Bot) Join(name string) {
	newchan := &Channel{Name: name, Flood: 0}
	c.Channels = append(c.Channels, newchan)
}

func (c *Bot) Part(name string) {
	for num, pchan := range c.Channels {
		if strings.ToLower(pchan.Name) == strings.ToLower(name) {
			c.Channels = append(c.Channels[:num], c.Channels[num+1:]...)
		}
	}
	if c.RejoinOnKick {
		c.con.PrintLine("JOIN %s\r\n", name)
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

func (c *Channel) FloodControl() int {
	c.Flood += 1
	go c.DecreaseFlood()
	return c.Flood
}

func (c *Channel) DecreaseFlood() {
	time.Sleep((time.Millisecond * 500) + (time.Second * time.Duration(c.Flood)))
	c.Flood = c.Flood - 1
}
