package goirc

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

type Conn struct {
	con    net.Conn
	reader *bufio.Reader
	err    chan error
}

func Connect(addr string, port int) (*Conn, error) {
	address := fmt.Sprintf("%s:%d", addr, port)
	for {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			fmt.Printf("Error: %s, retrying in 5 seconds\r\n", err)
			time.Sleep(5 * time.Second)
		} else {
			con := &Conn{
				con:    conn,
				reader: bufio.NewReader(conn),
			}
			return con, nil
		}
	}
}

func (c *Conn) ReadLine() (string, error) {
	for {
		c.con.SetDeadline(time.Now().Add(300 * time.Second)) //If nothing happens in 5 minutes, there's a "ping timeout"
		line, err := c.reader.ReadString(byte('\n'))
		if err != nil {
			return "", err
		} else {
			line = strings.Trim(line, "\r\n")
			return line, nil
		}
	}
}

func (c *Conn) WriteLine(data string) error { //These need to be more fault tolerant
	fmt.Println("Sending: ", data)
	c.con.Write([]byte(data + "\r\n"))
	return nil
}

func (c *Conn) PrintLine(msg string, args ...interface{}) { //This needs to be more fault tolerant as well.
	fmt.Printf("Sending: "+msg, args...)
	fmt.Fprintf(c.con, msg, args...)
}

/*
* This may cause issues when called, the bot will try and reconnect
* We should set a value to make sure we're "really" disconnecting
 */
func (c *Conn) Disconnect() {
	fmt.Println("Sending Disconnect")
	c.con.Close()
}
