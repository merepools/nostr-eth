package log

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// NostrEventType represents the type of Nostr event for transaction logs
const (
	EventTypeTxLogCreated = "tx_log_created"
	EventTypeTxLogUpdated = "tx_log_updated"
)

// TxLogEvent represents a Nostr event for transaction logs
type TxLogEvent struct {
	LogData   map[string]interface{} `json:"log_data"`
	EventType string                 `json:"event_type"`
	Tags      []string               `json:"tags,omitempty"`
}

// NostrEvent represents a generic Nostr event structure
type NostrEvent struct {
	ID        string     `json:"id"`
	PubKey    string     `json:"pubkey"`
	CreatedAt int64      `json:"created_at"`
	Kind      int        `json:"kind"`
	Tags      [][]string `json:"tags"`
	Content   string     `json:"content"`
	Sig       string     `json:"sig"`
}

// DataOutputter defines an interface for outputting map[string]interface{} data
type DataOutputter interface {
	OutputData() (map[string]interface{}, error)
}

// MapDataOutputter is a simple implementation of DataOutputter for map[string]interface{}
type MapDataOutputter struct {
	data map[string]interface{}
}

// NewMapDataOutputter creates a new MapDataOutputter
func NewMapDataOutputter(data map[string]interface{}) *MapDataOutputter {
	return &MapDataOutputter{data: data}
}

// OutputData returns the underlying data map
func (m *MapDataOutputter) OutputData() (map[string]interface{}, error) {
	return m.data, nil
}

// CreateTxLogEvent creates a new Nostr event for a transaction log
func CreateTxLogEvent(log DataOutputter, privateKey string) (*NostrEvent, error) {
	logData, err := log.OutputData()
	if err != nil {
		return nil, err
	}

	// Create the event data
	eventData := TxLogEvent{
		LogData:   logData,
		EventType: EventTypeTxLogCreated,
		Tags:      []string{"tx_log", "ethereum"},
	}

	// Marshal the event data
	content, err := json.Marshal(eventData)
	if err != nil {
		return nil, err
	}

	// Create the Nostr event
	evt := &NostrEvent{
		PubKey:    "", // Will be derived from private key
		CreatedAt: time.Now().Unix(),
		Kind:      30000, // Custom kind for transaction logs
		Tags:      make([][]string, 0),
		Content:   string(content),
	}

	// Add tags for better indexing and filtering
	if hash, ok := logData["hash"].(string); ok {
		evt.Tags = append(evt.Tags, []string{"d", hash}) // Identifier
	}

	// Type and category tags
	evt.Tags = append(evt.Tags, []string{"t", "tx_log"})   // Type
	evt.Tags = append(evt.Tags, []string{"t", "ethereum"}) // Blockchain

	// Chain-specific tag
	if chainId, ok := logData["chain_id"].(string); ok {
		evt.Tags = append(evt.Tags, []string{"t", chainId}) // Chain ID
	}

	// Status tag
	if status, ok := logData["status"].(string); ok {
		evt.Tags = append(evt.Tags, []string{"t", status}) // Status
	}

	// Reference tags for transaction hash
	if txHash, ok := logData["tx_hash"].(string); ok {
		evt.Tags = append(evt.Tags, []string{"r", txHash}) // Transaction hash as reference
	}

	// Address tags (using "p" for pubkey-like addresses)
	if sender, ok := logData["sender"].(string); ok {
		evt.Tags = append(evt.Tags, []string{"p", sender}) // Sender address
	}

	if to, ok := logData["to"].(string); ok && to != "" {
		evt.Tags = append(evt.Tags, []string{"p", to}) // Recipient address
	}

	// Amount/value tag for filtering
	if value, ok := logData["value"].(string); ok && value != "0" {
		evt.Tags = append(evt.Tags, []string{"amount", value})
	}

	// Timestamp for time-based filtering
	if createdAt, ok := logData["created_at"].(int64); ok {
		evt.Tags = append(evt.Tags, []string{"timestamp", fmt.Sprintf("%d", createdAt)})
	}

	// Flatten data into tags
	if data, ok := logData["data"].(map[string]interface{}); ok {
		dataTags := flattenDataToTags(data)
		evt.Tags = append(evt.Tags, dataTags...)
	}

	// Note: In a real implementation, you would sign the event here
	// For now, we'll leave the ID and Sig empty as they require cryptographic operations

	return evt, nil
}

