package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
)

/**
 * This type represents each input sent to the server
 *
 * Content
 * =======
 * nickname   - The name the client uses in the chat
 * inputType  - What kind of input. Determines how server responds. Examples: msg, connection init data
 * message    - The content of the message (mostly relevant in chat)
 */
type ClientInput struct {
	Nickname  string
	InputType string
	Message   string
}

type ServerResponse struct {
	ResponseType string
	Message      string
}

const (
	SERVER_ADDRESS = "127.0.0.1:8000"
	SERVER_TYPE    = "tcp"
)

/**
 * Initializing the ClientInfo struct so that the basic necessities
 * are ready when initiating chat communication.
 */
func initializeClient() string {
	var nickname string
	fmt.Println("Enter your nickname: ")
	fmt.Print(" > ")
	fmt.Scanf("%s", &nickname)
	return nickname
}

/**
 * The data client sends to a server when a connection is initialized.
 * This allows the server to connection various pieces of data to
 * offer better service.
 */
func sendClientDataToServer(connection net.Conn, nickname string) {
	clientInput := ClientInput{
		Nickname:  nickname,
		InputType: "initialize",
		Message:   "",
	}
	sendMessageToServer(connection, clientInput)
}

/**
 * Dials the server and checks for an error. If an error is encountered,
 * the program exits with error status 1. If no errors are found, the
 * connection is returned.
 */
func initiateConnection() net.Conn {
	connection, error := net.Dial(SERVER_TYPE, SERVER_ADDRESS)

	// Exiting with status 1 if connecting fails
	if error != nil {
		fmt.Println("Error when trying to dial the server!")
		os.Exit(1)
	}

	return connection
}

/**
 * A function for receiving inputs from the user. Uses bufio
 * to include text that comes after a whitespace.
 */
func receiveUserInput() string {
	input := bufio.NewReader(os.Stdin)
	clientMessage, _ := input.ReadString('\n')
	return clientMessage
}

func sendMessage(connection net.Conn, nickname string) {
	fmt.Println("Message: ")
	fmt.Print(nickname + " > ")
	clientMessage := receiveUserInput()

	fmt.Println(clientMessage)

	var userChoice string

	// Asking whether the message should be public or private
	for {
		fmt.Println("Will it be private? y = yes | n = no")
		fmt.Print("Input > ")
		fmt.Scanln(&userChoice)

		// A switch for handling scenarios based on whether the user wants the
		// message to be private
		switch userChoice {
		case "y":
			recipient := ""
			fmt.Println("Enter recipient name: ")
			fmt.Scanf("%s", &recipient)

			// The "---//---" symbol separates the recipient and the message
			clientInput := ClientInput{
				Nickname:  nickname,
				Message:   recipient + "---//---" + clientMessage + "\n",
				InputType: "private-message",
			}

			sendMessageToServer(connection, clientInput)

		case "n":
			clientInput := ClientInput{
				Nickname:  nickname,
				Message:   clientMessage + "\n",
				InputType: "message",
			}

			sendMessageToServer(connection, clientInput)
		default:
			fmt.Println("Invalid input received. Try again!")
			continue
		}
		break
	}

}

/**
 * Serializes the client input, checks for serialization errors and
 * sends the message to the server.
 */
func sendMessageToServer(connection net.Conn, clientInput ClientInput) {
	jsonInput, error := json.Marshal(clientInput)
	if error != nil {
		fmt.Println(error)
		return
	}
	// Sending the input to the server
	fmt.Fprint(connection, string(jsonInput)+"\n")
}

/**
 * A function that prints a response received from the server.
 * Only takes the active connection as a parameter.
 *
 * This function executes in its own thread to ensure responses
 * can be received constantly.
 */
func handleServerResponse(connection net.Conn, nickname string) {
	for {
		response, error := bufio.NewReader(connection).ReadString('\n')
		if error != nil {
			fmt.Println(error)
			os.Exit(1)
		}

		var serverResponse ServerResponse
		error = json.Unmarshal([]byte(response), &serverResponse)
		if error != nil {
			fmt.Println("Deserialization error:", error)
		}

		fmt.Println("")

		// A switch-case for handling various response types
		switch serverResponse.ResponseType {
		case "initialize-success":
			fmt.Println("Connection with server initialized successfully!")
			fmt.Println("Server: " + serverResponse.Message)

		case "client-message":
			fmt.Println(serverResponse.Message)

		case "error":
			fmt.Println("Server encountered an error.")
			fmt.Println("Server: " + serverResponse.Message)

		case "client-list-entry":
			fmt.Println("Found client: " + serverResponse.Message)

		default:
			fmt.Println("Server: " + serverResponse.Message)
		}
		fmt.Print(nickname + " > ")
	}
}

/**
 * Allows the client to change their channel based on the channel
 * number they type in.
 */
func changeChannel(connection net.Conn, nickname string) {
	var channelNumber int
	fmt.Println("Enter channel number.")
	fmt.Print(nickname + " > ")

	fmt.Scanln(&channelNumber)

	clientInput := ClientInput{
		Nickname:  nickname,
		Message:   strconv.Itoa(channelNumber),
		InputType: "change-channel",
	}
	sendMessageToServer(connection, clientInput)
}

func disconnectFromServer(connection net.Conn, nickname string) {
	clientInput := ClientInput{
		Nickname:  nickname,
		Message:   "",
		InputType: "disconnect",
	}
	sendMessageToServer(connection, clientInput)
}

func main() {
	fmt.Println("Welcome! You can chat in any of the following channels: 1, 2 or 3.")
	fmt.Println("Instructions: 1) Send message | 2) Select channel | 0) Disconnect")

	nickname := initializeClient()
	fmt.Println("Welcome", nickname+"!")

	// Initiating connection
	connection := initiateConnection()
	defer connection.Close()

	// A thread for receiving server responses
	go handleServerResponse(connection, nickname)

	// Sending client data to server
	sendClientDataToServer(connection, nickname)

	for {
		var userInput int

		fmt.Scanln(&userInput)

		/**
			 * Handling the user inputs
			 * Case 1 - User sends a message (either public or private)
		     * Case 2 - The user changes the channel
		     * Case 0 - The user disconnects and the application exits with code 0
		*/
		switch userInput {
		case 1:
			sendMessage(connection, nickname)

		case 2:
			changeChannel(connection, nickname)

		case 0:
			disconnectFromServer(connection, nickname)
			os.Exit(0)

		default:
			fmt.Println("Invalid input. Try again.")
		}
	}
}
