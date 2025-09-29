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
go get github.com/citizenapp2/nostr-eth
```

## Current Implementation: Ethereum Transaction Logs and NIP-29 Groups

The current implementation includes two main features:

1. **Ethereum Transaction Logs (kind 111000)**: Generated when smart contracts emit events after transaction execution
2. **NIP-29 Groups (kinds 39000-39004)**: Full implementation of the NIP-29 group specification for decentralized group chat functionality

This serves as the first use case and will be referenced in the upcoming NIP.

### Key Features

#### Ethereum Transaction Logs
- **Ethereum Log Events**: Create Nostr events from Ethereum smart contract logs
- **Status Tracking**: Track transaction status changes (pending, confirmed, failed, cancelled)
- **Rich Tagging**: Comprehensive tag system for filtering and indexing
- **Event Parsing**: Parse Nostr events back into structured log data
- **Flexible Data Interface**: JSON-based interface for different data sources

#### NIP-29 Groups
- **Group Metadata Events (kind 39000)**: Create and manage group information
- **Group Messages (kind 39001)**: Send messages to groups with mentions and replies
- **Group Join/Leave Events (kinds 39002-39003)**: Track user membership changes
- **Group Moderation (kind 39004)**: Admin actions like bans, mutes, and role changes
- **Group Utilities**: Parse group identifiers, validate group IDs, and filter events

## Usage

### Creating Ethereum Log Events

```go
package main

