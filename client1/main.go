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
		// Attempt to connect to the server
		connection, err := net.Dial("tcp", "localhost:3000")
		if err != nil {
			fmt.Printf("client1 => Error connecting to server: %v\n", err)
			return
		}
		defer connection.Close() // Ensure the connection is closed when the client exits

		fmt.Println("client1 => Connected to server. Enter commands:")

		// Readers for user input and server response
		userInputReader := bufio.NewReader(os.Stdin)
		serverResponseReader := bufio.NewReader(connection)

		for {
			// Prompt user for input
			fmt.Print("> ")
			command, _ := userInputReader.ReadString('\n')
			command = strings.TrimSpace(command)

			// Send the command to the server
			_, err := connection.Write([]byte(command + "\n"))
			if err != nil {
				fmt.Printf("client1 => Error sending data: %v\n", err)
				break
			}

			// Read and display the server's response
			response, err := serverResponseReader.ReadString('\n')
			if err != nil {
				fmt.Printf("client1 => Server disconnected: %v\n", err)
				break
			}
			fmt.Println("client1 => Response:", strings.TrimSpace(response))
		}

		// Attempt to reconnect to the server
		fmt.Println("client1 => Reconnecting...")
	}
}
