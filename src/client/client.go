package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type ClientInfo struct {
	nickname string
	channel  string
}

type ClientInput struct {
	clientInfo ClientInfo
	inputType  string
	message    string
	isPrivate  bool
}

const (
	SERVER_ADDRESS = "127.0.0.1:8000"
	SERVER_TYPE    = "tcp"
)

/**
 * Initializing the ClientInfo struct so that the basic necessities
 * are ready when initiating chat communication.
 */
func initializeClient() ClientInfo {
	var nickname string
	fmt.Println("Enter your nickname: ")
	fmt.Scanf("%s", &nickname)

	clientInfo := ClientInfo{
		nickname: nickname,
		channel:  "",
	}
	return clientInfo
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

func sendMessage(connection net.Conn, clientInfo ClientInfo) {
	var clientMessage string
	fmt.Println("Message: ")
	fmt.Scanln(&clientMessage)

	var isPrivate bool
	var userChoice string

	// Asking whether the message should be public or private
	for {
		fmt.Println("Will it be private? y = yes | n = no")
		fmt.Scanln(&userChoice)
		switch userChoice {
		case "y":
			isPrivate = true
			break
		case "n":
			isPrivate = false
			break
		default:
			fmt.Println("Invalid input received. Try again!")
			continue
		}
		break
	}

	clientInput := ClientInput{
		clientInfo: clientInfo,
		message:    clientMessage + "\n",
		inputType:  "message",
		isPrivate:  isPrivate,
	}

	fmt.Fprint(connection, clientInput)
}

/**
 * A function that prints a response received from the server.
 * Only takes the active connection as a parameter.
 */
func handleServerResponse(connection net.Conn) {
	response, error := bufio.NewReader(connection).ReadString('\n')
	if error != nil {
		fmt.Println(error)
	}
	fmt.Println(string(response))
}

func main() {
	clientInfo := initializeClient()
	fmt.Println("Welcome, " + clientInfo.nickname + "!")

	// Initiating connection
	connection := initiateConnection()
	defer connection.Close()

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
			sendMessage(connection, clientInfo)
			handleServerResponse(connection)

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
