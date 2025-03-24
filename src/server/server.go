package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
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
}

// Server response structs
type ServerResponse struct {
	ResponseType string
	Message      string
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

			serverResponse := ServerResponse{
				ResponseType: "initialize-success",
				Message:      "You have been added to channel " + strconv.Itoa(chatInstance.Channels[0].ChannelId) + ".",
			}

			jsonInput, error := json.Marshal(serverResponse)
			if error != nil {
				fmt.Println(error)
				fmt.Fprintln(connection, "Server could not respond to your initialization request.")
				return
			}

			fmt.Fprintln(connection, string(jsonInput)+"\n")

		case "message":

			fmt.Println(clientInput.Message)

			for _, i := range chatInstance.Channels[0].Clients {

				// Client will not receive their own message
				if i.Nickname == clientInput.Nickname {
					continue
				}

				serverResponse := ServerResponse{
					ResponseType: "client-message",
					Message:      clientInput.Nickname + " > " + clientInput.Message,
				}

				jsonInput, error := json.Marshal(serverResponse)
				if error != nil {
					fmt.Println(error)
					fmt.Fprintln(connection, "Server could not send a message")
					return
				}

				fmt.Fprintln(i.Connection, string(jsonInput)+"\n")
			}
			continue

		case "private-message":
			split := strings.Split(clientInput.Message, "---//---")
			recipient := split[0]
			message := split[1]

			for _, i := range chatInstance.Channels[0].Clients {

				// Skipping those who are not the recipient
				if i.Nickname != recipient {
					continue
				}

				serverResponse := ServerResponse{
					ResponseType: "client-message",
					Message:      "[Private] " + clientInput.Nickname + " > " + message,
				}

				jsonInput, error := json.Marshal(serverResponse)
				if error != nil {
					fmt.Println(error)
					fmt.Fprintln(connection, "Server could not send a message")
					return
				}

				fmt.Fprintln(i.Connection, string(jsonInput)+"\n")
				break
			}

		case "client-list":
			for _, channel := range chatInstance.Channels {
				for _, client := range channel.Clients {
					if client.Nickname != clientInput.Nickname {
						continue
					}

					for _, selectClient := range channel.Clients {
						if selectClient.Nickname == clientInput.Nickname {
							continue
						}

						serverResponse := ServerResponse{
							ResponseType: "client-list-entry",
							Message:      selectClient.Nickname,
						}
						jsonInput, error := json.Marshal(serverResponse)
						if error != nil {
							fmt.Println(error)
							fmt.Fprintln(connection, "Server could not send a message")
							return
						}

						fmt.Fprintln(connection, string(jsonInput)+"\n")
					}
				}
			}

		default:
			serverResponse := ServerResponse{
				ResponseType: "error",
				Message:      "Request rejected due to unknown type. Try again.",
			}
			jsonInput, error := json.Marshal(serverResponse)
			if error != nil {
				fmt.Println(error)
				fmt.Fprintln(connection, "Server could not send a message")
				return
			}

			fmt.Fprintln(connection, string(jsonInput)+"\n")
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
