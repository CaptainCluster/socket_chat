package main

import (
	"fmt"
)

type ClientInfo struct {
	nickname  string
	channel   string
	ipAddress string
}

type ClientInput struct {
	message   string
	isPrivate bool
}

/**
 * Initializing the ClientInfo struct so that the basic necessities
 * are ready when initiating chat communication.
 */
func initializeClient() ClientInfo {
	var nickname string
	fmt.Println("Enter your nickname: ")
	fmt.Scanf("%s", &nickname)

	clientInfo := ClientInfo{
		nickname:  nickname,
		ipAddress: "wee",
		channel:   "",
	}
	return clientInfo
}

func main() {
	clientInfo := initializeClient()
	fmt.Println("Welcome, " + clientInfo.nickname + "!")
	for {
		var userInput int
		fmt.Println("1) Send message | 2) Select channel | 0) Disconnect: ")
		fmt.Scanln(&userInput)

		switch userInput {
		case 1:
			fmt.Println("Send message")

		case 2:
			fmt.Println("Select channel")

		case 0:
			fmt.Println("Disconnect")
			return

		default:
			fmt.Println("Invalid input. Try again.")
		}
	}
}
