package main

import (
	"flag"
	"log"
	"net"
)

// entry point fo the program
func main() {

	// parse command-line arguments ( default : 8080 )
	var addr = flag.String("addr", ":8080", "Address for the app")
	flag.Parse()

	// instantiate a server
	s := newServer()
	log.Printf("created new server")

	// run server in separate go routine
	go s.run()

	// start listening...
	listener, err := net.Listen("tcp", *addr)

	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	} else {
		log.Printf("listening to port : %s", *addr)
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
