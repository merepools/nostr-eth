# nostr-tx

A standalone Go module that provides basic functionality for Nostr transactions.

## Installation

```bash
go get github.com/kevin/nostr-tx
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/kevin/nostr-tx/pkg/tx"
)

func main() {
    // Get a simple hello world message
    message := tx.HelloWorld()
    fmt.Println(message) // Output: Hello, World!
    
    // Get a personalized message
    personalized := tx.HelloWorldWithName("Alice")
    fmt.Println(personalized) // Output: Hello, Alice!
    
    // Print directly to stdout
    tx.PrintHelloWorld() // Output: Hello, World!
    tx.PrintHelloWorldWithName("Bob") // Output: Hello, Bob!
}
```

## Available Functions

### `HelloWorld() string`
Returns a simple "Hello, World!" message.

### `HelloWorldWithName(name string) string`
Returns a personalized hello message with the provided name.

### `PrintHelloWorld()`
Prints the hello world message directly to stdout.

### `PrintHelloWorldWithName(name string)`
Prints a personalized hello message directly to stdout.

## Testing

Run the tests with:

```bash
go test ./pkg/tx
```

## Example

See the `example/` directory for a complete working example.

## License

This project is licensed under the MIT License - see the LICENSE file for details.