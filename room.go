package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {

	// channel to store incoming messages ( which are later forwarded to other clients )
	forward chan []byte

	// channel for clients to join room
	join chan *client

	// channel for clients to exit room
	leave chan *client

	// map for storing current clients
	clients map[*client]bool
}

// newRoom makes new room
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {

	// run forever ~ in background ( separate go routine )
	for {

		// keep watching the three channels
		// and accordingly execute code for each case
		select {

		// join : client joins the room
		case client := <-r.join:
			// store client info
			r.clients[client] = true

		// exit : client leaves the room
		case client := <-r.leave:

			// delete client
			delete(r.clients, client)

			// close send channel
			close(client.send)

		// forward :
		case msg := <-r.forward:

			// forward arrived message to all clients
			for client := range r.clients {
				client.send <- msg
			}

		}
	}
}

// turning a room into an http handler
const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

// upgrade HTTP connection
var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

// method for room to act as handler
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// on receiving a request, upgrade to socket
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatalln("ServeHTTP:", err)
		return
	}

	// create client
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}

	// join channel
	r.join <- client
	// clean up after user leaves
	defer func() { r.leave <- client }()

	// write in separate go routine
	go client.write()

	// read in the main routine ( blocking operations ~ keeps connection alive )
	client.read()

}
