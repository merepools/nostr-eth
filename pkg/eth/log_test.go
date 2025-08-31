package log

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nbd-wtf/go-nostr"
)

func TestCreateTxLogEvent(t *testing.T) {
	// Create log data
	logData := map[string]interface{}{
		"hash":       "0x1234567890abcdef",
		"tx_hash":    "0xabcdef1234567890",
		"created_at": time.Now().Unix(),
		"updated_at": time.Now().Unix(),
		"nonce":      12345,
		"sender":     "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
		"to":         "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7",
		"value":      "0",
		"data":       nil,
		"status":     "pending",
	}

	// Create a Nostr event
	dataOutputter := NewJSONOutputter(logData)
	event, err := CreateTxLogEvent(dataOutputter, "private_key_here")
	if err != nil {
		t.Fatalf("Failed to create Nostr event: %v", err)
	}

	// Verify the event structure
	if event.Kind != 30000 {
		t.Errorf("Expected kind %d, got %d", 30000, event.Kind)
	}

	if event.CreatedAt == 0 {
		t.Error("Expected CreatedAt to be set")
	}

	// Check tags
	foundTxLog := false
	foundEthereum := false
	for _, tag := range event.Tags {
		if len(tag) >= 2 && tag[0] == "t" && tag[1] == "tx_log" {
			foundTxLog = true
		}
		if len(tag) >= 2 && tag[0] == "t" && tag[1] == "ethereum" {
			foundEthereum = true
		}
	}

	if !foundTxLog {
		t.Error("Expected tx_log tag not found")
	}
	if !foundEthereum {
		t.Error("Expected ethereum tag not found")
	}
}

func TestUpdateLogStatusEvent(t *testing.T) {
	logData := map[string]interface{}{
		"hash":       "0x1234567890abcdef",
		"tx_hash":    "0xabcdef1234567890",
		"created_at": time.Now().Unix(),
		"updated_at": time.Now().Unix(),
		"nonce":      12345,
		"sender":     "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
		"to":         "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7",
		"value":      "0",
		"data":       nil,
		"status":     "pending",
	}

	originalUpdatedAt := logData["updated_at"]

	// Update status without referencing original event
	event, err := UpdateLogStatusEvent(logData, "confirmed", "private_key_here")
	if err != nil {
		t.Fatalf("Failed to update log status: %v", err)
	}

	// Parse the content to verify the updated data
	var txLogEvent TxLogEvent
	err = json.Unmarshal([]byte(event.Content), &txLogEvent)
	if err != nil {
		t.Fatalf("Failed to parse event content: %v", err)
	}

	updatedLogData := txLogEvent.LogData
	if updatedLogData["status"] != "confirmed" {
		t.Errorf("Expected status %s, got %v", "confirmed", updatedLogData["status"])
	}

	if updatedLogData["updated_at"] == originalUpdatedAt {
		t.Error("Expected UpdatedAt to be updated")
	}

	// Test with original event reference
	originalEventID := "original_event_id_123"
	eventWithRef, err := UpdateLogStatusEvent(logData, "failed", "private_key_here", originalEventID)
	if err != nil {
		t.Fatalf("Failed to update log status with reference: %v", err)
	}

	// Check that the event reference is included
	foundEventRef := false
	for _, tag := range eventWithRef.Tags {
		if len(tag) >= 3 && tag[0] == "e" && tag[1] == originalEventID && tag[2] == "reply" {
			foundEventRef = true
			break
		}
	}

	if !foundEventRef {
		t.Error("Expected event reference tag not found")
	}
}

