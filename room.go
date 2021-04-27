package main

import (
	"net"
)

type room struct {
	name    string
	members map[net.Addr]*client
	host    *client
	game    *game
}

func (r *room) broadcastFromClient(sender *client, msg string) {
	for addr, m := range r.members {
		if sender.conn.RemoteAddr() != addr {
			m.msg(msg)
		}
	}
}

func (r *room) broadcastFromServer(msg string) {
	for _, m := range r.members {
		m.msg(msg)
	}
}
