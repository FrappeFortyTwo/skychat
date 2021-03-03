package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {

	// rooms will be known by client name
	rooms    map[string]*client
	commands chan command
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*client),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args[1])
		case CMD_JOIN:
			s.join(cmd.client, cmd.args[1])
		case CMD_ROOMS:
			s.listRooms(cmd.client)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("new client has joined: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
		contact:  "",
	}

	c.readInput()
}

func (s *server) nick(c *client, nick string) {

	c.nick = nick
	s.rooms[nick] = c

	// msg to self
	c.msg(c, fmt.Sprintf("all right, I will call you %s", nick))
}

func (s *server) join(c *client, roomName string) {

	// roomName means the key for clients on server

	c.contact = roomName
	// r, ok := s.rooms[roomName]
	// if !ok {
	// 	r = &room{
	// 		name:    roomName,
	// 		members: make(map[net.Addr]*client),
	// 	}
	// 	s.rooms[roomName] = r
	// }
	// r.members[c.conn.RemoteAddr()] = c

	// s.quitCurrentRoom(c)
	// c.room = r

	// r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))

	c.msg(c, fmt.Sprintf("You are now talking to :%s", roomName))

}

func (s *server) listRooms(c *client) {

	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(c, fmt.Sprintf("available rooms: %s", strings.Join(rooms, ", ")))
}

func (s *server) msg(c *client, args []string) {
	msg := strings.Join(args[1:len(args)], " ")

	// self message to other
	c.msg(s.rooms[c.contact], c.nick+": "+msg+"from server")
}

func (s *server) quit(c *client) {
	log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg(c, "sad to see you go =(")
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	c.msg(c, "quitting current room..")
}
