package main

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
