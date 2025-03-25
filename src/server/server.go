package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
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

func sendResponseToClient(connection net.Conn, serverResponse ServerResponse) {
	jsonInput, error := json.Marshal(serverResponse)
	if error != nil {
		fmt.Println(error)
		return
	}
	// Sending the input to the server
	fmt.Fprint(connection, string(jsonInput)+"\n")
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

			channelNum := rand.Intn(len(chatInstance.Channels)-1) + 1
			chatInstance.Channels[channelNum].Clients = append(chatInstance.Channels[channelNum].Clients, client)

			serverResponse := ServerResponse{
				ResponseType: "initialize-success",
				Message:      "You have been added to channel " + strconv.Itoa(channelNum) + ".",
			}

			sendResponseToClient(connection, serverResponse)

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

				sendResponseToClient(i.Connection, serverResponse)
				break
			}

			serverResponseMainClient := ServerResponse{
				ResponseType: "default",
				Message:      "Message successfully sent",
			}
			sendResponseToClient(connection, serverResponseMainClient)
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

				sendResponseToClient(i.Connection, serverResponse)
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
						sendResponseToClient(connection, serverResponse)
					}
				}
			}

		case "change-channel":
			channelIndex := 0
			clientIndex := 0
			for _, channel := range chatInstance.Channels {
				for _, client := range channel.Clients {
					if client.Nickname != clientInput.Nickname {
						clientIndex++
						continue
					}
					// Removing the client
					channel.Clients = append(channel.Clients[:clientIndex], channel.Clients[clientIndex:]...)
					break
				}
				channelIndex++
				clientIndex = 0
			}

			message, error := strconv.Atoi(clientInput.Message)
			if error != nil {
				fmt.Println("Error when parsing int.")
			}

			chatInstance.Channels = append(chatInstance.Channels, chatInstance.Channels[message])

			serverResponse := ServerResponse{
				ResponseType: "channel-changed",
				Message:      "Your channel is now " + clientInput.Message,
			}
			sendResponseToClient(connection, serverResponse)

		default:
			serverResponse := ServerResponse{
				ResponseType: "error",
				Message:      "Request rejected due to unknown type. Try again.",
			}
			sendResponseToClient(connection, serverResponse)
		}
	}
}

func createChannels() []Channel {
	channels := []Channel{}

	for i := 1; i < 4; i++ {
		channel := Channel{
			ChannelId: i,
			Clients:   make([]Client, 0),
		}
		channels = append(channels, channel)
	}
	return channels
}

func main() {
	// Creating the chat instance before the server socket runs. It initially has
	// only one channel, but more will be created to serve the client needs.
	channels := createChannels()

	chat := ChatInstance{
		Channels: channels,
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
