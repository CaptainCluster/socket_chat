package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"slices"
	"strconv"
	"strings"
	"sync"
)

const (
	SERVER_ADDRESS = ":8000"
	SERVER_TYPE    = "tcp"
	SERVER_HOST    = "127.0.0.1"
	CHANNEL_AMOUNT = 3
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

func findClientChannel(chatInstance ChatInstance, nickname string) (Channel, bool) {
	for _, channel := range chatInstance.Channels {
		for _, client := range channel.Clients {
			if client.Nickname != nickname {
				continue
			}
			return channel, true
		}
	}
	return chatInstance.Channels[0], false
}

func handleMessage(connection net.Conn, chatInstance ChatInstance, mutex *sync.Mutex) {
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
			client := Client{
				Nickname:   clientInput.Nickname,
				Connection: connection,
			}
			channelNum := rand.Intn(len(chatInstance.Channels))

			if len(clientInput.Message) > 0 {
				channelNum, error = strconv.Atoi(clientInput.Message)
				if error != nil {
					fmt.Println("Error when parsing int.")
					return
				}
			}

			chatInstance.Channels[channelNum].Clients = append(chatInstance.Channels[channelNum].Clients, client)

			serverResponse := ServerResponse{
				ResponseType: "initialize-success",
				Message:      "You have been added to channel " + strconv.Itoa(chatInstance.Channels[channelNum].ChannelId) + ".",
			}

			sendResponseToClient(connection, serverResponse)

		case "message":

			channel, channelFound := findClientChannel(chatInstance, clientInput.Nickname)

			// Ensuring the sender is in a channel before the message is sent
			if !channelFound {
				fmt.Println("The user " + clientInput.Nickname + " tried to send a message to a channel, but was not a part of the channel.")
				return
			}

			// Sending a message to each client in the channel
			for _, client := range channel.Clients {

				// Client will not receive their own message
				if client.Nickname == clientInput.Nickname {
					continue
				}

				serverResponse := ServerResponse{
					ResponseType: "client-message",
					Message:      clientInput.Nickname + " > " + clientInput.Message,
				}

				sendResponseToClient(client.Connection, serverResponse)
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

			channel, channelFound := findClientChannel(chatInstance, clientInput.Nickname)

			// Ensuring the sender is in a channel before the message is sent
			if !channelFound {
				fmt.Println("The user " + clientInput.Nickname + " tried to send a message to a channel, but was not a part of the channel.")
				return
			}

			userExists := false

			// Sending a message to each client in the channel
			for _, client := range channel.Clients {

				// Skipping those who are not the recipient
				if client.Nickname != recipient {
					continue
				}

				serverResponse := ServerResponse{
					ResponseType: "client-message",
					Message:      "[Private] " + clientInput.Nickname + " > " + message,
				}

				// Sending the response to the recipient. Writing true to userExists to avoid error response
				fmt.Println("sending")
				sendResponseToClient(client.Connection, serverResponse)

				serverResponse = ServerResponse{
					ResponseType: "client-message",
					Message:      "Message sent successfully",
				}
				sendResponseToClient(connection, serverResponse)

				userExists = true
				break
			}

			// If the user was not found, an error response is sent to the client
			if !userExists {
				serverResponse := ServerResponse{
					ResponseType: "error",
					Message:      "User " + recipient + " not found.",
				}
				sendResponseToClient(connection, serverResponse)
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

			message, error := strconv.Atoi(clientInput.Message)
			if error != nil {
				fmt.Println("Error when parsing int.")
				return
			}

			if message > len(chatInstance.Channels) || message <= 0 {
				serverResponse := ServerResponse{
					ResponseType: "error",
					Message:      "Invalid channel name. Try again!",
				}
				sendResponseToClient(connection, serverResponse)
				return
			}

			channel, channelFound := findClientChannel(chatInstance, clientInput.Nickname)
			if !channelFound {
				fmt.Println("The user " + clientInput.Nickname + " tried to send a message to a channel, but was not a part of the channel.")
				return
			}

			mutex.Lock()
			for i := 0; i < len(channel.Clients); i++ {
				if channel.Clients[i].Nickname != clientInput.Nickname {
					continue
				}
				for j := 0; j < len(chatInstance.Channels); j++ {
					if chatInstance.Channels[j].ChannelId != message {
						continue
					}
					// Leaving if the new channel would be the same
					if i == j {
						break
					}
					chatInstance.Channels[j].Clients = append(chatInstance.Channels[j].Clients, channel.Clients[i])
					channel.Clients = slices.Delete(channel.Clients, i, i+1)
				}
			}
			mutex.Unlock()

			serverResponse := ServerResponse{
				ResponseType: "channel-changed",
				Message:      "Your channel is now " + clientInput.Message,
			}
			sendResponseToClient(connection, serverResponse)

		default:
			serverResponse := ServerResponse{
				ResponseType: "error",
				Message:      "Request rejected due to an unknown type. Try again.",
			}
			sendResponseToClient(connection, serverResponse)
		}
	}
}

func createChannels() []Channel {
	channels := []Channel{}

	for i := 1; i < CHANNEL_AMOUNT+1; i++ {
		channel := Channel{
			ChannelId: i,
			Clients:   make([]Client, 0),
		}
		channels = append(channels, channel)
	}
	return channels
}

func main() {
	var mutex sync.Mutex

	// Creating the chat instance before the server socket runs. It initially has
	// only one channel, but more will be created to serve the client needs.
	channels := createChannels()

	chat := ChatInstance{
		Channels: channels,
	}
	fmt.Println("A chat instance with " + strconv.Itoa(len(channels)) + " has been created.")

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

		// Returning, in case a connection fails
		if error != nil {
			fmt.Println(error)
			return
		}

		fmt.Println("Connection - " + connection.RemoteAddr().String())

		// Creating a goroutine (Golang's own version of thread) to handle
		// client interaction
		go handleMessage(connection, chat, &mutex)
	}
}
