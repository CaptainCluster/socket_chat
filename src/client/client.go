package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

/**
 * This type represents each input sent to the server
 *
 * Content
 * =======
 * clientInfo - Information about the client
 * inputType  - What kind of input. Determines how server responds. Examples: msg, connection init data
 * message    - The content of the message (mostly relevant in chat)
 * isPrivate  - Private or public message (mostly relevant in chat)
 */
type ClientInput struct {
	Nickname  string
	InputType string
	Message   string
	IsPrivate bool
	Recipient string
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
	fmt.Scanf("%s", &nickname)
	return nickname
}

func sendClientDataToServer(connection net.Conn, nickname string) {
	clientInput := ClientInput{
		Nickname:  nickname,
		InputType: "initialize",
		Message:   "Connected to the server\n",
		IsPrivate: false,
		Recipient: "",
	}
	sendMessageToServer(connection, clientInput)
}

func initiateConnection() net.Conn {
	connection, error := net.Dial(SERVER_TYPE, SERVER_ADDRESS)

	// Exiting with status 1 if connecting fails
	if error != nil {
		fmt.Println("Error when trying to dial the server!")
		os.Exit(1)
	}

	return connection
}

func sendMessage(connection net.Conn, nickname string) {
	var clientMessage string
	fmt.Println("Message: ")
	fmt.Scanln(&clientMessage)

	var isPrivate bool
	recipient := ""
	var userChoice string

	// Asking whether the message should be public or private
	for {
		fmt.Println("Will it be private? y = yes | n = no")
		fmt.Scanln(&userChoice)
		switch userChoice {

		case "y":
			isPrivate = true

			clientsReq := ClientInput{
				Nickname:  nickname,
				Message:   clientMessage + "\n",
				InputType: "client-list",
				IsPrivate: isPrivate,
				Recipient: "",
			}
			sendMessageToServer(connection, clientsReq)

			fmt.Println("Enter recipient name: ")
			fmt.Scanf("%s", &recipient)

		case "n":
			isPrivate = false
		default:
			fmt.Println("Invalid input received. Try again!")
			continue
		}
		break
	}

	clientInput := ClientInput{
		Nickname:  nickname,
		Message:   clientMessage + "\n",
		InputType: "message",
		IsPrivate: isPrivate,
		Recipient: recipient,
	}

	// Converting input to JSON
	sendMessageToServer(connection, clientInput)
}

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
 */
func handleServerResponse(connection net.Conn) {
	for {
		response, error := bufio.NewReader(connection).ReadString('\n')
		if error != nil {
			fmt.Println(error)
			os.Exit(1)
		}
		fmt.Println(string(response))
	}
}

func main() {
	nickname := initializeClient()
	fmt.Println("Welcome", nickname+"!")

	// Initiating connection
	connection := initiateConnection()
	defer connection.Close()

	// A thread for receiving server responses
	go handleServerResponse(connection)

	// Sending client data to server
	sendClientDataToServer(connection, nickname)

	for {
		var userInput int
		fmt.Println("1) Send message | 2) Select channel | 0) Disconnect: ")
		fmt.Scanln(&userInput)

		/**
			 * Handling the user inputs
			 * Case 1 - User sends a message (either public or private)
		     * Case 2 - The user changes the channel
		     * Case 0 - The user disconnects. They can reconnect or exit after this.
		*/
		switch userInput {
		case 1:
			sendMessage(connection, nickname)

		case 2:
			fmt.Println("Select channel")

		case 0:
			var userInput string
			connection.Close()
			fmt.Println("The connection has been closed.")
			fmt.Println("Would you like to end the program? y = yes | n = no")
			fmt.Scanln(&userInput)

			if userInput == "n" {
				fmt.Println("Goodbye! Thank you for using the application!")
				os.Exit(0)
			}

			// Since the program didn't exit, the connection will be re-initialized
			connection = initiateConnection()
			fmt.Println("The connection has been re-initialized.")

		default:
			fmt.Println("Invalid input. Try again.")
		}
	}
}
