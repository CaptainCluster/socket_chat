package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

const (
	SERVER_ADDRESS = ":8000"
	SERVER_TYPE    = "tcp"
	SERVER_HOST    = "127.0.0.1"
)

type ClientInfo struct {
	Nickname string
	Channel  string
}

type ClientInput struct {
	ClientInfo ClientInfo
	InputType  string
	Message    string
	IsPrivate  bool
}

type Channel struct {
	ChannelId   int
	Clients     []ClientInfo
	Connections []net.Conn
}

type ChatInstance struct {
	Channels []Channel
}

func handleMessage(connection net.Conn, chatInstance ChatInstance) {
	defer connection.Close()

	// A loop that reads data from the buffer
	for {
		var clientInput ClientInput

		// Receiving data from a client
		clientData, error := bufio.NewReader(connection).ReadString('\n')

		// In case an error occurs, both client and server are notified
		if error != nil {
			fmt.Println("Connection error:", error)
			fmt.Fprintln(connection, "Failed to handle your request.")
			return
		}

		// Deserializing the data
		error = json.Unmarshal([]byte(clientData), &clientInput)

		// In case an error occurs, both client and server are notified
		if error != nil {
			fmt.Println("Deserialization error:", error)
			fmt.Fprintln(connection, "Failed to handle your request.")
			continue
		}

		/**
		 * Cases
		 * =====
		 * initialize - Puts the client on a channel
		 * message    - Sends a message to the channel
		 */
		switch clientInput.InputType {
		case "initialize":
			chatInstance.Channels[0].Clients = append(chatInstance.Channels[0].Clients, clientInput.ClientInfo)
			fmt.Fprintln(connection, "You have been added to channel ", chatInstance.Channels[0].ChannelId, ".")

		case "message":

			// Sending a public message to everyone in the channel
			if !clientInput.IsPrivate {
				for _, i := range chatInstance.Channels[0].Connections {
					fmt.Fprintln(i, string(clientInput.ClientInfo.Nickname+": "+clientInput.Message)+"\n")
				}
				continue
			}

			// A private message to a specified individual
			for _, i := range chatInstance.Channels[0].Connections {
				fmt.Fprintln(i, string(clientInput.ClientInfo.Nickname+": "+clientInput.Message)+"\n")
			}
		}

	}
}

func main() {
	var channelIdCounter int = 1
	// Creating the chat instance before the server socket runs. It initially has
	// only one channel, but more will be created to serve the client needs.
	channel := Channel{
		ChannelId: channelIdCounter,
		Clients:   make([]ClientInfo, 0),
	}
	chat := ChatInstance{
		Channels: []Channel{channel},
	}
	fmt.Println("Channel with the id", chat.Channels[0].ChannelId, "has been created.")

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
		chat.Channels[0].Connections = append(chat.Channels[0].Connections, connection)

		if error != nil {
			fmt.Println(error)
			return
		}

		// Creating a goroutine (Golang's own version of thread) to handle
		// client interaction
		go handleMessage(connection, chat)
	}
}
