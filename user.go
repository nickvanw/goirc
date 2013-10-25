package goirc

import (
	"errors"
	"strings"
)

func (c *Bot) JoinUser(channel string, name string) {
	addChannel, err := c.GetChan(channel)
	if err != nil {
		return
	}
	addChannel.UserJoin(name)
}

func (c *Bot) PartUser(channel string, name string) {
	remChannel, err := c.GetChan(channel)
	if err != nil {
		return
	}
	remChannel.UserDel(name)
}

func (c *Bot) QuitUser(name string) {
	for _, channel := range c.Channels {
		channel.UserDel(name)
	}
}

func (c *Channel) UserJoin(name string) {
	isop := false
	isvoice := false
	switch c := name[0]; c {
	case '~', '&', '@', '%', '+':
		switch op := name[0]; op {
		case '~', '&', '@', '%':
			isop = true
		case '+':
			isvoice = true
		}
		name = name[1:]
	}
	if c.InChan(name) != true {
		c.AddUser(name)
	}
	user := c.GetUser(name)
	user.Op = isop
	user.Voice = isvoice
}

func (c *Channel) UserDel(name string) {
	for num, user := range c.Users {
		if strings.ToLower(user.Name) == strings.ToLower(name) {
			c.Users = append(c.Users[:num], c.Users[num+1:]...)
		}
	}
}

func (c *Bot) GetUser(channel string, user string) (*User, error) {
	getChan, err := c.GetChan(channel)
	if err != nil {
		return &User{}, err
	}
	for _, usertry := range getChan.Users {
		if strings.ToLower(usertry.Name) == strings.ToLower(user) {
			return usertry, nil
		}
	}
	return &User{}, errors.New("Not a user!")
}

func (c *Channel) GetUser(name string) *User {
	for _, user := range c.Users {
		if strings.ToLower(user.Name) == strings.ToLower(name) {
			return user
		}
	}
	return &User{}
}

func (c *Channel) InChan(name string) bool {
	for _, user := range c.Users {
		if strings.ToLower(user.Name) == strings.ToLower(name) {
			return true
		}
	}
	return false
}

func (c *Channel) AddUser(name string) {
	newUser := &User{Name: name}
	c.Users = append(c.Users, newUser)
}
