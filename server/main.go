package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

// Shared in-memory key-value store
var (
	dataStore  = make(map[string]string) // The key-value store
	storeMutex sync.Mutex                // Mutex to synchronize access to the store
)

// Function to handle each client connection
func handleClient(conn net.Conn) {
	defer conn.Close() // Ensure the connection is closed when the function exits
	reader := bufio.NewReader(conn)

	for {
		// Read input from the client
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Client disconnected: %v\n", err)
			return
		}

		// Trim and split the message into parts
		message = strings.TrimSpace(message)
		parts := strings.Fields(message)
		if len(parts) == 0 {
			conn.Write([]byte("Invalid command\n"))
			continue
		}

		// Extract the command and process it
		command := strings.ToUpper(parts[0])
		switch command {
		case "PUT":
			// Validate the command format
			if len(parts) != 3 {
				conn.Write([]byte("Usage: PUT <key> <value>\n"))
				continue
			}
			key, value := parts[1], parts[2]

			// Store the key-value pair
			storeMutex.Lock()
			dataStore[key] = value
			storeMutex.Unlock()
			conn.Write([]byte("OK\n"))
		case "GET":
			// Validate the command format
			if len(parts) != 2 {
				conn.Write([]byte("Usage: GET <key>\n"))
				continue
			}
			key := parts[1]

			// Retrieve the value for the given key
			storeMutex.Lock()
			value, exists := dataStore[key]
			storeMutex.Unlock()
			if exists {
				conn.Write([]byte(fmt.Sprintf("%s\n", value)))
			} else {
				conn.Write([]byte("Key not found\n"))
			}
		case "DELETE":
			// Validate the command format
			if len(parts) != 2 {
				conn.Write([]byte("Usage: DELETE <key>\n"))
				continue
			}
			key := parts[1]

			// Remove the key-value pair if it exists
			storeMutex.Lock()
			_, exists := dataStore[key]
			if exists {
				delete(dataStore, key)
				conn.Write([]byte("Deleted\n"))
			} else {
				conn.Write([]byte("Key not found\n"))
			}
			storeMutex.Unlock()
		case "LIST":
			// List all key-value pairs in the store
			storeMutex.Lock()
			for k, v := range dataStore {
				conn.Write([]byte(fmt.Sprintf("%s: %s\n", k, v)))
			}
			storeMutex.Unlock()
		default:
			// Handle unknown commands
			conn.Write([]byte("Unknown command\n"))
		}
	}
}

func main() {
	// Start the server on port 3000
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server is running on port 3000...")

	// Accept and handle incoming client connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go handleClient(conn) // Spawn a goroutine for each client
	}
}
