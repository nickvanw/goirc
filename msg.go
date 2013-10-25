package goirc

import (
	"strings"
)

type Line struct {
	Nick, Ident, Host, Src string
	Cmd, Raw               string
	Args                   []string
}

/*
This is all shamelessly stolen from another IRC project
I just didn't feel like writing a message parser!
*/
func ParseLine(s string) *Line {
	line := &Line{Raw: s}
	if s[0] == ':' {
		// remove a source and parse it
		if idx := strings.Index(s, " "); idx != -1 {
			line.Src, s = s[1:idx], s[idx+1:]
		} else {
			// pretty sure we shouldn't get here ...
			return nil
		}

		// src can be the hostname of the irc server or a nick!user@host
		line.Host = line.Src
		nidx, uidx := strings.Index(line.Src, "!"), strings.Index(line.Src, "@")
		if uidx != -1 && nidx != -1 {
			line.Nick = line.Src[:nidx]
			line.Ident = line.Src[nidx+1 : uidx]
			line.Host = line.Src[uidx+1:]
		}
	}

	// now we're here, we've parsed a :nick!user@host or :server off
	// s should contain "cmd args[] :text"
	args := strings.SplitN(s, " :", 2)
	if len(args) > 1 {
		args = append(strings.Fields(args[0]), args[1])
	} else {
		args = strings.Fields(args[0])
	}
	line.Cmd = strings.ToUpper(args[0])
	if len(args) > 1 {
		line.Args = args[1:]
	}

	// So, I think CTCP and (in particular) CTCP ACTION are better handled as
	// separate events as opposed to forcing people to have gargantuan
	// handlers to cope with the possibilities.
	if (line.Cmd == "PRIVMSG" || line.Cmd == "NOTICE") &&
		len(line.Args[1]) > 2 &&
		strings.HasPrefix(line.Args[1], "\001") &&
		strings.HasSuffix(line.Args[1], "\001") {
		// WOO, it's a CTCP message
		t := strings.SplitN(strings.Trim(line.Args[1], "\001"), " ", 2)
		if len(t) > 1 {
			// Replace the line with the unwrapped CTCP
			line.Args[1] = t[1]
		}
		if c := strings.ToUpper(t[0]); c == "ACTION" && line.Cmd == "PRIVMSG" {
			// make a CTCP ACTION it's own event a-la PRIVMSG
			line.Cmd = c
		} else {
			// otherwise, dispatch a generic CTCP/CTCPREPLY event that
			// contains the type of CTCP in line.Args[0]
			if line.Cmd == "PRIVMSG" {
				line.Cmd = "CTCP"
			} else {
				line.Cmd = "CTCPREPLY"
			}
			line.Args = append([]string{c}, line.Args...)
		}
	}
	return line
}
