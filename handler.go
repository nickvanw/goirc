package goirc

import (
	"fmt"
	"strings"
)

type Message struct {
	Nick, Chan string
	Message    []string
	ChanObj    *Channel
	Type       int
	out        chan string
	Host       string
	Ident      string
	Admin      bool
}

//TODO: Create structs for common messages.
//TODO keep track of bot own name for PRIVMSG purposes
func (c *Bot) events() { //Where all the reading magic happens
	for {
		select {
		case msg := <-c.raw:
			//Now we have the message event, it's time we handle it?
			switch msg.Cmd {
			case "PRIVMSG":
				//Handle Private Messages
				//Fix the to/from issue either here or in the output message
				admin := CheckAdmin(msg.Nick, msg.Ident, msg.Host, c.AdminAddr)
				ch, _ := c.GetChan(msg.Args[0])
				out := &Message{
					Nick:    msg.Nick,
					Chan:    msg.Args[0],
					ChanObj: ch,
					Message: strings.Split(msg.Args[1], " "),
					out:     c.Write,
					Host:    msg.Host,
					Ident:   msg.Ident,
					Admin:   admin,
				}
				if msg.Args[0] == c.Name {
					out.Type = 1
					out.Chan = ""
				} else {
					out.Type = 0
				}
				c.Message <- out
			case "001":
				if len(c.Channels) > 0 {
					OldChannels := c.GetChanList()
					c.Channels = []*Channel{}
					c.con.PrintLine("JOIN %s\r\n", OldChannels)
				}
			case "NOTICE":
				admin := CheckAdmin(msg.Nick, msg.Ident, msg.Host, c.AdminAddr)
				out := &Message{
					Nick:    msg.Nick,
					Message: strings.Split(msg.Args[1], " "),
					Chan:    "",
					ChanObj: &Channel{},
					Type:    2,
					out:     c.Write,
					Host:    msg.Host,
					Ident:   msg.Ident,
					Admin:   admin,
				}
				c.Message <- out
			case "NICK":
				newname := msg.Args[0]
				if msg.Nick == c.Name {
					c.Name = newname
				} else {
					c.RenameUser(msg.Nick, newname)
				}
			case "JOIN":
				if msg.Nick == c.Name {
					c.Join(msg.Args[0])
				} else {
					c.JoinUser(msg.Args[0], msg.Nick)
				}
			case "PART":
				if msg.Nick == c.Name {
					c.Part(msg.Args[0])
				} else {
					c.PartUser(msg.Args[0], msg.Nick)
				}
			case "KICK":
				if msg.Args[1] == c.Name {
					c.Part(msg.Args[0])
				} else {
					c.PartUser(msg.Args[0], msg.Args[1])
				}
			case "QUIT":
				c.QuitUser(msg.Nick)
			case "353":
				c.Names(msg.Args)
			case "MODE":
				c.ParseMode(msg.Args)
			case discon:
				fmt.Println("Disconnected: Ending goworker()")
				return
			default:
				//fmt.Println(msg.Raw)
			}
		case out := <-c.Write:
			c.con.WriteLine(out)
		}
	}
}

func split(s string, N int) []string {
	var r []string
	if len(s) <= N {
		return append(r, s)
	}
	for len(s) > N {
		r, s = append(r, s[:N]), s[N:]
	}
	return append(r, s)
}

func (msg *Message) Return(out string) {
	d := split(out, 400)
	for i := 0; i < len(d) && i <= 1; i++ {
		out := d[i]
		switch msg.Type {
		case 0:
			msg.out <- fmt.Sprintf("PRIVMSG %s :%s", msg.Chan, out)
		case 1:
			msg.out <- fmt.Sprintf("PRIVMSG %s :%s", msg.Nick, out)
		case 2:
			msg.out <- fmt.Sprintf("NOTICE %s :%s", msg.Nick, out)
		default:
			fmt.Println("Not Implemented!")
		}
	}
}

func (msg *Message) Send(to string, message string) {
	msg.out <- fmt.Sprintf("PRIVMSG %s :%s", to, message)
}

func CheckAdmin(nick string, ident string, addr string, adminaddr string) bool {
	RequesterAddr := fmt.Sprintf("%s!%s@%s", nick, ident, addr)
	if RequesterAddr == adminaddr {
		return true
	}
	return false
}
