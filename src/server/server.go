package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const (
	SERVER_ADDRESS = ":8000"
	SERVER_TYPE    = "tcp"
	SERVER_HOST    = "127.0.0.1"
)

func handleMessage(connection net.Conn) {
	defer connection.Close()

	// A loop that reads data from the buffer
	for {
		clientData, error := bufio.NewReader(connection).ReadString('\n')
		if error != nil {
			fmt.Println("Connection error:", error)
			return
		}
		fmt.Println("Received:", string(clientData))
		fmt.Fprintln(connection, "Received")
	}
}

func main() {
	// Building the address string and listening through the socket
	listener, error := net.Listen(SERVER_TYPE, SERVER_ADDRESS)
	fmt.Println(strings.ToUpper(SERVER_TYPE) + " socket running on " + listener.Addr().String() + ".")

	// Printing an error in case socket can't be listened to
	if error != nil {
		fmt.Println(error)
		return
	}
	defer listener.Close()

	for {
		connection, error := listener.Accept()
		fmt.Println("Connection - " + connection.RemoteAddr().String())

		if error != nil {
			fmt.Println(error)
			return
		}
		go handleMessage(connection)
	}
}
