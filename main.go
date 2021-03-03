package main

import (
	"log"
	"net"
)

// entry point fo the program
func main() {

	// instantiate a server
	s := newServer()
	log.Printf("created new server")

	// run server in separate go routine
	go s.run()

	// start listening...
	listener, err := net.Listen("tcp", ":8888")

	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	} else {
		log.Printf("listening to port : 8888")
	}

	defer listener.Close()
	// sometime later.. close listener

	// continuously accept new connections
	for {

		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %s", err.Error())
			continue
		}

		go s.newClient(conn)
		log.Printf("added new client : %s", conn.RemoteAddr().String())
	}
}

// quick error check function
func checkErr(e error, m string) {
	if e != nil {
		log.Fatalln(m)
	}
}
