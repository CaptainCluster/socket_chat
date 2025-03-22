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

// Client input structs
type ClientInput struct {
	Nickname  string
	InputType string
	Message   string
	IsPrivate bool
}

// Chat instance structs
type Client struct {
	Nickname   string
	Connection net.Conn
}

type Channel struct {
	ChannelId int
	Clients   []Client
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

		fmt.Println(clientInput.InputType)
		/**
		 * Cases
		 * =====
		 * initialize - Puts the client on a channel
		 * message    - Sends a message to the channel
		 */
		switch clientInput.InputType {
		case "initialize":
			client := Client{
				Nickname:   clientInput.Nickname,
				Connection: connection,
			}
			chatInstance.Channels[0].Clients = append(chatInstance.Channels[0].Clients, client)
			fmt.Fprintln(connection, "You have been added to channel ", chatInstance.Channels[0].ChannelId, ".")

		case "message":
			// Sending a public message to everyone in the channel
			if !clientInput.IsPrivate {
				for _, i := range chatInstance.Channels[0].Clients {
					fmt.Fprintln(i.Connection, string(clientInput.Nickname+": "+clientInput.Message)+"\n")
				}
				continue
			}

			// A private message to a specified individual
			for _, i := range chatInstance.Channels[0].Clients {
				fmt.Fprintln(i.Connection, string(clientInput.Nickname+": "+clientInput.Message)+"\n")
			}

		case "client-list":
			for _, channel := range chatInstance.Channels {
				for _, client := range channel.Clients {
					if client.Nickname != clientInput.Nickname {
						continue
					}

					fmt.Fprintln(connection, "Here are the clients you can send a private message to:")
					for _, selectClient := range channel.Clients {
						if selectClient.Nickname == clientInput.Nickname {
							continue
						}
						fmt.Fprintln(connection, selectClient.Nickname)
					}
				}
			}

		default:
			fmt.Fprintln(connection, "Request rejected due to unknown type. Try again.")
		}
	}
}

func main() {
	var channelIdCounter int = 1
	// Creating the chat instance before the server socket runs. It initially has
	// only one channel, but more will be created to serve the client needs.
	channel := Channel{
		ChannelId: channelIdCounter,
		Clients:   make([]Client, 0),
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

		if error != nil {
			fmt.Println(error)
			return
		}

		// Creating a goroutine (Golang's own version of thread) to handle
		// client interaction
		go handleMessage(connection, chat)
	}
}
