package main

import (
	"encoding/json"
	"fmt"
	"log"

	tx "github.com/citizenwallet/nostr-eth/pkg/eth"
)

func main() {
	fmt.Println("=== Nostr Ethereum Transaction Log Example ===")

	var err error

	// Example 1: Create a generic transaction log event
	fmt.Println("1. Creating Generic Transaction Log Event:")
	genericLogData := map[string]interface{}{
		"hash":       "0x1234567890abcdef1234567890abcdef12345678",
		"tx_hash":    "0xabcdef1234567890abcdef1234567890abcdef12",
		"created_at": 1640995200,
		"updated_at": 1640995200,
		"nonce":      12345,
		"sender":     "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
		"to":         "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7",
		"value":      "1000000000000000000",
		"status":     "pending",
		"data": map[string]interface{}{
			"topic": "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
			"from":  "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
			"to":    "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7",
			"value": "1000000000000000000",
		},
	}

	genericDataOutputter := tx.NewMapDataOutputter(genericLogData)
	var genericEvent *tx.NostrEvent
	genericEvent, err = tx.CreateTxLogEvent(genericDataOutputter, "your_private_key_here")
	if err != nil {
		log.Fatalf("Failed to create generic transaction event: %v", err)
	}

	// Parse the event to get the log data
	var txLogEvent tx.TxLogEvent
	err = json.Unmarshal([]byte(genericEvent.Content), &txLogEvent)
	if err != nil {
		log.Fatalf("Failed to parse event content: %v", err)
	}

	logData := txLogEvent.LogData
	fmt.Printf("   Log Hash: %v\n", logData["hash"])
	fmt.Printf("   Transaction Hash: %v\n", logData["tx_hash"])
	fmt.Printf("   From: %v\n", logData["sender"])
	fmt.Printf("   To: %v\n", logData["to"])
	fmt.Printf("   Status: %v\n", logData["status"])

	// Get transfer data
	transferData, err := tx.GetTransferData(logData)
	if err != nil {
		log.Fatalf("Failed to get transfer data: %v", err)
	}
	if transferData != nil {
		fmt.Printf("   Transfer Value: %v\n", transferData["value"])
		fmt.Printf("   Transfer Topic: %v\n", transferData["topic"])
	}

	// Example 2: Create a Nostr event for custom log data
	fmt.Println("\n2. Creating Nostr Event for Custom Log Data:")
	customLogData := map[string]interface{}{
		"hash":       "0x9876543210fedcba9876543210fedcba98765432",
		"tx_hash":    "0xfedcba0987654321fedcba0987654321fedcba09",
		"created_at": 1640995200, // Unix timestamp
		"updated_at": 1640995200,
		"nonce":      67890,
		"sender":     "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b8",
		"to":         "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b9",
		"value":      "1000000000000000000", // 1 ETH
		"data":       nil,
		"status":     "pending",
	}

	customDataOutputter := tx.NewMapDataOutputter(customLogData)
	var nostrEvent *tx.NostrEvent
	nostrEvent, err = tx.CreateTxLogEvent(customDataOutputter, "your_private_key_here")
	if err != nil {
		log.Fatalf("Failed to create Nostr event: %v", err)
	}

	fmt.Printf("   Event Kind: %d\n", nostrEvent.Kind)
	fmt.Printf("   Event Created At: %d\n", nostrEvent.CreatedAt)
	fmt.Printf("   Event Tags Count: %d\n", len(nostrEvent.Tags))

	// Print tags
	fmt.Println("   Event Tags:")
	for _, tag := range nostrEvent.Tags {
		if len(tag) >= 2 {
			fmt.Printf("     %s: %s\n", tag[0], tag[1])
		}
	}

	// Example 3: Update log status and create update event
	fmt.Println("\n3. Updating Log Status:")
	var updateEvent *tx.NostrEvent
	updateEvent, err = tx.UpdateLogStatusEvent(customLogData, "confirmed", "your_private_key_here")
	if err != nil {
		log.Fatalf("Failed to create update event: %v", err)
	}

	fmt.Printf("   Update Event Kind: %d\n", updateEvent.Kind)
	fmt.Printf("   Update Event Tags Count: %d\n", len(updateEvent.Tags))

	// Example 4: Parse event back
	fmt.Println("\n4. Parsing Event Back:")
	var parsedEvent *tx.TxLogEvent
	parsedEvent, err = tx.ParseTxLogEvent(nostrEvent)
	if err != nil {
		log.Fatalf("Failed to parse event: %v", err)
	}

	fmt.Printf("   Parsed Event Type: %s\n", parsedEvent.EventType)
	fmt.Printf("   Parsed Log Hash: %v\n", parsedEvent.LogData["hash"])
	fmt.Printf("   Parsed Log Status: %v\n", parsedEvent.LogData["status"])

	// Example 5: JSON representation
	fmt.Println("\n5. JSON Representation:")
	var eventJSON []byte
	eventJSON, err = json.MarshalIndent(nostrEvent, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal event: %v", err)
	}
	fmt.Printf("   Nostr Event JSON:\n%s\n", string(eventJSON))

	// Example 6: Check if it's an ERC20 transfer

	// Example 2: Dynamic data flattening with different event types
	fmt.Println("\n=== Example 2: Dynamic Data Flattening ===")

	// ERC20 Transfer event
	erc20Data := map[string]interface{}{
		"hash":   "0xerc20hash1234567890abcdef",
		"status": "pending",
		"data": map[string]interface{}{
			"topic": "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
			"from":  "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
			"to":    "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7",
			"value": "1000000000000000000",
			"token": "0x1234567890123456789012345678901234567890",
		},
	}

	var erc20Event *tx.NostrEvent
	erc20Event, err = tx.CreateTxLogEvent(tx.NewMapDataOutputter(erc20Data), "private_key")
	if err != nil {
		fmt.Printf("Error creating ERC20 event: %v\n", err)
	} else {
		fmt.Println("ERC20 Transfer Event Tags:")
		for _, tag := range erc20Event.Tags {
			if len(tag) >= 2 {
				fmt.Printf("  %s: %s\n", tag[0], tag[1])
			}
		}
	}

	// NFT Transfer event
	nftData := map[string]interface{}{
		"hash":   "0xnfthash1234567890abcdef",
		"status": "confirmed",
		"data": map[string]interface{}{
			"topic":    "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
			"from":     "0x1111111111111111111111111111111111111111",
			"to":       "0x2222222222222222222222222222222222222222",
			"tokenId":  "12345",
			"contract": "0x9876543210987654321098765432109876543210",
		},
	}

	var nftEvent *tx.NostrEvent
	nftEvent, err = tx.CreateTxLogEvent(tx.NewMapDataOutputter(nftData), "private_key")
	if err != nil {
		fmt.Printf("Error creating NFT event: %v\n", err)
	} else {
		fmt.Println("\nNFT Transfer Event Tags:")
		for _, tag := range nftEvent.Tags {
			if len(tag) >= 2 {
				fmt.Printf("  %s: %s\n", tag[0], tag[1])
			}
		}
	}

	// Complex nested data event
	complexData := map[string]interface{}{
		"hash":   "0xcomplexhash1234567890abcdef",
		"status": "pending",
		"data": map[string]interface{}{
			"topic": "0xcomplexhash1234567890abcdef",
			"metadata": map[string]interface{}{
				"name":     "Test Token",
				"symbol":   "TEST",
				"decimals": int64(18),
			},
			"participants": []interface{}{
				"0x3333333333333333333333333333333333333333",
				"0x4444444444444444444444444444444444444444",
			},
			"amount": "500000000000000000",
			"owner":  "0x5555555555555555555555555555555555555555",
		},
	}

	var complexEvent *tx.NostrEvent
	complexEvent, err = tx.CreateTxLogEvent(tx.NewMapDataOutputter(complexData), "private_key")
	if err != nil {
		fmt.Printf("Error creating complex event: %v\n", err)
	} else {
		fmt.Println("\nComplex Event Tags:")
		for _, tag := range complexEvent.Tags {
			if len(tag) >= 2 {
				fmt.Printf("  %s: %s\n", tag[0], tag[1])
			}
		}
	}

	fmt.Println("\n=== Example Complete ===")
}