func TestUpdateTxLogEventWithReference(t *testing.T) {
	// Create log data
	logData := map[string]interface{}{
		"hash":       "0x1234567890abcdef",
		"tx_hash":    "0xabcdef1234567890",
		"created_at": time.Now().Unix(),
		"updated_at": time.Now().Unix(),
		"nonce":      12345,
		"sender":     "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
		"to":         "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7",
		"value":      "0",
		"data":       nil,
		"status":     "confirmed",
	}

	originalEventID := "original_event_id_456"

	// Create update event with reference
	event, err := UpdateTxLogEvent(logData, "private_key_here", originalEventID)
	if err != nil {
		t.Fatalf("Failed to create update event: %v", err)
	}

	// Check that the event reference is included
	foundEventRef := false
	for _, tag := range event.Tags {
		if len(tag) >= 3 && tag[0] == "e" && tag[1] == originalEventID && tag[2] == "reply" {
			foundEventRef = true
			break
		}
	}

	if !foundEventRef {
		t.Error("Expected event reference tag not found")
	}

	// Check that update tag is present
	foundUpdateTag := false
	for _, tag := range event.Tags {
		if len(tag) >= 2 && tag[0] == "t" && tag[1] == "update" {
			foundUpdateTag = true
			break
		}
	}

	if !foundUpdateTag {
		t.Error("Expected update tag not found")
	}

	// Parse the content to verify the event type
	var txLogEvent TxLogEvent
	err = json.Unmarshal([]byte(event.Content), &txLogEvent)
	if err != nil {
		t.Fatalf("Failed to parse event content: %v", err)
	}

	if txLogEvent.EventType != EventTypeTxLogUpdated {
		t.Errorf("Expected event type %s, got %s", EventTypeTxLogUpdated, txLogEvent.EventType)
	}
}

func TestParseTxLogEvent(t *testing.T) {
	// Create log data
	logData := map[string]interface{}{
		"hash":       "0x1234567890abcdef",
		"tx_hash":    "0xabcdef1234567890",
		"created_at": time.Now().Unix(),
		"updated_at": time.Now().Unix(),
		"nonce":      12345,
		"sender":     "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
		"to":         "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7",
		"value":      "0",
		"data":       nil,
		"status":     "pending",
	}

	// Create event data
	eventData := TxLogEvent{
		LogData:   logData,
		EventType: EventTypeTxLogCreated,
		Tags:      []string{"tx_log", "ethereum"},
	}

	// Marshal to JSON
	content, err := json.Marshal(eventData)
	if err != nil {
		t.Fatalf("Failed to marshal event data: %v", err)
	}

	// Create Nostr event
	nostrEvent := &nostr.Event{
		ID:        "event_id",
		PubKey:    "pubkey",
		CreatedAt: nostr.Timestamp(time.Now().Unix()),
		Kind:      30000,
		Tags:      []nostr.Tag{},
		Content:   string(content),
		Sig:       "signature",
	}

	// Parse back
	parsedEvent, err := ParseTxLogEvent(nostrEvent)
	if err != nil {
		t.Fatalf("Failed to parse Nostr event: %v", err)
	}

	// Verify parsed data
	if parsedEvent.EventType != EventTypeTxLogCreated {
		t.Errorf("Expected event type %s, got %s", EventTypeTxLogCreated, parsedEvent.EventType)
	}

	if parsedEvent.LogData["hash"] != logData["hash"] {
		t.Errorf("Expected hash %v, got %v", logData["hash"], parsedEvent.LogData["hash"])
	}
}

type TestJSONOutputter json.RawMessage

func NewJSONOutputter(data map[string]interface{}) JSONOutputter {
	b, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	return TestJSONOutputter(b)
}

func (t TestJSONOutputter) ToJSON() []byte {
	return []byte(t)
}

