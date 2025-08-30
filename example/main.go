package main

import (
	"fmt"

	tx "github.com/citizenwallet/nostr-tx"
)

func main() {
	// Use the basic hello world function
	message := tx.HelloWorld()
	fmt.Println("Basic message:", message)

	// Use the personalized hello world function
	personalized := tx.HelloWorldWithName("Alice")
	fmt.Println("Personalized message:", personalized)

	// Use the print functions
	fmt.Println("Printing hello world:")
	tx.PrintHelloWorld()

	fmt.Println("Printing personalized hello world:")
	tx.PrintHelloWorldWithName("Bob")
}
