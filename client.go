package main

import "github.com/gorilla/websocket"

// structure for a user ( client )
type client struct {

	// websocket for the client
	socket *websocket.Conn

	// buffered channel on which received messages are queued to be forwarded to user ( via socket )
	send chan []byte

	// refers to the room the client is chatting in ( used to forward messages to everyone else in the room)
	room *room
}

// read from socket
func (c *client) read() {
	defer c.socket.Close()

	// continually..
	for {

		// read messages from socket
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			// if reading fails, socket is closed
			return
		}

		// send received messages to forward channel
		c.room.forward <- msg
	}
}

// write to socket
func (c *client) write() {
	defer c.socket.Close()

	// write received messages to socket
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)

		if err != nil {
			// if writing fails, socket is closed
			return
		}
	}
}
