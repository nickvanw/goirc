package goirc

import (
	"fmt"
	"io"
	"net"
)

//TODO implement debug goroutine

const discon string = "DC\r\n"

type Bot struct {
	con          *Conn
	raw          chan *Line
	Write        chan string
	err          chan error
	Message      chan *Message
	Name         string
	Server       string
	port         int
	RejoinOnKick bool
	Channels     []*Channel
	AdminAddr    string
}

type Channel struct {
	Name  string
	Users []*User
	Flood int
}

type User struct {
	Name  string
	Op    bool
	Voice bool
	Admin bool
}

func Create(name string, server string, port int) (*Bot, error) {
	bot := &Bot{
		Name:         name,
		Server:       server,
		port:         port,
		err:          make(chan error),
		raw:          make(chan *Line),
		Message:      make(chan *Message),
		Write:        make(chan string),
		RejoinOnKick: true,
		AdminAddr:    "manacit!~manacit@XVXVqUzRsSStYx",
	}
	_, err := bot.connect()
	if err != nil {
		return nil, err
	}
	go bot.error()
	return bot, nil
}

func (c *Bot) connect() (*Conn, error) {
	conn, err := Connect(c.Server, c.port)
	if err != nil {
		return nil, err
	}
	conn.PrintLine("USER %s 8 * :%s\r\n", c.Name, c.Name)
	conn.PrintLine("NICK %s\r\n", c.Name)
	conn.err = c.err
	c.con = conn
	go c.process()
	return conn, nil
}
func (c *Bot) process() { //Read Middleware
	go c.events()
	for {
		msg, err := c.con.ReadLine()
		if err != nil {
			//For now, we self destruct at any error -- should we change this?
			c.err <- err
			line := &Line{Cmd: discon}
			c.raw <- line
			return
		}
		if len(msg) > 0 {
			m := ParseLine(msg)
			switch m.Cmd {
			case "PING":
				c.con.PrintLine("PONG %s\r\n", m.Args[0])
			default:
				c.raw <- m
			}
		}
	}
}
func isTimeout(err error) bool {
	neterr, ok := err.(net.Error)
	return ok && neterr.Timeout()
}
func (c *Bot) error() {
	for {
		err := <-c.err
		fmt.Println("Error Detected: ", err)
		if err == io.EOF || isTimeout(err) {
			c.con.Disconnect()
			conn, err := c.connect()
			if err != nil {
				fmt.Println(err)
				panic("Too Many Retries?")
			} else {
				c.con = conn
			}
		}
	}
}
