package event

import (
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/comunifi/nostr-eth/pkg/neth"
	"github.com/nbd-wtf/go-nostr"
)

func TestCreateTxLogEvent(t *testing.T) {
	// Create log data
	logData := neth.Log{
		Hash:      "0x1234567890abcdef",
		TxHash:    "0xabcdef1234567890",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Nonce:     12345,
		Sender:    "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
		To:        "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7",
		Value:     big.NewInt(0),
		Data:      nil,
	}

	// Create a Nostr event
	event, err := CreateTxLogEvent(logData)
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

func TestUpdateTxLogEventWithReference(t *testing.T) {
	// Create log data
	logData := neth.Log{
		Hash:      "0x1234567890abcdef",
		TxHash:    "0xabcdef1234567890",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Nonce:     12345,
		Sender:    "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
		To:        "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7",
		Value:     big.NewInt(0),
		Data:      nil,
		ChainID:   "1",
	}

	originalEventID := "original_event_id_456"

	// Create update event with reference
	event, err := UpdateTxLogEvent(logData, originalEventID)
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
	logData := neth.Log{
		Hash:      "0x1234567890abcdef",
		TxHash:    "0xabcdef1234567890",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Nonce:     12345,
		Sender:    "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
		To:        "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7",
		Value:     big.NewInt(0),
		Data:      nil,
		ChainID:   "1",
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

	if parsedEvent.LogData.Hash != logData.Hash {
		t.Errorf("Expected hash %v, got %v", logData.Hash, parsedEvent.LogData.Hash)
	}
}

func NewJSONOutputter(data map[string]interface{}) json.RawMessage {
	b, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	return b
}

func TestDataFlatteningWithAddresses(t *testing.T) {
	data, err := json.Marshal(map[string]interface{}{
		"topic": "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
		"from":  "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6", // Should become "p" tag
		"to":    "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7", // Should become "p" tag
		"value": "1000000000000000000",                        // Should become regular tag
		"token": "0x1234567890123456789012345678901234567890", // Should become "p" tag
	})
	if err != nil {
		t.Fatalf("Failed to marshal data: %v", err)
	}

	var dataJSON json.RawMessage
	dataJSON = data

	// Create log data with dynamic data containing addresses
	logData := neth.Log{
		Hash:      "0x1234567890abcdef1234567890abcdef12345678",
		TxHash:    "0xabcdef1234567890abcdef1234567890abcdef12",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Nonce:     12345,
		Sender:    "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
		To:        "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7",
		Value:     big.NewInt(0),
		ChainID:   "1",
		Data:      &dataJSON,
	}

	// Create a Nostr event
	event, err := CreateTxLogEvent(logData)
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
	data1, err := json.Marshal(map[string]interface{}{
		"topic": "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
		"from":  "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
		"to":    "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7",
		"value": "1000000000000000000",
	})
	if err != nil {
		t.Fatalf("Failed to marshal data: %v", err)
	}

	var data1JSON json.RawMessage
	data1JSON = data1

	data2, err := json.Marshal(map[string]interface{}{
		"topic":    "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
		"from":     "0x1234567890123456789012345678901234567890",
		"to":       "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		"tokenId":  "12345",
		"contract": "0x9876543210987654321098765432109876543210",
	})
	if err != nil {
		t.Fatalf("Failed to marshal data: %v", err)
	}

	var data2JSON json.RawMessage
	data2JSON = data2

	data3, err := json.Marshal(map[string]interface{}{
		"topic":    "0xcomplexhash1234567890abcdef",
		"name":     "Test Token",
		"symbol":   "TEST",
		"decimals": 18,
		"amount":   "500000000000000000",
	})
	if err != nil {
		t.Fatalf("Failed to marshal data: %v", err)
	}

	var data3JSON json.RawMessage
	data3JSON = data3

	// Test with different types of dynamic data
	testCases := []struct {
		name     string
		logData  neth.Log
		expected map[string]string // key -> expected tag value
	}{
		{
			name: "ERC20 Transfer",
			logData: neth.Log{
				Hash: "0x1234567890abcdef",
				Data: &data1JSON,
			},
			expected: map[string]string{
				"topic": "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
				"value": "1000000000000000000",
			},
		},
		{
			name: "NFT Transfer",
			logData: neth.Log{
				Hash:    "0xabcdef1234567890",
				ChainID: "1",
				Data:    &data2JSON,
			},
			expected: map[string]string{
				"topic":   "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
				"tokenId": "12345",
			},
		},
		{
			name: "Complex Nested Data",
			logData: neth.Log{
				Hash:    "0xcomplex123456",
				ChainID: "1",
				Data:    &data3JSON,
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
			event, err := CreateTxLogEvent(tc.logData)
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
			var ldata map[string]interface{}
			err = json.Unmarshal(*tc.logData.Data, &ldata)
			if err != nil {
				t.Fatalf("Failed to unmarshal data: %v", err)
			}

			for _, value := range ldata {
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
		})
	}
}