// UpdateTxLogEvent creates a Nostr event for updating a transaction log status
func UpdateTxLogEvent(logData map[string]interface{}, privateKey string, originalEventID ...string) (*NostrEvent, error) {
	// Create the event data
	eventData := TxLogEvent{
		LogData:   logData,
		EventType: EventTypeTxLogUpdated,
		Tags:      []string{"tx_log", "ethereum", "update"},
	}

	// Marshal the event data
	content, err := json.Marshal(eventData)
	if err != nil {
		return nil, err
	}

	// Create the Nostr event
	evt := &NostrEvent{
		PubKey:    "", // Will be derived from private key
		CreatedAt: time.Now().Unix(),
		Kind:      30000, // Custom kind for transaction logs
		Tags:      make([][]string, 0),
		Content:   string(content),
	}

	// Add reference to original event if provided
	if len(originalEventID) > 0 && originalEventID[0] != "" {
		evt.Tags = append(evt.Tags, []string{"e", originalEventID[0], "reply"}) // Reference to original event
	}

	// Add tags for better indexing and filtering
	if hash, ok := logData["hash"].(string); ok {
		evt.Tags = append(evt.Tags, []string{"d", hash}) // Identifier
	}

	// Type and category tags
	evt.Tags = append(evt.Tags, []string{"t", "tx_log"})   // Type
	evt.Tags = append(evt.Tags, []string{"t", "ethereum"}) // Blockchain
	evt.Tags = append(evt.Tags, []string{"t", "update"})   // Update marker

	// Chain-specific tag
	if chainId, ok := logData["chain_id"].(string); ok {
		evt.Tags = append(evt.Tags, []string{"t", chainId}) // Chain ID
	}

	// Status tag
	if status, ok := logData["status"].(string); ok {
		evt.Tags = append(evt.Tags, []string{"t", status}) // Status
	}

	// Reference tags for transaction hash
	if txHash, ok := logData["tx_hash"].(string); ok {
		evt.Tags = append(evt.Tags, []string{"r", txHash}) // Transaction hash as reference
	}

	// Address tags (using "p" for pubkey-like addresses)
	if sender, ok := logData["sender"].(string); ok {
		evt.Tags = append(evt.Tags, []string{"p", sender}) // Sender address
	}

	if to, ok := logData["to"].(string); ok && to != "" {
		evt.Tags = append(evt.Tags, []string{"p", to}) // Recipient address
	}

	// Amount/value tag for filtering
	if value, ok := logData["value"].(string); ok && value != "0" {
		evt.Tags = append(evt.Tags, []string{"amount", value})
	}

	// Timestamp for time-based filtering
	if createdAt, ok := logData["created_at"].(int64); ok {
		evt.Tags = append(evt.Tags, []string{"timestamp", fmt.Sprintf("%d", createdAt)})
	}

	// Flatten data into tags
	if data, ok := logData["data"].(map[string]interface{}); ok {
		dataTags := flattenDataToTags(data)
		evt.Tags = append(evt.Tags, dataTags...)
	}

	// Note: In a real implementation, you would sign the event here
	// For now, we'll leave the ID and Sig empty as they require cryptographic operations

	return evt, nil
}

// ParseTxLogEvent parses a Nostr event back into a TxLogEvent
func ParseTxLogEvent(evt *NostrEvent) (*TxLogEvent, error) {
	var txLogEvent TxLogEvent
	err := json.Unmarshal([]byte(evt.Content), &txLogEvent)
	if err != nil {
		return nil, err
	}
	return &txLogEvent, nil
}

// UpdateLogStatusEvent creates a Nostr event for updating log status
func UpdateLogStatusEvent(logData map[string]interface{}, newStatus string, privateKey string, originalEventID ...string) (*NostrEvent, error) {
	// Update the status and timestamp
	updatedLogData := make(map[string]interface{})
	for k, v := range logData {
		updatedLogData[k] = v
	}
	updatedLogData["status"] = newStatus
	updatedLogData["updated_at"] = time.Now().Unix()

	return UpdateTxLogEvent(updatedLogData, privateKey, originalEventID...)
}

// GetTransferData extracts transfer data from a log
func GetTransferData(logData map[string]interface{}) (map[string]interface{}, error) {
	data, ok := logData["data"].(map[string]interface{})
	if !ok {
		return nil, nil
	}

	return data, nil
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
func flattenDataToTags(data map[string]interface{}) [][]string {
	var tags [][]string

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
			// Convert float64 to string
			tags = append(tags, []string{key, fmt.Sprintf("%f", floatValue)})
		} else if boolValue, ok := value.(bool); ok {
			// Convert bool to string
			tags = append(tags, []string{key, fmt.Sprintf("%t", boolValue)})
		} else if mapValue, ok := value.(map[string]interface{}); ok {
			// Recursively flatten nested maps
			nestedTags := flattenDataToTags(mapValue)
			// Prefix nested keys with parent key to avoid conflicts
			for _, tag := range nestedTags {
				if len(tag) >= 2 {
					prefixedKey := fmt.Sprintf("%s_%s", key, tag[0])
					tags = append(tags, []string{prefixedKey, tag[1]})
				}
			}
		} else if sliceValue, ok := value.([]interface{}); ok {
			// Handle slices by joining with comma
			var strValues []string
			for _, item := range sliceValue {
				if strItem, ok := item.(string); ok {
					strValues = append(strValues, strItem)
				} else {
					strValues = append(strValues, fmt.Sprintf("%v", item))
				}
			}
			if len(strValues) > 0 {
				tags = append(tags, []string{key, strings.Join(strValues, ",")})
			}
		}
	}

	return tags
}
