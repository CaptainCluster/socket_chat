package main

import (
	"fmt"
	"net"
)

const (
	SERVER_PORT = "8080"
	SERVER_TYPE = "tcp"
	SERVER_HOST = "127.0.0.1"
)

func main() {

	// Building the address string and listening through the socket
	var address string = ":" + SERVER_PORT
	connection, error := net.Listen(SERVER_TYPE, address)

	// Printing an error in case socket can't be listened to
	if error != nil {
		fmt.Println(error)
		return
	}

	for {
	}

	fmt.Println("Listening to " + SERVER_HOST + " on port " + SERVER_PORT + " .")

}
