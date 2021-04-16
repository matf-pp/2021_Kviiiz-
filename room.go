package main

import (
	"net"
)

type room struct {
	name    string
	members map[net.Addr]*client
	// TODO game pokazivac, svaki put novi
}

func (r *room) broadcast(sender *client, msg string) {
	for addr, m := range r.members {
		if sender.conn.RemoteAddr() != addr {
			m.msg(msg)
		}
	}
}

func (r *room) broadcastAll(msg string) {
	for _, m := range r.members {
		m.msg(msg)
	}
}