func TestDataFlatteningWithAddresses(t *testing.T) {
	// Create log data with dynamic data containing addresses
	logData := map[string]interface{}{
		"hash":       "0x1234567890abcdef1234567890abcdef12345678",
		"tx_hash":    "0xabcdef1234567890abcdef1234567890abcdef12",
		"created_at": time.Now().Unix(),
		"updated_at": time.Now().Unix(),
		"nonce":      12345,
		"sender":     "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
		"to":         "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7",
		"value":      "1000000000000000000",
		"status":     "pending",
		"data": map[string]interface{}{
			"topic": "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
			"from":  "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6", // Should become "p" tag
			"to":    "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7", // Should become "p" tag
			"value": "1000000000000000000",                        // Should become regular tag
			"token": "0x1234567890123456789012345678901234567890", // Should become "p" tag
		},
	}

	// Create a Nostr event
	dataOutputter := NewJSONOutputter(logData)
	event, err := CreateTxLogEvent(dataOutputter, "private_key_here")
	if err != nil {
		t.Fatalf("Failed to create Nostr event: %v", err)
	}

	// Check that addresses from data are converted to "p" tags
	foundFromAsP := false
	foundToAsP := false
	foundTokenAsP := false
	foundValueAsRegular := false
	foundTopicAsRegular := false

	for _, tag := range event.Tags {
		if len(tag) >= 2 {
			if tag[0] == "p" && tag[1] == "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6" {
				foundFromAsP = true
			}
			if tag[0] == "p" && tag[1] == "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7" {
				foundToAsP = true
			}
			if tag[0] == "p" && tag[1] == "0x1234567890123456789012345678901234567890" {
				foundTokenAsP = true
			}
			if tag[0] == "value" && tag[1] == "1000000000000000000" {
				foundValueAsRegular = true
			}
			if tag[0] == "topic" && tag[1] == "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" {
				foundTopicAsRegular = true
			}
		}
	}

	if !foundFromAsP {
		t.Error("Expected 'from' address to be converted to 'p' tag")
	}
	if !foundToAsP {
		t.Error("Expected 'to' address to be converted to 'p' tag")
	}
	if !foundTokenAsP {
		t.Error("Expected 'token' address to be converted to 'p' tag")
	}
	if !foundValueAsRegular {
		t.Error("Expected 'value' to remain as regular tag")
	}
	if !foundTopicAsRegular {
		t.Error("Expected 'topic' to remain as regular tag")
	}

	// Print all tags for debugging
	t.Logf("Event tags:")
	for _, tag := range event.Tags {
		if len(tag) >= 2 {
			t.Logf("  %s: %s", tag[0], tag[1])
		}
	}
}

func TestDynamicDataFlattening(t *testing.T) {
	// Test with different types of dynamic data
	testCases := []struct {
		name     string
		logData  map[string]interface{}
		expected map[string]string // key -> expected tag value
	}{
		{
			name: "ERC20 Transfer",
			logData: map[string]interface{}{
				"hash":   "0x1234567890abcdef",
				"status": "pending",
				"data": map[string]interface{}{
					"topic": "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
					"from":  "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
					"to":    "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7",
					"value": "1000000000000000000",
				},
			},
			expected: map[string]string{
				"topic": "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
				"value": "1000000000000000000",
			},
		},
		{
			name: "NFT Transfer",
			logData: map[string]interface{}{
				"hash":   "0xabcdef1234567890",
				"status": "confirmed",
				"data": map[string]interface{}{
					"topic":    "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
					"from":     "0x1234567890123456789012345678901234567890",
					"to":       "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
					"tokenId":  "12345",
					"contract": "0x9876543210987654321098765432109876543210",
				},
			},
			expected: map[string]string{
				"topic":   "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
				"tokenId": "12345",
			},
		},
		{
			name: "Complex Nested Data",
			logData: map[string]interface{}{
				"hash":   "0xcomplex123456",
				"status": "pending",
				"data": map[string]interface{}{
					"topic":    "0xcomplexhash1234567890abcdef",
					"name":     "Test Token",
					"symbol":   "TEST",
					"decimals": int64(18),
					"amount":   "500000000000000000",
				},
			},
			expected: map[string]string{
				"topic":    "0xcomplexhash1234567890abcdef",
				"name":     "Test Token",
				"symbol":   "TEST",
				"decimals": "18",
				"amount":   "500000000000000000",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a Nostr event
			dataOutputter := NewJSONOutputter(tc.logData)
			event, err := CreateTxLogEvent(dataOutputter, "private_key_here")
			if err != nil {
				t.Fatalf("Failed to create Nostr event: %v", err)
			}

			// Check that expected tags are present
			for expectedKey, expectedValue := range tc.expected {
				found := false
				for _, tag := range event.Tags {
					if len(tag) >= 2 && tag[0] == expectedKey && tag[1] == expectedValue {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected tag %s=%s not found", expectedKey, expectedValue)
				}
			}

			// Check that 0x addresses are converted to "p" tags
			data, ok := tc.logData["data"].(map[string]interface{})
			if ok {
				for _, value := range data {
					if strValue, ok := value.(string); ok && isEthereumAddress(strValue) {
						foundAsP := false
						for _, tag := range event.Tags {
							if len(tag) >= 2 && tag[0] == "p" && tag[1] == strValue {
								foundAsP = true
								break
							}
						}
						if !foundAsP {
							t.Errorf("Expected address %s to be converted to 'p' tag", strValue)
						}
					}
				}
			}
		})
	}
}
