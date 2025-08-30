package tx

// Re-export all functions from the tx package
import "github.com/citizenwallet/nostr-tx/pkg/tx"

// HelloWorld returns a simple hello world message
func HelloWorld() string {
	return tx.HelloWorld()
}

// HelloWorldWithName returns a personalized hello message
func HelloWorldWithName(name string) string {
	return tx.HelloWorldWithName(name)
}

// PrintHelloWorld prints the hello world message to stdout
func PrintHelloWorld() {
	tx.PrintHelloWorld()
}

// PrintHelloWorldWithName prints a personalized hello message to stdout
func PrintHelloWorldWithName(name string) {
	tx.PrintHelloWorldWithName(name)
}
