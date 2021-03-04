# Skychat

Skychat is a basic encrypted messaging system, where clients connect to server via tcp.

## Skychat in action : ( Click to view video demo )

[![youtube demo](/README/thumb.png)](https://www.youtube.com/watch?v=30-eSvU4Rfk)


## Features

- Clients connect to server via tcp
- Command-line based
- Light weight : low on system resources.

## Installation

The application is madeup of a single binary so to get it running :
1. Download binary
2. Grant access : `chmod +x skychat`
3. Execute it ( default address ) : `./skychat`
4. Execute it ( custom  address ) : `./skychat -addr :1313`

That's all !

See [releases](github.com/FrappeFortyTwo/skychat/releases) for available packages.

## Building the project

1. [Install Go](https://goverse.dev/p1/)
2. Clone or fork this repository.
3. Open the directory with this repo's contents.
4. Run the command : `go build` ( this should create a new file within the present working directory ).
5. Finally, you can run the executable : `./skychat`. This should start the server at :8888 (default).
6. connect to this server using telnet : enter following command into a commandline while server is running.
7. `telnet localhost 8888`

## Sky Chat Commands :

Usage : /`<command>` [arguments]

* `name` : Specify your name.
* `list` : List connected users.
* `join` : Specify message recepient.
* `msg`  : Send message to recepient.
* `quit` : Exit Skychat.
* `help` : List help commands.

## How long it took ?
It took close to 10 hours to finish the task completely including researching, coding, documenting and everything else.
