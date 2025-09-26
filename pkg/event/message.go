package event

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nbd-wtf/go-nostr"
)

// NostrEventType represents the type of Nostr event for messages
const (
	EventTypeMessageCreated EventTypeMessage = "message_created"
	EventTypeMessageUpdated EventTypeMessage = "message_updated"
)

type EventTypeMessage string

// MessageEvent represents a Nostr event for simple text messages
type MessageEvent struct {
	MessageData MessageData      `json:"message_data"`
	EventType   EventTypeMessage `json:"event_type"`
	Tags        []string         `json:"tags,omitempty"`
}

// MessageData represents the core message data
type MessageData struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	ChainID   string    `json:"chain_id"`
	TxHash    string    `json:"tx_hash"`
	Group     string    `json:"group"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateMessageEvent creates a new Nostr event for a simple text message
func CreateMessageEvent(content, chainID, txHash, group, author string) (*nostr.Event, error) {
	// Generate a unique ID for the message
	messageID := generateMessageID(chainID, txHash, content, author)

	// Create the message data
	messageData := MessageData{
		ID:        messageID,
		Content:   content,
		ChainID:   chainID,
		TxHash:    txHash,
		Group:     group,
		Author:    author,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create the event data
	eventData := MessageEvent{
		MessageData: messageData,
		EventType:   EventTypeMessageCreated,
		Tags:        []string{"message", "group"},
	}

	// Marshal the event data
	contentJSON, err := json.Marshal(eventData)
	if err != nil {
		return nil, err
	}

	// Create the Nostr event
	evt := &nostr.Event{
		PubKey:    "", // Will be derived from private key
		CreatedAt: nostr.Timestamp(messageData.CreatedAt.Unix()),
		Kind:      1, // Standard kind for text messages
		Tags:      make([]nostr.Tag, 0),
		Content:   string(contentJSON),
	}

	// Add tags for better indexing and filtering
	evt.Tags = append(evt.Tags, []string{"d", messageID}) // Identifier

	// Type and category tags
	evt.Tags = append(evt.Tags, []string{"t", "message"}) // Type
	evt.Tags = append(evt.Tags, []string{"t", "text"})    // Content type

	// Group tag for filtering by group (NIP-29 compliant)
	evt.Tags = append(evt.Tags, []string{"h", group}) // Group ID

	// Chain-specific tag
	evt.Tags = append(evt.Tags, []string{"chain", chainID}) // Chain ID

	// Reference tags for transaction hash
	evt.Tags = append(evt.Tags, []string{"r", txHash}) // Transaction hash as reference

	// Author tag
	evt.Tags = append(evt.Tags, []string{"p", author}) // Author address

	return evt, nil
}

// UpdateMessageEvent creates a Nostr event for updating a message
func UpdateMessageEvent(messageData MessageData, event *nostr.Event) (*nostr.Event, error) {
	// Update the timestamp
	messageData.UpdatedAt = time.Now()

	// Create the event data
	eventData := MessageEvent{
		MessageData: messageData,
		EventType:   EventTypeMessageUpdated,
		Tags:        []string{"message", "group", "update"},
	}

	// Marshal the event data
	contentJSON, err := json.Marshal(eventData)
	if err != nil {
		return nil, err
	}

	// Create the Nostr event
	evt := &nostr.Event{
		PubKey:    "", // Will be derived from private key
		CreatedAt: nostr.Timestamp(messageData.UpdatedAt.Unix()),
		Kind:      1, // Standard kind for text messages
		Tags:      make([]nostr.Tag, 0),
		Content:   string(contentJSON),
	}

	// Add reference to original event if provided
	if event != nil {
		evt.Tags = append(evt.Tags, []string{"e", event.ID, "reply"}) // Reference to original event
	}

	// Add tags for better indexing and filtering
	evt.Tags = append(evt.Tags, []string{"d", messageData.ID}) // Identifier

	// Type and category tags
	evt.Tags = append(evt.Tags, []string{"t", "message"}) // Type
	evt.Tags = append(evt.Tags, []string{"t", "text"})    // Content type
	evt.Tags = append(evt.Tags, []string{"t", "update"})  // Update marker

	// Group tag for filtering by group (NIP-29 compliant)
	evt.Tags = append(evt.Tags, []string{"h", messageData.Group}) // Group ID

	// Chain-specific tag
	evt.Tags = append(evt.Tags, []string{"chain", messageData.ChainID}) // Chain ID

	// Reference tags for transaction hash
	evt.Tags = append(evt.Tags, []string{"r", messageData.TxHash}) // Transaction hash as reference

	// Author tag
	evt.Tags = append(evt.Tags, []string{"p", messageData.Author}) // Author address

	return evt, nil
}

// ParseMessageEvent parses a Nostr event back into a MessageEvent
func ParseMessageEvent(evt *nostr.Event) (*MessageEvent, error) {
	var messageEvent MessageEvent
	err := json.Unmarshal([]byte(evt.Content), &messageEvent)
	if err != nil {
		return nil, err
	}
	return &messageEvent, nil
}

// generateMessageID creates a unique ID for a message based on its content and metadata
func generateMessageID(chainID, txHash, content, author string) string {
	// Create a simple hash-like ID by combining key fields
	// In a real implementation, you might want to use a proper hash function
	return fmt.Sprintf("%s_%s_%s_%d", chainID, txHash, author, time.Now().UnixNano())
}

// GetGroupFromEvent extracts the group ID from a Nostr event (NIP-29 compliant)
func GetGroupFromEvent(evt *nostr.Event) (string, error) {
	for _, tag := range evt.Tags {
		if len(tag) >= 2 && tag[0] == "h" {
			return tag[1], nil
		}
	}
	return "", fmt.Errorf("group tag (h) not found in event")
}

// GetChainIDFromEvent extracts the chain ID from a Nostr event
func GetChainIDFromEvent(evt *nostr.Event) (string, error) {
	for _, tag := range evt.Tags {
		if len(tag) >= 2 && tag[0] == "chain" {
			return tag[1], nil
		}
	}
	return "", fmt.Errorf("chain tag not found in event")
}

// GetTxHashFromEvent extracts the transaction hash from a Nostr event
func GetTxHashFromEvent(evt *nostr.Event) (string, error) {
	for _, tag := range evt.Tags {
		if len(tag) >= 2 && tag[0] == "r" {
			return tag[1], nil
		}
	}
	return "", fmt.Errorf("transaction hash tag not found in event")
}

// IsMessageEvent checks if a Nostr event is a message event
func IsMessageEvent(evt *nostr.Event) bool {
	return evt.Kind == 1
}

// FilterEventsByGroup filters a list of events by group alias
func FilterEventsByGroup(events []*nostr.Event, group string) []*nostr.Event {
	var filtered []*nostr.Event
	for _, evt := range events {
		if groupTag, err := GetGroupFromEvent(evt); err == nil && groupTag == group {
			filtered = append(filtered, evt)
		}
	}
	return filtered
}

// FilterEventsByChainID filters a list of events by chain ID
func FilterEventsByChainID(events []*nostr.Event, chainID string) []*nostr.Event {
	var filtered []*nostr.Event
	for _, evt := range events {
		if chainTag, err := GetChainIDFromEvent(evt); err == nil && chainTag == chainID {
			filtered = append(filtered, evt)
		}
	}
	return filtered
}

// FilterEventsByTxHash filters a list of events by transaction hash
func FilterEventsByTxHash(events []*nostr.Event, txHash string) []*nostr.Event {
	var filtered []*nostr.Event
	for _, evt := range events {
		if txTag, err := GetTxHashFromEvent(evt); err == nil && txTag == txHash {
			filtered = append(filtered, evt)
		}
	}
	return filtered
}
