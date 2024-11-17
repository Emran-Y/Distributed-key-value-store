package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	for {
		// Attempt to establish a connection with the server
		connection, err := net.Dial("tcp", "localhost:3000")
		if err != nil {
			fmt.Printf("client3 => Error connecting to server: %v\n", err)
			return
		}
		defer connection.Close() // Ensure connection is closed on exit

		fmt.Println("client3 => Connected to server. Enter commands:")

		// Initialize readers for user input and server responses
		userInputReader := bufio.NewReader(os.Stdin)
		serverResponseReader := bufio.NewReader(connection)

		for {
			// Prompt the user for a command
			fmt.Print("> ")
			command, _ := userInputReader.ReadString('\n')
			command = strings.TrimSpace(command) // Remove extra spaces or newlines

			// Send the user's command to the server
			_, err := connection.Write([]byte(command + "\n"))
			if err != nil {
				fmt.Printf("client3 => Error sending data: %v\n", err)
				break
			}

			// Read and print the server's response
			response, err := serverResponseReader.ReadString('\n')
			if err != nil {
				fmt.Printf("client3 => Server disconnected: %v\n", err)
				break
			}
			fmt.Println("client3 => Response:", strings.TrimSpace(response))
		}

		// Reconnection logic after the loop ends
		fmt.Println("client3 => Reconnecting...")
	}
}
