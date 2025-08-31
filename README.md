# nostr-eth

A Go module that provides functionality for creating and managing Nostr events for Ethereum transaction logs. This serves as the reference implementation for an upcoming Nostr Improvement Proposal (NIP) that will standardize how blockchain transaction logs are stored in Nostr-native format.

## Overview

This module establishes a foundation for storing blockchain transaction logs in Nostr format, starting with Ethereum smart contract events. When smart contracts execute transactions, they emit events (logs) that contain structured data about the transaction. This implementation provides a standardized way to store these logs as Nostr events, enabling:

- **Decentralized Storage**: Store transaction logs in Nostr relays
- **Cross-Platform Compatibility**: Access logs from any Nostr client
- **Rich Querying**: Use Nostr's tag-based filtering system
- **Event Relationships**: Link related transactions and updates

## Installation

```bash
go get github.com/citizenwallet/nostr-eth
```

## Current Implementation: Ethereum Transaction Logs

The current implementation focuses on Ethereum transaction logs (kind 30000), which are generated when smart contracts emit events after transaction execution. This serves as the first use case and will be referenced in the upcoming NIP.

### Key Features

- **Ethereum Log Events**: Create Nostr events from Ethereum smart contract logs
- **Status Tracking**: Track transaction status changes (pending, confirmed, failed, cancelled)
- **Rich Tagging**: Comprehensive tag system for filtering and indexing
- **Event Parsing**: Parse Nostr events back into structured log data
- **Flexible Data Interface**: JSON-based interface for different data sources

## Usage

### Creating Ethereum Log Events

```go
package main

import (
    "fmt"
    "time"
    "github.com/citizenwallet/nostr-eth/pkg/eth/log"
)

func main() {
    // Example Ethereum log data from a smart contract event
    logData := map[string]interface{}{
        "hash":       "0x1234567890abcdef1234567890abcdef12345678", // Log hash
        "tx_hash":    "0xabcdef1234567890abcdef1234567890abcdef12", // Transaction hash
        "sender":     "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6", // From address
        "to":         "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7", // To address
        "value":      "1000000000000000000", // Value in wei (1 ETH)
        "nonce":      12345,
        "status":     "pending",
        "chain_id":   "1", // Ethereum mainnet
        "created_at": time.Now().Unix(),
        "data": map[string]interface{}{
            "topic": "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925", // Event signature
            "from":  "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
            "to":    "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7",
            "value": "1000000000000000000",
        },
    }

    // Create a JSON outputter
    jsonOutputter := log.NewGenericJSONOutputter(logData)

    // Create a Nostr event for the Ethereum log
    event, err := log.CreateTxLogEvent(jsonOutputter, "your_private_key_here")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Created Ethereum log event with kind: %d\n", event.Kind)
    fmt.Printf("Event ID: %s\n", event.ID)
}
```

### Updating Transaction Status

```go
// Update the log status when transaction is confirmed
logData["status"] = "confirmed"
logData["updated_at"] = time.Now().Unix()

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

### Parsing Events

```go
// Parse a Nostr event back into structured log data
parsedEvent, err := log.ParseTxLogEvent(nostrEvent)
if err != nil {
    panic(err)
}

fmt.Printf("Event type: %s\n", parsedEvent.EventType)
fmt.Printf("Log hash: %s\n", parsedEvent.LogData["hash"])
fmt.Printf("Transaction hash: %s\n", parsedEvent.LogData["tx_hash"])
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

### JSONOutputter Interface

Interface for outputting JSON data:

```go
type JSONOutputter interface {
    ToJSON() []byte
}
```

### GenericJSONOutputter

Simple implementation of JSONOutputter for map[string]interface{}:

```go
type GenericJSONOutputter json.RawMessage
```

## Available Functions

### Core Functions

- `CreateTxLogEvent(log JSONOutputter, privateKey string) (*nostr.Event, error)`
  - Creates a new Nostr event for an Ethereum transaction log