import (
    "fmt"
    "time"
    "github.com/citizenapp2/nostr-eth/pkg/eth/log"
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

## NIP-29 Groups Usage

### Creating Group Metadata Events

```go
package main

import (
    "fmt"
    "github.com/comunifi/nostr-eth"
)

func main() {
    // Create a group with metadata
    groupID := "my-awesome-group"
    admins := []string{"admin1@example.com", "admin2@example.com"}
    moderators := []string{"mod1@example.com"}
    
    metadataEvent, err := nostreth.CreateGroupMetadataEvent(
        groupID,
        "My Awesome Group",
        "This is a group for awesome people",
        "https://example.com/group-pic.jpg",
        admins,
        moderators,
        true,  // private group
        false, // not closed
    )
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Created group metadata event: %s\n", metadataEvent.ID)
}
```

### Sending Group Messages

```go
// Send a message to the group
mentions := []string{"user1@example.com", "user2@example.com"}
messageEvent, err := nostreth.CreateGroupMessageEvent(
    groupID,
    "Hello everyone! How's it going?",
    "", // not a reply
    mentions,
)
if err != nil {
    panic(err)
}

// Reply to a message
replyEvent, err := nostreth.CreateGroupMessageEvent(
    groupID,
    "Thanks for the update!",
    "original-message-event-id", // reply to specific message
    []string{}, // no mentions
)
```

### Managing Group Membership

```go
// User joins the group
joinEvent, err := nostreth.CreateGroupJoinEvent(
    groupID,
    "newuser@example.com",
    "member", // role: admin, moderator, or member
)

// User leaves the group
leaveEvent, err := nostreth.CreateGroupLeaveEvent(
    groupID,
    "leavinguser@example.com",
    "Personal reasons", // optional reason
)
```

### Group Moderation

```go
// Ban a user for 24 hours
banEvent, err := nostreth.CreateGroupModerationEvent(
    groupID,
    "ban",
    "spammer@example.com",
    "Spamming the group",
    86400, // 24 hours in seconds
)

// Mute a user
muteEvent, err := nostreth.CreateGroupModerationEvent(
    groupID,
    "mute",
    "noisyuser@example.com",
    "Too many messages",
    3600, // 1 hour in seconds
)

// Promote user to moderator
promoteEvent, err := nostreth.CreateGroupModerationEvent(
    groupID,
    "promote",
    "helpfuluser@example.com",
    "Active and helpful member",
    0, // permanent action
)
```

### Parsing Group Events

```go
// Parse group metadata
metadata, err := nostreth.ParseGroupMetadataEvent(metadataEvent)
if err != nil {
    panic(err)
}
fmt.Printf("Group name: %s\n", metadata.Name)
fmt.Printf("Group admins: %v\n", metadata.Admins)

// Parse group message
message, err := nostreth.ParseGroupMessageEvent(messageEvent)
if err != nil {
    panic(err)
}
fmt.Printf("Message content: %s\n", message.Content)
fmt.Printf("Mentions: %v\n", message.Mentions)
```

### Group Utilities

```go
// Check if an event is a group event
if nostreth.IsGroupEvent(someEvent) {
    fmt.Println("This is a group event")
}

// Get group ID from any group event
groupID, err := nostreth.GetGroupIDFromEvent(someEvent)
if err != nil {
    panic(err)
}

// Filter events by group
allEvents := []*nostr.Event{event1, event2, event3}
groupEvents := nostreth.FilterGroupEventsByGroupID(allEvents, "my-group")

// Parse group identifier (format: "host'group-id")
host, groupID, err := nostreth.ParseGroupIdentifier("wss://groups.example.com'my-group")
if err != nil {
    panic(err)
}

// Format group identifier
identifier := nostreth.FormatGroupIdentifier("wss://groups.example.com", "my-group")
// Result: "wss://groups.example.com'my-group"

// Validate group ID
err = nostreth.ValidateGroupID("my-group")
if err != nil {
    fmt.Printf("Invalid group ID: %v\n", err)
}
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

### NIP-29 Group Data Structures

#### GroupMetadata

Represents group metadata for kind 39000 events:

```go
type GroupMetadata struct {
    Name        string   `json:"name"`
    About       string   `json:"about,omitempty"`
    Picture     string   `json:"picture,omitempty"`
    Admins      []string `json:"admins,omitempty"`
    Moderators  []string `json:"moderators,omitempty"`
    Private     bool     `json:"private,omitempty"`
    Closed      bool     `json:"closed,omitempty"`
    CreatedAt   int64    `json:"created_at"`
    UpdatedAt   int64    `json:"updated_at"`
}
```

#### GroupMessage

Represents a group message for kind 39001 events:

```go
type GroupMessage struct {
    Content   string   `json:"content"`
    ReplyTo   string   `json:"reply_to,omitempty"`
    Mentions  []string `json:"mentions,omitempty"`
    CreatedAt int64    `json:"created_at"`
}
```

#### GroupJoin

Represents a user joining a group for kind 39002 events:

```go
type GroupJoin struct {
    User      string `json:"user"`
    JoinedAt  int64  `json:"joined_at"`
    Role      string `json:"role,omitempty"` // "admin", "moderator", "member"
}
```

#### GroupLeave

Represents a user leaving a group for kind 39003 events:

```go
type GroupLeave struct {
    User     string `json:"user"`
    LeftAt   int64  `json:"left_at"`
    Reason   string `json:"reason,omitempty"`
}
```

#### GroupModeration

Represents moderation actions for kind 39004 events:

```go
type GroupModeration struct {
    Action    string `json:"action"` // "ban", "unban", "mute", "unmute", "promote", "demote"
    Target    string `json:"target"` // User being moderated
    Reason    string `json:"reason,omitempty"`
    Duration  int64  `json:"duration,omitempty"` // Duration in seconds for temporary actions
    CreatedAt int64  `json:"created_at"`
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

### NIP-29 Group Functions

#### Group Event Creation

- `CreateGroupMetadataEvent(groupID, name, about, picture string, admins, moderators []string, private, closed bool) (*nostr.Event, error)`
  - Creates a group metadata event (kind 39000)

- `CreateGroupMessageEvent(groupID, content, replyTo string, mentions []string) (*nostr.Event, error)`
  - Creates a group message event (kind 39001)

- `CreateGroupJoinEvent(groupID, user, role string) (*nostr.Event, error)`
  - Creates a group join event (kind 39002)

- `CreateGroupLeaveEvent(groupID, user, reason string) (*nostr.Event, error)`
  - Creates a group leave event (kind 39003)

- `CreateGroupModerationEvent(groupID, action, target, reason string, duration int64) (*nostr.Event, error)`
  - Creates a group moderation event (kind 39004)

#### Group Event Parsing

- `ParseGroupMetadataEvent(evt *nostr.Event) (*GroupMetadata, error)`
  - Parses a group metadata event

- `ParseGroupMessageEvent(evt *nostr.Event) (*GroupMessage, error)`
  - Parses a group message event

- `ParseGroupJoinEvent(evt *nostr.Event) (*GroupJoin, error)`
  - Parses a group join event

- `ParseGroupLeaveEvent(evt *nostr.Event) (*GroupLeave, error)`
  - Parses a group leave event

- `ParseGroupModerationEvent(evt *nostr.Event) (*GroupModeration, error)`
  - Parses a group moderation event

#### Group Utility Functions

- `IsGroupEvent(evt *nostr.Event) bool`
  - Checks if an event is a group-related event

- `GetGroupIDFromEvent(evt *nostr.Event) (string, error)`
  - Extracts the group ID from a group event

- `FilterGroupEventsByGroupID(events []*nostr.Event, groupID string) []*nostr.Event`
  - Filters events by group ID

- `ParseGroupIdentifier(groupIdentifier string) (host, groupID string, err error)`
  - Parses a group identifier in format "host'group-id"

- `FormatGroupIdentifier(host, groupID string) string`
  - Formats a host and group ID into a group identifier

- `ValidateGroupID(groupID string) error`
  - Validates that a group ID is properly formatted

- `GetEventTypeFromGroupEvent(evt *nostr.Event) string`
  - Determines the type of group event (metadata, message, join, leave, moderation)

## Nostr Event Structure (Kind 111000)

The module creates Nostr events with kind 111000 for Ethereum transaction logs:

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

## NIP-29 Group Event Structures

The module implements the full NIP-29 specification for group functionality using the following event kinds:

### Group Metadata Events (Kind 39000)

Group metadata events contain information about the group itself:

#### Event Content
```json
{
  "name": "Example Group",
  "about": "This is an example group for demonstrating NIP-29",
  "picture": "https://example.com/group-picture.jpg",
  "admins": ["admin1@example.com", "admin2@example.com"],
  "moderators": ["mod1@example.com", "mod2@example.com"],
  "private": true,
  "closed": false,
  "created_at": 1640995200,
  "updated_at": 1640995200
}
```

#### Standard Tags
- **`h`**: Group ID (unique identifier)
- **`p`**: Admin and moderator pubkeys with role markers
- **`t`**: Type tags (`group`, `metadata`, `private`, `closed`)

### Group Message Events (Kind 39001)

Group message events contain the actual messages sent to groups:

#### Event Content
```json
{
  "content": "Hello everyone! This is a test message in our group.",
  "reply_to": "original-message-event-id",
  "mentions": ["user1@example.com", "user2@example.com"],
  "created_at": 1640995200
}
```

#### Standard Tags
- **`h`**: Group ID
- **`e`**: Reference to original message (for replies)
- **`p`**: Mentioned users with `mention` marker
- **`t`**: Type tags (`group`, `message`)

### Group Join Events (Kind 39002)

Group join events track when users join groups:

#### Event Content
```json
{
  "user": "newuser@example.com",
  "joined_at": 1640995200,
  "role": "member"
}
```

#### Standard Tags
- **`h`**: Group ID
- **`p`**: User pubkey with `member` marker
- **`t`**: Type tags (`group`, `join`, role)

### Group Leave Events (Kind 39003)

Group leave events track when users leave groups:

#### Event Content
```json
{
  "user": "leavinguser@example.com",
  "left_at": 1640995200,
  "reason": "Personal reasons"
}
```

#### Standard Tags
- **`h`**: Group ID
- **`p`**: User pubkey with `former_member` marker
- **`t`**: Type tags (`group`, `leave`)

### Group Moderation Events (Kind 39004)

Group moderation events track admin actions:

#### Event Content
```json
{
  "action": "ban",
  "target": "spammer@example.com",
  "reason": "Spamming the group",
  "duration": 86400,
  "created_at": 1640995200
}
```

#### Standard Tags
- **`h`**: Group ID
- **`p`**: Target user pubkey with `target` marker
- **`t`**: Type tags (`group`, `moderation`, action type)

### Group Identifier Format

Groups are identified using the format: `host'group-id`

Examples:
- `wss://groups.example.com'my-awesome-group`
- `relay.nostr.com'dev-team`
- `localhost:8080'test-group`

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

## Implemented NIPs

This implementation provides full support for the following NIPs:

1. **NIP-29: Groups** (Fully implemented)
   - Group metadata events (kind 39000)
   - Group message events (kind 39001)
   - Group join/leave events (kinds 39002-39003)
   - Group moderation events (kind 39004)
   - Complete group identifier parsing and validation
   - Full event filtering and utility functions

2. **NIP-XX: Ethereum Transaction Logs** (Current implementation)
   - Standardizes kind 111000 for Ethereum logs
   - Defines tag structure and content format
   - Establishes event relationships

## Future NIPs

Additional NIPs planned for future implementation:

1. **Additional blockchain support** (Planned)
   - Bitcoin transaction logs
   - Solana transaction logs
   - Other blockchain ecosystems

2. **Advanced event types** (Planned)
   - DeFi protocol events (swaps, liquidity, etc.)
   - Governance events (voting, proposals, etc.)
   - NFT marketplace events

3. **Enhanced filtering and querying** (Planned)
   - Advanced tag-based filtering
   - Cross-chain event correlation
   - Real-time event streaming

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