package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	nostreth "github.com/comunifi/nostr-eth"
	"github.com/nbd-wtf/go-nostr"
)

func main() {
	fmt.Println("=== Nostr nostreth.ereum Transaction Log Example ===")

	var err error

	data, err := json.Marshal(map[string]interface{}{
		"topic": "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
		"from":  "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
		"to":    "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7",
		"value": "1000000000000000000",
	})
	if err != nil {
		log.Fatalf("Failed to marshal data: %v", err)
	}

	var dataJSON json.RawMessage
	dataJSON = data

	// Example 1: Create a generic transaction log event
	fmt.Println("1. Creating Generic Transaction Log Event:")
	genericLogData := nostreth.Log{
		Hash:      "0x1234567890abcdef1234567890abcdef12345678",
		TxHash:    "0xabcdef1234567890abcdef1234567890abcdef12",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Nonce:     12345,
		Sender:    "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
		To:        "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7",
		Value:     big.NewInt(1000000000000000000),
		Data:      &dataJSON,
	}

	var genericEvent *nostr.Event
	genericEvent, err = nostreth.CreateTxLogEvent(genericLogData)
	if err != nil {
		log.Fatalf("Failed to create generic transaction event: %v", err)
	}

	// Parse the event to get the log data
	var txLogEvent nostreth.TxLogEvent
	err = json.Unmarshal([]byte(genericEvent.Content), &txLogEvent)
	if err != nil {
		log.Fatalf("Failed to parse event content: %v", err)
	}

	logData := txLogEvent.LogData
	fmt.Printf("   Log Hash: %v\n", logData.Hash)
	fmt.Printf("   Transaction Hash: %v\n", logData.TxHash)
	fmt.Printf("   From: %v\n", logData.Sender)
	fmt.Printf("   To: %v\n", logData.To)

	// Get transfer data
	transferData, err := logData.GetEventData()
	if err != nil {
		log.Fatalf("Failed to get transfer data: %v", err)
	}
	if transferData != nil {
		fmt.Printf("   Transfer Value: %v\n", transferData["value"])
		fmt.Printf("   Transfer Topic: %v\n", transferData["topic"])
	}

	// Example 2: Create a Nostr event for custom log data
	fmt.Println("\n2. Creating Nostr Event for Custom Log Data:")
	customLogData := nostreth.Log{
		Hash:      "0x9876543210fedcba9876543210fedcba98765432",
		TxHash:    "0xfedcba0987654321fedcba0987654321fedcba09",
		CreatedAt: time.Now(), // Unix timestamp
		UpdatedAt: time.Now(),
		Nonce:     67890,
		Sender:    "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b8",
		To:        "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b9",
		Value:     big.NewInt(1000000000000000000), // 1 nostreth.
		Data:      nil,
	}

	var nostrEvent *nostr.Event
	nostrEvent, err = nostreth.CreateTxLogEvent(customLogData)
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
	var updateEvent *nostr.Event
	updateEvent, err = nostreth.UpdateTxLogEvent(customLogData, &nostr.Event{ID: "event_id"})
	if err != nil {
		log.Fatalf("Failed to create update event: %v", err)
	}

	fmt.Printf("   Update Event Kind: %d\n", updateEvent.Kind)
	fmt.Printf("   Update Event Tags Count: %d\n", len(updateEvent.Tags))

	// Example 4: Parse event back
	fmt.Println("\n4. Parsing Event Back:")
	var parsedEvent *nostreth.TxLogEvent
	parsedEvent, err = nostreth.ParseTxLogEvent(nostrEvent)
	if err != nil {
		log.Fatalf("Failed to parse event: %v", err)
	}

	fmt.Printf("   Parsed Event Type: %s\n", parsedEvent.EventType)
	fmt.Printf("   Parsed Log Hash: %v\n", parsedEvent.LogData.Hash)

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

	data2, err := json.Marshal(map[string]interface{}{
		"topic": "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
		"from":  "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
		"to":    "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7",
		"value": "1000000000000000000",
		"token": "0x1234567890123456789012345678901234567890",
	})
	if err != nil {
		log.Fatalf("Failed to marshal data: %v", err)
	}

	var data2JSON json.RawMessage
	data2JSON = data2

	// ERC20 Transfer event
	erc20Data := nostreth.Log{
		Hash: "0xerc20hash1234567890abcdef",
		Data: &data2JSON,
	}

	var erc20Event *nostr.Event
	erc20Event, err = nostreth.CreateTxLogEvent(erc20Data)
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

	data3, err := json.Marshal(map[string]interface{}{
		"topic":    "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
		"from":     "0x1111111111111111111111111111111111111111",
		"to":       "0x2222222222222222222222222222222222222222",
		"tokenId":  "12345",
		"contract": "0x9876543210987654321098765432109876543210",
	})
	if err != nil {
		log.Fatalf("Failed to marshal data: %v", err)
	}

	var data3JSON json.RawMessage
	data3JSON = data3

	// NFT Transfer event
	nftData := nostreth.Log{
		Hash: "0xnfthash1234567890abcdef",
		Data: &data3JSON,
	}

	var nftEvent *nostr.Event
	nftEvent, err = nostreth.CreateTxLogEvent(nftData)
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

	data4, err := json.Marshal(map[string]interface{}{
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
	})
	if err != nil {
		log.Fatalf("Failed to marshal data: %v", err)
	}

	var data4JSON json.RawMessage
	data4JSON = data4

	// Complex nested data event
	complexData := nostreth.Log{
		Hash: "0xcomplexhash1234567890abcdef",
		Data: &data4JSON,
	}

	var complexEvent *nostr.Event
	complexEvent, err = nostreth.CreateTxLogEvent(complexData)
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

	// Example 7: NIP-29 Group Implementation
	fmt.Println("\n=== NIP-29 Group Implementation Example ===")

	// Example group ID and host
	groupID := "example-group-123"
	host := "wss://groups.example.com"
	groupIdentifier := nostreth.FormatGroupIdentifier(host, groupID)

	fmt.Printf("Group Identifier: %s\n\n", groupIdentifier)

	// Create Group Metadata Event (kind 39000)
	fmt.Println("1. Creating Group Metadata Event...")
	admins := []string{"admin1@example.com", "admin2@example.com"}
	moderators := []string{"mod1@example.com", "mod2@example.com"}

	metadataEvent, err := nostreth.CreateGroupMetadataEvent(
		groupID,
		"Example Group",
		"This is an example group for demonstrating NIP-29",
		"https://example.com/group-picture.jpg",
		admins,
		moderators,
		true,  // private
		false, // not closed
	)
	if err != nil {
		log.Fatalf("Failed to create group metadata event: %v", err)
	}

	// Set a dummy pubkey for demonstration
	metadataEvent.PubKey = "dummy-pubkey-for-demo"
	metadataEvent.ID = "metadata-event-id"

	fmt.Printf("   Group Metadata Event:\n")
	fmt.Printf("   - Kind: %d\n", metadataEvent.Kind)
	fmt.Printf("   - Group ID: %s\n", groupID)
	fmt.Printf("   - Content: %s\n", metadataEvent.Content)
	fmt.Printf("   - Tags: %v\n\n", metadataEvent.Tags)

	// Create Group Message Event (kind 39001)
	fmt.Println("2. Creating Group Message Event...")
	mentions := []string{"user1@example.com", "user2@example.com"}

	messageEvent, err := nostreth.CreateGroupMessageEvent(
		groupID,
		"Hello everyone! This is a test message in our group.",
		"", // not a reply
		mentions,
	)
	if err != nil {
		log.Fatalf("Failed to create group message event: %v", err)
	}

	messageEvent.PubKey = "dummy-pubkey-for-demo"
	messageEvent.ID = "message-event-id"

	fmt.Printf("   Group Message Event:\n")
	fmt.Printf("   - Kind: %d\n", messageEvent.Kind)
	fmt.Printf("   - Group ID: %s\n", groupID)
	fmt.Printf("   - Content: %s\n", messageEvent.Content)
	fmt.Printf("   - Tags: %v\n\n", messageEvent.Tags)

	// Create Group Join Event (kind 39002)
	fmt.Println("3. Creating Group Join Event...")
	joinEvent, err := nostreth.CreateGroupJoinEvent(
		groupID,
		"newuser@example.com",
		"member",
	)
	if err != nil {
		log.Fatalf("Failed to create group join event: %v", err)
	}

	joinEvent.PubKey = "dummy-pubkey-for-demo"
	joinEvent.ID = "join-event-id"

	fmt.Printf("   Group Join Event:\n")
	fmt.Printf("   - Kind: %d\n", joinEvent.Kind)
	fmt.Printf("   - Group ID: %s\n", groupID)
	fmt.Printf("   - Content: %s\n", joinEvent.Content)
	fmt.Printf("   - Tags: %v\n\n", joinEvent.Tags)

	// Create Group Leave Event (kind 39003)
	fmt.Println("4. Creating Group Leave Event...")
	leaveEvent, err := nostreth.CreateGroupLeaveEvent(
		groupID,
		"leavinguser@example.com",
		"Personal reasons",
	)
	if err != nil {
		log.Fatalf("Failed to create group leave event: %v", err)
	}

	leaveEvent.PubKey = "dummy-pubkey-for-demo"
	leaveEvent.ID = "leave-event-id"

	fmt.Printf("   Group Leave Event:\n")
	fmt.Printf("   - Kind: %d\n", leaveEvent.Kind)
	fmt.Printf("   - Group ID: %s\n", groupID)
	fmt.Printf("   - Content: %s\n", leaveEvent.Content)
	fmt.Printf("   - Tags: %v\n\n", leaveEvent.Tags)

	// Create Group Moderation Event (kind 39004)
	fmt.Println("5. Creating Group Moderation Event...")
	moderationEvent, err := nostreth.CreateGroupModerationEvent(
		groupID,
		"ban",
		"spammer@example.com",
		"Spamming the group",
		86400, // 24 hours in seconds
	)
	if err != nil {
		log.Fatalf("Failed to create group moderation event: %v", err)
	}

	moderationEvent.PubKey = "dummy-pubkey-for-demo"
	moderationEvent.ID = "moderation-event-id"

	fmt.Printf("   Group Moderation Event:\n")
	fmt.Printf("   - Kind: %d\n", moderationEvent.Kind)
	fmt.Printf("   - Group ID: %s\n", groupID)
	fmt.Printf("   - Content: %s\n", moderationEvent.Content)
	fmt.Printf("   - Tags: %v\n\n", moderationEvent.Tags)

	// Parse and validate events
	fmt.Println("6. Parsing and Validating Events...")

	// Parse metadata event
	metadata, err := nostreth.ParseGroupMetadataEvent(metadataEvent)
	if err != nil {
		log.Fatalf("Failed to parse group metadata: %v", err)
	}
	fmt.Printf("   Parsed Group Metadata:\n")
	fmt.Printf("   - Name: %s\n", metadata.Name)
	fmt.Printf("   - About: %s\n", metadata.About)
	fmt.Printf("   - Private: %t\n", metadata.Private)
	fmt.Printf("   - Admins: %v\n", metadata.Admins)
	fmt.Printf("   - Moderators: %v\n\n", metadata.Moderators)

	// Parse message event
	message, err := nostreth.ParseGroupMessageEvent(messageEvent)
	if err != nil {
		log.Fatalf("Failed to parse group message: %v", err)
	}
	fmt.Printf("   Parsed Group Message:\n")
	fmt.Printf("   - Content: %s\n", message.Content)
	fmt.Printf("   - Mentions: %v\n\n", message.Mentions)

	// Utility functions demonstration
	fmt.Println("7. Utility Functions Demonstration...")

	// Check if events are group events
	fmt.Printf("   Is metadata event a group event: %t\n", nostreth.IsGroupEvent(metadataEvent))
	fmt.Printf("   Is message event a group event: %t\n", nostreth.IsGroupEvent(messageEvent))

	// Get group ID from events
	metadataGroupID, err := nostreth.GetGroupIDFromEvent(metadataEvent)
	if err != nil {
		log.Fatalf("Failed to get group ID from metadata event: %v", err)
	}
	fmt.Printf("   Group ID from metadata event: %s\n", metadataGroupID)

	// Get event type
	fmt.Printf("   Event type from metadata event: %s\n", nostreth.GetEventTypeFromGroupEvent(metadataEvent))
	fmt.Printf("   Event type from message event: %s\n", nostreth.GetEventTypeFromGroupEvent(messageEvent))

	// Validate group ID
	err = nostreth.ValidateGroupID(groupID)
	if err != nil {
		log.Fatalf("Group ID validation failed: %v", err)
	}
	fmt.Printf("   Group ID validation: PASSED\n")

	// Parse group identifier
	parsedHost, parsedGroupID, err := nostreth.ParseGroupIdentifier(groupIdentifier)
	if err != nil {
		log.Fatalf("Failed to parse group identifier: %v", err)
	}
	fmt.Printf("   Parsed group identifier - Host: %s, Group ID: %s\n\n", parsedHost, parsedGroupID)

	// Filter events by group
	fmt.Println("8. Filtering Events by Group...")
	allEvents := []*nostr.Event{metadataEvent, messageEvent, joinEvent, leaveEvent, moderationEvent}
	groupEvents := nostreth.FilterGroupEventsByGroupID(allEvents, groupID)
	fmt.Printf("   Total events: %d\n", len(allEvents))
	fmt.Printf("   Events in group '%s': %d\n", groupID, len(groupEvents))

	for i, evt := range groupEvents {
		eventType := nostreth.GetEventTypeFromGroupEvent(evt)
		fmt.Printf("   Event %d: %s (kind %d)\n", i+1, eventType, evt.Kind)
	}

	fmt.Println("\n✅ NIP-29 Group Implementation Example Completed Successfully!")
	fmt.Println("\nKey Features Implemented:")
	fmt.Println("- Group metadata events (kind 39000)")
	fmt.Println("- Group message events (kind 39001)")
	fmt.Println("- Group join events (kind 39002)")
	fmt.Println("- Group leave events (kind 39003)")
	fmt.Println("- Group moderation events (kind 39004)")
	fmt.Println("- Event parsing and validation")
	fmt.Println("- Group filtering and utility functions")
	fmt.Println("- Group identifier parsing and formatting")
}
