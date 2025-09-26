package event

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/comunifi/nostr-eth/pkg/neth"
	"github.com/nbd-wtf/go-nostr"
)

// NostrEventType represents the type of Nostr event for transaction logs
const (
	KindTxLog = 30000

	EventTypeTxLogCreated EventTypeTxLog = "tx_log_created"
	EventTypeTxLogUpdated EventTypeTxLog = "tx_log_updated"
)

type EventTypeTxLog string

// TxLogEvent represents a Nostr event for transaction logs
type TxLogEvent struct {
	LogData   neth.Log       `json:"log_data"`
	EventType EventTypeTxLog `json:"event_type"`
	Tags      []string       `json:"tags,omitempty"`
}

// CreateTxLogEvent creates a new Nostr event for a transaction log
func CreateTxLogEvent(log neth.Log) (*nostr.Event, error) {
	// Create the event data
	eventData := TxLogEvent{
		LogData:   log,
		EventType: EventTypeTxLogCreated,
		Tags:      []string{"tx_log", "ethereum"},
	}

	// Marshal the event data
	content, err := json.Marshal(eventData)
	if err != nil {
		return nil, err
	}

	// Create the Nostr event
	evt := &nostr.Event{
		PubKey:    "", // Will be derived from private key
		CreatedAt: nostr.Timestamp(log.CreatedAt.Unix()),
		Kind:      KindTxLog, // Custom kind for transaction logs
		Tags:      make([]nostr.Tag, 0),
		Content:   string(content),
	}

	// Add tags for better indexing and filtering
	evt.Tags = append(evt.Tags, []string{"d", log.Hash}) // Identifier

	// Type and category tags
	evt.Tags = append(evt.Tags, []string{"t", "tx_log"})   // Type
	evt.Tags = append(evt.Tags, []string{"t", "ethereum"}) // Blockchain

	// Chain-specific tag
	evt.Tags = append(evt.Tags, []string{"t", log.ChainID}) // Chain ID

	// Reference tags for transaction hash
	evt.Tags = append(evt.Tags, []string{"r", log.TxHash}) // Transaction hash as reference

	// Address tags (using "p" for pubkey-like addresses)
	evt.Tags = append(evt.Tags, []string{"p", log.Sender}) // Sender address
	evt.Tags = append(evt.Tags, []string{"p", log.To})     // Recipient/Contract address

	// Amount/value tag for filtering
	evt.Tags = append(evt.Tags, []string{"amount", log.Value.String()})

	// Flatten data into tags
	if log.Data != nil {
		dataTags := flattenDataToTags(*log.Data)
		evt.Tags = append(evt.Tags, dataTags...)
	}

	return evt, nil
}

// UpdateTxLogEvent creates a Nostr event for updating a transaction log status
func UpdateTxLogEvent(log neth.Log, event *nostr.Event) (*nostr.Event, error) {
	// Create the event data
	eventData := TxLogEvent{
		LogData:   log,
		EventType: EventTypeTxLogUpdated,
		Tags:      []string{"tx_log", "ethereum", "update"},
	}

	// Marshal the event data
	content, err := json.Marshal(eventData)
	if err != nil {
		return nil, err
	}

	// Create the Nostr event
	evt := &nostr.Event{
		PubKey:    "", // Will be derived from private key
		CreatedAt: nostr.Timestamp(log.CreatedAt.Unix()),
		Kind:      KindTxLog, // Custom kind for transaction logs
		Tags:      make([]nostr.Tag, 0),
		Content:   string(content),
	}

	// Add reference to original event if provided
	if event != nil {
		evt.Tags = append(evt.Tags, []string{"e", event.ID, "reply"}) // Reference to original event
	}

	// Add tags for better indexing and filtering
	evt.Tags = append(evt.Tags, []string{"d", log.Hash}) // Identifier

	// Type and category tags
	evt.Tags = append(evt.Tags, []string{"t", "tx_log"})   // Type
	evt.Tags = append(evt.Tags, []string{"t", "ethereum"}) // Blockchain
	evt.Tags = append(evt.Tags, []string{"t", "update"})   // Update marker

	// Chain-specific tag
	evt.Tags = append(evt.Tags, []string{"t", log.ChainID}) // Chain ID

	// Reference tags for transaction hash
	evt.Tags = append(evt.Tags, []string{"r", log.TxHash}) // Transaction hash as reference

	// Address tags (using "p" for pubkey-like addresses)
	evt.Tags = append(evt.Tags, []string{"p", log.Sender}) // Sender address
	evt.Tags = append(evt.Tags, []string{"p", log.To})     // Recipient/Contract address

	// Amount/value tag for filtering
	evt.Tags = append(evt.Tags, []string{"amount", log.Value.String()})

	// Flatten data into tags
	if log.Data != nil {
		dataTags := flattenDataToTags(*log.Data)
		evt.Tags = append(evt.Tags, dataTags...)
	}

	return evt, nil
}

// ParseTxLogEvent parses a Nostr event back into a TxLogEvent
func ParseTxLogEvent(evt *nostr.Event) (*TxLogEvent, error) {
	var txLogEvent TxLogEvent
	err := json.Unmarshal([]byte(evt.Content), &txLogEvent)
	if err != nil {
		return nil, err
	}
	return &txLogEvent, nil
}

// isEthereumAddress checks if a string looks like an Ethereum address
func isEthereumAddress(value string) bool {
	// Ethereum addresses are 42 characters long (0x + 40 hex chars)
	if len(value) != 42 {
		return false
	}

	// Must start with 0x
	if !strings.HasPrefix(value, "0x") {
		return false
	}

	// Must contain only hex characters after 0x
	for _, char := range value[2:] {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
			return false
		}
	}

	return true
}

// flattenDataToTags flattens the data map into tags, using "p" for address values
// The data is dynamic and can be any event, but will always contain "topic" which is a hash
func flattenDataToTags(b []byte) []nostr.Tag {
	var tags []nostr.Tag

	var data map[string]interface{}
	err := json.Unmarshal(b, &data)
	if err != nil {
		return tags
	}

	for key, value := range data {
		if strValue, ok := value.(string); ok {
			if isEthereumAddress(strValue) {
				// Use "p" tag for 0x addresses
				tags = append(tags, []string{"p", strValue})
			} else if key == "topic" {
				// Topic is always a hash, keep as regular tag for better indexing
				tags = append(tags, []string{key, strValue})
			} else {
				// Use the key as tag name for other string values
				tags = append(tags, []string{key, strValue})
			}
		} else if intValue, ok := value.(int64); ok {
			// Convert int64 to string
			tags = append(tags, []string{key, fmt.Sprintf("%d", intValue)})
		} else if floatValue, ok := value.(float64); ok {
			// Convert float64 to string, but check if it's actually an integer
			if floatValue == float64(int64(floatValue)) {
				// It's an integer, format without decimal places
				tags = append(tags, []string{key, fmt.Sprintf("%d", int64(floatValue))})
			} else {
				// It's a real float, format with decimal places
				tags = append(tags, []string{key, fmt.Sprintf("%f", floatValue)})
			}
		} else if boolValue, ok := value.(bool); ok {
			// Convert bool to string
			tags = append(tags, []string{key, fmt.Sprintf("%t", boolValue)})
		}
	}

	return tags
}
