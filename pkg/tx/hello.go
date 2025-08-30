package tx

import "fmt"

// HelloWorld returns a simple hello world message
func HelloWorld() string {
	return "Hello, World!"
}

// HelloWorldWithName returns a personalized hello message
func HelloWorldWithName(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

// PrintHelloWorld prints the hello world message to stdout
func PrintHelloWorld() {
	fmt.Println(HelloWorld())
}

// PrintHelloWorldWithName prints a personalized hello message to stdout
func PrintHelloWorldWithName(name string) {
	fmt.Println(HelloWorldWithName(name))
}
