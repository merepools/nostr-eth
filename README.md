# nostr-eth

A Go module that provides functionality for creating and managing Nostr events for Ethereum transaction logs.

## Installation

```bash
go get github.com/citizenwallet/nostr-eth
```

## Features

- **Ethereum Transaction Log Support**: Create and manage Nostr events for Ethereum transaction logs
- **ERC20 Transfer Support**: Specialized handling for ERC20 token transfers
- **Status Updates**: Track transaction status changes (pending, confirmed, failed, cancelled)
- **Nostr Event Creation**: Generate properly formatted Nostr events with appropriate tags
- **Event Parsing**: Parse Nostr events back into structured data

## Usage

### Basic Transaction Log Creation

```go
package main

import (
    "fmt"
    "github.com/citizenwallet/nostr-eth/pkg/eth"
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
    }

    // Create a Nostr event
    event, err := eth.CreateTxLogEvent(eth.NewMapDataOutputter(logData), "your_private_key_here")
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
updateEvent, err := eth.UpdateTxLogEvent(logData, "your_private_key_here")
if err != nil {
    panic(err)
}
```

### Event Parsing

```go
// Parse a Nostr event back into structured data
parsedEvent, err := eth.ParseTxLogEvent(nostrEvent)
if err != nil {
    panic(err)
}

fmt.Printf("Event type: %s\n", parsedEvent.EventType)
fmt.Printf("Log hash: %s\n", parsedEvent.LogData["hash"])
```

## Data Structures

### Log

Represents an Ethereum transaction log:

```go
type Log struct {
    Hash      string           `json:"hash"`
    TxHash    string           `json:"tx_hash"`
    CreatedAt time.Time        `json:"created_at"`
    UpdatedAt time.Time        `json:"updated_at"`
    Nonce     int64            `json:"nonce"`
    Sender    string           `json:"sender"`
    To        string           `json:"to"`
    Value     *big.Int         `json:"value"`
    Data      *json.RawMessage `json:"data"`
    ExtraData *json.RawMessage `json:"extra_data"`
    Status    LogStatus        `json:"status"`
}
```

### LogTransferData

Represents ERC20 transfer event data:

```go
type LogTransferData struct {
    To    string `json:"to"`
    From  string `json:"from"`
    Topic string `json:"topic"`
    Value string `json:"value"`
}
```

### LogStatus

Transaction status constants:

```go
const (
    LogStatusPending   LogStatus = "pending"
    LogStatusConfirmed LogStatus = "confirmed"
    LogStatusFailed    LogStatus = "failed"
    LogStatusCancelled LogStatus = "cancelled"
)
```

## Available Functions

### Core Functions

- `CreateERC20TransferLog(hash, txHash, sender, to, value string, nonce int64) (*Log, error)`
  - Creates a new log with ERC20 transfer data

- `CreateTxLogEvent(log Log, privateKey string) (*NostrEvent, error)`
  - Creates a new Nostr event for a transaction log

- `UpdateTxLogEvent(log Log, privateKey string) (*NostrEvent, error)`
  - Creates a Nostr event for updating a transaction log status

- `UpdateLogStatus(log *Log, status LogStatus)`
  - Updates the status of a log and sets the updated timestamp

### Utility Functions

- `ParseTxLogEvent(evt *NostrEvent) (*TxLogEvent, error)`
  - Parses a Nostr event back into structured data

- `IsERC20Transfer(log *Log) bool`
  - Checks if the log data represents an ERC20 transfer

- `GetTransferData(log *Log) (*LogTransferData, error)`
  - Extracts transfer data from a log

## Nostr Event Structure

The module creates Nostr events with:

- **Kind**: 30000 (custom kind for transaction logs)
- **Tags**: Properly formatted tags for indexing and filtering
  - `d`: Log hash (identifier)
  - `t`: Type tags (tx_log, ethereum, status, update)
  - `tx_hash`: Transaction hash
  - `sender`: Sender address
  - `to`: Recipient address (if available)

## Testing

Run the tests with:

```bash
go test ./pkg/eth -v
```

## Example

See the `example/` directory for a complete working example that demonstrates:

1. Creating ERC20 transfer logs
2. Generating Nostr events
3. Updating transaction status
4. Parsing events back
5. JSON representation
6. ERC20 transfer detection

Run the example with:

```bash
go run example/main.go
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.