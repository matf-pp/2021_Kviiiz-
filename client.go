package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	name     string
	room     *room
	commands chan<- command
}

var help_message string = `Commands:
/name new_name -> change name
/join room_name -> join room
/rooms -> list available rooms
/start -> start the quiz
/quit -> leave the room
/help -> see help`

var welcome_string string = `Welcome to Kviiiz!
You can use following commands
/name new_name -> change name
/join room_name -> join room
/rooms -> list available rooms
/start -> start the quiz
/quit -> leave the room
/help -> see help`

func (c *client) readInput() {
	c.msg(welcome_string)

	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/name":
			c.commands <- command{
				id:     CMD_NAME,
				client: c,
				args:   args,
			}
		case "/join":
			c.commands <- command{
				id:     CMD_JOIN,
				client: c,
				args:   args,
			}
		case "/rooms":
			c.commands <- command{
				id:     CMD_ROOMS,
				client: c,
			}

		case "/start":
			c.commands <- command{
				id:     CMD_START_GAME,
				client: c,
				args:   args,
			}
		case "/quit":
			c.commands <- command{
				id:     CMD_QUIT,
				client: c,
			}
		case "/help":
			c.msg(help_message)
		default:
			if cmd[0] == '/' {
				c.err(fmt.Errorf("unknown command: %s", cmd))
			} else {
				//msg
				c.commands <- command{
					id:     CMD_MSG,
					client: c,
					args:   args,
				}
			}
		}
	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte("err: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}
