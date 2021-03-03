package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

// structure of server
type server struct {

	// map of clients connected to the server :
	// client name (key) & client (value)
	contacts map[string]*client

	// channel on which server receives commands from clients
	commands chan command
}

// function to instantiate new server
func newServer() *server {

	return &server{
		contacts: make(map[string]*client),
		commands: make(chan command),
	}
}

// function to run server
func (s *server) run() {

	log.Printf("running server...")

	// loop through incoming commands..
	for cmd := range s.commands {

		// based on the command id, execute desired functions
		switch cmd.id {
		case cmdName:
			// update client name to input
			s.name(cmd.client, cmd.args[1])
		case cmdJoin:
			// update client contact to input
			s.join(cmd.client, cmd.args[1])
		case cmdList:
			// return list of users (clients) connected to the server
			s.listRooms(cmd.client)
		case cmdMsg:
			// send input to client contact
			s.msg(cmd.client, cmd.args)
		case cmdQuit:
			// quit chat system
			s.quit(cmd.client)
		}
	}
}

// function to instantiate new client :
// called when a new client joins the server
func (s *server) newClient(conn net.Conn) {

	// instantiate client
	c := &client{
		conn:     conn,
		name:     "anonymous",
		commands: s.commands,
		contact:  "",
	}

	log.Printf("new client has joined : %s", conn.RemoteAddr().String())

	// start reading for input ( this is a blocking call on a separte go routine )
	c.readInput()
}

// function to assign an identifer (name) to a newly created client
func (s *server) name(c *client, name string) {

	// assign name to client
	c.name = name

	// update server guest list i.e currently connected users (clients)
	s.contacts[name] = c

	// give user feedback message
	c.msg(c, fmt.Sprintf("you will be known as %s", name))
}

// function to assign contact ( who a client is currently talkig to ) :
func (s *server) join(c *client, contactName string) {

	// check if a client by the given name exists on the server contacts list
	for k := range s.contacts {

		// if such a client exists...
		if contactName == k {

			// update client contact ( this contact is who messages will be sent to )
			c.contact = contactName

			// pass feedback
			c.msg(c, fmt.Sprintf("You are now talking to :%s", c.contact))
			break
		} else {

			// otherwise, pass feedback
			c.msg(c, fmt.Sprintf("No such user exists. check available users again."))
		}
	}
}

// function to display list of connected users :
// these clients are who you (a client) can join and then msg
func (s *server) listRooms(c *client) {

	var contacts []string

	// loop through available users
	for name := range s.contacts {

		// fetch all users except current client
		if name != c.name {
			contacts = append(contacts, name)
		}

	}

	// pass message
	c.msg(c, fmt.Sprintf("available rooms: %s", strings.Join(contacts, ", ")))
}

// function to pass a message to specified user (client)
func (s *server) msg(c *client, args []string) {

	// check if a user for given name exists on the server contacts map
	_, ok := s.contacts[c.name]

	// is so...
	if ok && c.contact != "" {

		log.Printf("attempting to send message")

		// join the entire mesage
		msg := strings.Join(args[1:], " ")

		// send the message
		c.msg(s.contacts[c.contact], c.name+" : "+msg)

	} else {
		// otherwise, prompt user to join to a user
		c.msg(c, "no one hears you. connect to a contact to initiate chat service. ")
	}

	println("from here")
}

// function to exit from chat
func (s *server) quit(c *client) {

	log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())

	// remove user from server contact list
	_, ok := s.contacts[c.name]
	if ok {
		delete(s.contacts, c.name)
	}

	// pass message
	c.msg(c, "skychat will miss you...")
	// close client connection
	c.conn.Close()
}
