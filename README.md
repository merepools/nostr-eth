# nostr-eth

A Go module that provides functionality for creating and managing Nostr events for Ethereum transaction logs.

## Installation

```bash
go get github.com/citizenwallet/nostr-eth
```

## Features

- **Ethereum Transaction Log Support**: Create and manage Nostr events for Ethereum transaction logs
- **Status Updates**: Track transaction status changes (pending, confirmed, failed, cancelled)
- **Nostr Event Creation**: Generate properly formatted Nostr events with appropriate tags
- **Event Parsing**: Parse Nostr events back into structured data
- **Data Outputter Interface**: Flexible interface for different data sources

## Usage

### Basic Transaction Log Creation

```go
package main

import (
    "fmt"
    "github.com/citizenwallet/nostr-eth/pkg/eth/log"
)

func main() {
    // Create log data
    logData := map[string]interface{}{
        "hash":       "0x1234567890abcdef1234567890abcdef12345678",
        "tx_hash":    "0xabcdef1234567890abcdef1234567890abcdef12",
        "sender":     "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
        "to":         "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7",
        "value":      "1000000000000000000", // value (1 token with 18 decimals)
        "nonce":      12345,
        "status":     "pending",
        "chain_id":   "1", // Ethereum mainnet
        "created_at": time.Now().Unix(),
    }

    // Create a data outputter
    dataOutputter := log.NewMapDataOutputter(logData)

    // Create a Nostr event
    event, err := log.CreateTxLogEvent(dataOutputter, "your_private_key_here")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Created event with kind: %d\n", event.Kind)
}
```

### Status Updates

```go
// Update the log status
logData["status"] = "confirmed"

// Create an update event
updateEvent, err := log.UpdateTxLogEvent(logData, "your_private_key_here")
if err != nil {
    panic(err)
}
```

### Status Update with Original Event Reference

```go
// Update log status with reference to original event
updateEvent, err := log.UpdateLogStatusEvent(logData, "confirmed", "your_private_key_here", "original_event_id")
if err != nil {
    panic(err)
}
```

### Event Parsing

```go
// Parse a Nostr event back into structured data
parsedEvent, err := log.ParseTxLogEvent(nostrEvent)
if err != nil {
    panic(err)
}

fmt.Printf("Event type: %s\n", parsedEvent.EventType)
fmt.Printf("Log hash: %s\n", parsedEvent.LogData["hash"])
```

## Data Structures

### TxLogEvent

Represents a Nostr event for transaction logs:

```go
type TxLogEvent struct {
    LogData   map[string]interface{} `json:"log_data"`
    EventType string                 `json:"event_type"`
    Tags      []string               `json:"tags,omitempty"`
}
```

### DataOutputter Interface

Interface for outputting map[string]interface{} data:

```go
type DataOutputter interface {
    OutputData() (map[string]interface{}, error)
}
```

### MapDataOutputter

Simple implementation of DataOutputter for map[string]interface{}:

```go
type MapDataOutputter struct {
    data map[string]interface{}
}
```

## Available Functions

### Core Functions

- `CreateTxLogEvent(log DataOutputter, privateKey string) (*nostr.Event, error)`
  - Creates a new Nostr event for a transaction log

- `UpdateTxLogEvent(logData map[string]interface{}, privateKey string, originalEventID ...string) (*nostr.Event, error)`
  - Creates a Nostr event for updating a transaction log status

- `UpdateLogStatusEvent(logData map[string]interface{}, newStatus string, privateKey string, originalEventID ...string) (*nostr.Event, error)`
  - Creates a Nostr event for updating log status with new status and timestamp

### Utility Functions

- `ParseTxLogEvent(evt *nostr.Event) (*TxLogEvent, error)`
  - Parses a Nostr event back into structured data

- `GetTransferData(logData map[string]interface{}) (map[string]interface{}, error)`
  - Extracts transfer data from a log

- `NewMapDataOutputter(data map[string]interface{}) *MapDataOutputter`
  - Creates a new MapDataOutputter instance

## Nostr Event Structure

The module creates Nostr events with:

- **Kind**: 30000 (custom kind for transaction logs)
- **Tags**: Properly formatted tags for indexing and filtering
  - `d`: Log hash (identifier)
  - `t`: Type tags (tx_log, ethereum, status, update, chain_id)
  - `r`: Transaction hash as reference
  - `p`: Address tags for sender and recipient (0x addresses)
  - `amount`: Value amount for filtering
  - `timestamp`: Created timestamp for time-based filtering
  - `e`: Reference to original event (for updates)

### Event Types

```go
const (
    EventTypeTxLogCreated = "tx_log_created"
    EventTypeTxLogUpdated = "tx_log_updated"
)
```

## Data Flattening

The module automatically flattens nested data structures into Nostr tags:

- Ethereum addresses (0x...) are tagged with `p`
- Hash values (like topics) are tagged with their key name
- Numeric values are converted to strings
- Nested maps are recursively flattened with prefixed keys
- Arrays are joined with commas

## Testing

Run the tests with:

```bash
go test ./pkg/eth -v
```

## Example

See the `example/` directory for a complete working example that demonstrates:

1. Creating transaction log data
2. Generating Nostr events
3. Updating transaction status
4. Parsing events back
5. Data outputter interface usage

Run the example with:

```bash
go run example/main.go
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.