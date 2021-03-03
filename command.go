package main

// custom type based on int ( for clarity )
type commandID int

const (

	// using iota to generate ever increasing numbers
	cmdName commandID = iota
	cmdJoin
	cmdList
	cmdMsg
	cmdQuit
	cmdHelp
)

// structure for a command
type command struct {
	id     commandID
	client *client
	args   []string
}
