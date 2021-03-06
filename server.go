package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NAME:
			s.name(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_START_GAME:
			s.startGame(cmd.client)
		case CMD_QUIT:
			s.quit(cmd.client)
		case CMD_HELP:
			s.msg(cmd.client, cmd.args)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("new client has joined: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		name:     "",
		commands: s.commands,
	}

	c.readInput()
}

func (s *server) name(c *client, args []string) {
	if len(args) < 2 {
		c.msg("name is required. usage: /name NAME")
		return
	}
	if args[1] == "" {
		c.msg("Invalid name.")
	} else {
		for _, room := range s.rooms {
			for _, cl := range room.members {
				if args[1] == cl.name {
					c.msg("The name already exists.")
					return
				}
			}
		}
		c.name = args[1]
		c.msg(fmt.Sprintf("Hello %s! :)", c.name))
	}
}

func (s *server) join(c *client, args []string) {
	if len(args) < 2 {
		c.msg("room name is required. usage: /join room_name")
		return
	}

	if c.name == "" {
		c.msg("You have to choose a name.\nCommand is /name your_name")
		return
	}

	roomName := args[1]

	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
			host:    c,
		}
		s.rooms[roomName] = r
	}

	if r.game == nil {
		r.members[c.conn.RemoteAddr()] = c

		s.quitCurrentRoom(c)
		c.room = r

		r.broadcastFromClient(c, fmt.Sprintf("%s joined the room", c.name))

		c.msg(fmt.Sprintf("welcome to %s", roomName))
	} else {
		c.msg("Can't join, game is in progress!")
	}

}

func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ", ")))
}

func (s *server) msg(c *client, args []string) {
	msg := strings.Join(args[0:], " ")
	room := c.room
	if room == nil {
		c.msg("You have to be in a room to send messages.")
	} else if room.game == nil {
		room.broadcastFromClient(c, c.name+": "+msg)
	} else {
		correct := room.game.attemptAnswer(c, msg)
		if correct == 1 {
			c.msg("Correct!")
		} else if correct == 0 {
			c.msg("Wrong answer!")
		} else if correct == -1 {
			c.msg("Invalid attempt!")
		}

		if room.game.moveToNextQuestion() {
			s.nextQuestion(room)
		}
	}
}

func (s *server) nextQuestion(r *room) {
	question, end := r.game.getNextQuestion()
	if !end {
		r.broadcastFromServer(question)
		rbrQuestion := r.game.br_pitanja
		checkFunc := func(rbr int, r **room) func() {
			return func() {
				// if r != nil && (*r) != nil && (*r).game != nil {
				// 	println((*r).game.br_pitanja, " ", rbr)
				// }
				if r != nil && (*r) != nil && (*r).game != nil && (*r).game.br_pitanja == rbr {
					s.nextQuestion(*r)
				}
			}
		}(rbrQuestion, &r)
		time.AfterFunc(20*time.Second, checkFunc)

	} else {
		r.broadcastFromServer("Game End!")
		r.broadcastFromServer("=======Points=======")
		for _, member := range r.members {
			r.broadcastFromServer("  " + member.name + ": " + strconv.Itoa(r.game.getPoints(member)))
		}
		r.game = nil
	}
}

func (s *server) startGame(c *client) {
	room := c.room
	if room == nil {
		c.msg("You have to join room first!\nCommand is /join room_name")
		return
	}
	if c != room.host {
		c.msg("Only the host can start the game!\nThe host is " + room.host.name)
		return
	}
	room.broadcastFromServer("Game start!")

	members_slice := make([]*client, 0)
	for _, member := range room.members {
		members_slice = append(members_slice, member)
	}
	fmt.Println(members_slice)
	room.game = newGame(members_slice)
	s.nextQuestion(room)

	// question, end := room.game.getNextQuestion()
	// if !end {
	//room.broadcastFromServer(question)
	// } else {
	// room.game = nil
	// room.broadcastFromServer("Game End!")
	// }
}

func (s *server) quit(c *client) {
	log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg("Bye bye :(")
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcastFromClient(c, fmt.Sprintf("%s has left the room", c.name))
		if len(oldRoom.members) == 0 {
			delete(s.rooms, c.room.name)
		} else {
			if oldRoom.game != nil {
				// Izbaci ga iz game-a
				oldRoom.game.leaveGame(c)
			}
			if c == oldRoom.host {
				// Dodeli novog hosta
				for _, cl := range oldRoom.members {
					oldRoom.host = cl
					cl.msg("You are now the host.")
					break
				}
			}
		}
	}
}