- `UpdateTxLogEvent(logData map[string]interface{}, privateKey string, originalEventID ...string) (*nostr.Event, error)`
  - Creates a Nostr event for updating a transaction log status

- `UpdateLogStatusEvent(logData map[string]interface{}, newStatus string, privateKey string, originalEventID ...string) (*nostr.Event, error)`
  - Creates a Nostr event for updating log status with new status and timestamp

### Utility Functions

- `ParseTxLogEvent(evt *nostr.Event) (*TxLogEvent, error)`
  - Parses a Nostr event back into structured data

- `GetTransferData(logData map[string]interface{}) (map[string]interface{}, error)`
  - Extracts transfer data from a log

- `NewGenericJSONOutputter(data map[string]interface{}) GenericJSONOutputter`
  - Creates a new GenericJSONOutputter instance

## Nostr Event Structure (Kind 30000)

The module creates Nostr events with kind 30000 for Ethereum transaction logs:

### Event Content
The event content contains a JSON object with:
- `log_data`: The complete Ethereum log data
- `event_type`: Either "tx_log_created" or "tx_log_updated"
- `tags`: Additional tags for categorization

### Standard Tags

- **`d`**: Log hash (unique identifier)
- **`t`**: Type tags for filtering
  - `tx_log`: Identifies this as a transaction log
  - `ethereum`: Identifies the blockchain
  - `update`: Present in update events
  - `{chain_id}`: Chain identifier (e.g., "1" for mainnet)
  - `{status}`: Transaction status (pending, confirmed, failed, cancelled)
- **`r`**: Transaction hash as reference
- **`p`**: Address tags for sender and recipient (0x addresses)
- **`amount`**: Value amount for filtering
- **`timestamp`**: Created timestamp for time-based filtering
- **`e`**: Reference to original event (for updates)

### Dynamic Tags from Log Data

The `data` field from the Ethereum log is flattened into tags:
- Ethereum addresses (0x...) are tagged with `p`
- Event topics (hashes) are tagged with their key name
- Numeric values are converted to strings
- Boolean values are converted to strings

## Example Ethereum Log Structure

```json
{
  "hash": "0x1234567890abcdef1234567890abcdef12345678",
  "tx_hash": "0xabcdef1234567890abcdef1234567890abcdef12",
  "sender": "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
  "to": "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7",
  "value": "1000000000000000000",
  "nonce": 12345,
  "status": "pending",
  "chain_id": "1",
  "created_at": 1640995200,
  "data": {
    "topic": "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925",
    "from": "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
    "to": "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7",
    "value": "1000000000000000000"
  }
}
```

## Event Types

```go
const (
    EventTypeTxLogCreated = "tx_log_created"
    EventTypeTxLogUpdated = "tx_log_updated"
)
```

## Future NIPs

This implementation serves as the foundation for multiple upcoming NIPs:

1. **NIP-XX: Ethereum Transaction Logs** (Current implementation)
   - Standardizes kind 30000 for Ethereum logs
   - Defines tag structure and content format
   - Establishes event relationships

2. **Future NIPs** (Planned)
   - Additional blockchain support (Bitcoin, Solana, etc.)
   - Different event types (transfers, swaps, governance, etc.)
   - Advanced filtering and querying capabilities

## Testing

Run the tests with:

```bash
go test ./pkg/eth -v
```

## Example

See the `example/` directory for a complete working example that demonstrates:

1. Creating Ethereum log data from smart contract events
2. Generating Nostr events with proper tagging
3. Updating transaction status
4. Parsing events back into structured data
5. Using the JSON outputter interface

Run the example with:

```bash
go run example/main.go
```

## Contributing

This is a reference implementation for an upcoming NIP. Contributions should focus on:

- Improving the Ethereum log implementation
- Adding comprehensive test coverage
- Enhancing documentation
- Preparing for NIP submission

## License

This project is licensed under the MIT License - see the LICENSE file for details.