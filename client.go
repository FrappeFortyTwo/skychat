package main

import (
	"bufio"
	"crypto/rsa"
	"fmt"
	"log"
	"net"
	"strings"
)

// structure of a client i.e a user ( a new connection, will have this structure )
type client struct {

	// client connection details
	conn net.Conn

	// identifier for the client :
	// client will be known on the server by this name
	name string

	// identifier for the (other) client :
	// the other client (person), this client is talking to currently
	contact string

	// commands to facilitate chat system
	commands chan<- command

	// private
	private *rsa.PrivateKey

	// public
	public rsa.PublicKey
}

// function to read input
func (c *client) readInput() {

	// continuously...
	for {

		// read user input
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			// abort if an error occurs
			return
		}

		// process input, to parse commands
		msg = strings.Trim(msg, "\r\n")
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		// update client command for desired command
		switch cmd {

		case "/name":
			// specify your name
			c.commands <- command{
				id:     cmdName,
				client: c,
				args:   args,
			}
		case "/join":
			// connect to another user :
			// to be able to chat with him/her
			c.commands <- command{
				id:     cmdJoin,
				client: c,
				args:   args,
			}
		case "/list":
			// display all the available users on the server :
			// these are ones you ( a client ) can join and chat to
			c.commands <- command{
				id:     cmdList,
				client: c,
			}
		case "/msg":
			// send a message to the user ( another client ) you have joined
			c.commands <- command{
				id:     cmdMsg,
				client: c,
				args:   args,
			}
		case "/quit":
			// exit the chat system
			c.commands <- command{
				id:     cmdQuit,
				client: c,
			}
		case "/help":
			// return command list
			c.commands <- command{
				id:     cmdHelp,
				client: c,
			}
			// for any other command
		default:
			c.err(fmt.Errorf("unknown command: %s", cmd))
			c.msg(c, "* use '/help' to   command to send message to selected user")
		}
	}
}

// writes an error message current client
func (c *client) err(err error) {

	_, e := c.conn.Write([]byte("err: " + err.Error() + "\n"))
	if e != nil {
		log.Fatalln("unable to write to connection", e)
	}
}

// writes a message to specified client
func (c *client) msg(x *client, msg string) {

	// if contacting other client
	if c.private != x.private {

		dMsg := decrypt(msg, *x.private)

		// write message to client
		_, e := x.conn.Write([]byte("> " + dMsg + "\n"))
		if e != nil {
			log.Fatalln("unable to write over client connection")
		}

	} else {
		// write message to client
		_, e := x.conn.Write([]byte("> " + msg + "\n"))
		if e != nil {
			log.Fatalln("unable to write over client connection")
		}
	}

}
