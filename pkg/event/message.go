package event

import (
	"fmt"
	"strconv"
	"time"

	"github.com/nbd-wtf/go-nostr"
)

const (
	KindGenericRepost = 16
)

// CreateMessageEvent creates a new Nostr event for a simple text message
func CreateMessageEvent(content string, group *string) (*nostr.Event, error) {

	// Create the Nostr event with plain text content
	evt := &nostr.Event{
		PubKey:    "", // Will be derived from private key
		CreatedAt: nostr.Timestamp(time.Now().Unix()),
		Kind:      1, // Standard kind for text messages
		Tags:      make([]nostr.Tag, 0),
		Content:   content, // Plain text content
	}

	// Add tags for better indexing and filtering

	// Type and category tags
	evt.Tags = append(evt.Tags, []string{"t", "message"}) // Type
	evt.Tags = append(evt.Tags, []string{"t", "text"})    // Content type

	// Group tag for filtering by group (NIP-29 compliant)
	if group != nil {
		evt.Tags = append(evt.Tags, []string{"h", *group}) // Group ID
	}

	return evt, nil
}

// UpdateMessageEvent creates a Nostr event for updating a message
func UpdateMessageEvent(content string, group *string, originalEvent *nostr.Event) (*nostr.Event, error) {

	// Create the Nostr event with plain text content
	evt := &nostr.Event{
		PubKey:    "", // Will be derived from private key
		CreatedAt: nostr.Timestamp(time.Now().Unix()),
		Kind:      1, // Standard kind for text messages
		Tags:      make([]nostr.Tag, 0),
		Content:   content, // Plain text content
	}

	// Add reference to original event if provided (NIP-10 compliant)
	if originalEvent != nil {
		// Use marked e tag format: [event-id, relay-url, marker, pubkey]
		evt.Tags = append(evt.Tags, []string{"e", originalEvent.ID, "", "reply", originalEvent.PubKey})
	}

	// Add tags for better indexing and filtering

	// Type and category tags
	evt.Tags = append(evt.Tags, []string{"t", "message"}) // Type
	evt.Tags = append(evt.Tags, []string{"t", "text"})    // Content type
	evt.Tags = append(evt.Tags, []string{"t", "update"})  // Update marker

	// Group tag for filtering by group (NIP-29 compliant)
	if group != nil {
		evt.Tags = append(evt.Tags, []string{"h", *group}) // Group ID
	}

	return evt, nil
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

// CreateReplyEvent creates a NIP-10 compliant reply event
func CreateReplyEvent(content string, group *string, replyTo *nostr.Event) (*nostr.Event, error) {

	// Create the Nostr event with plain text content
	evt := &nostr.Event{
		PubKey:    "", // Will be derived from private key
		CreatedAt: nostr.Timestamp(time.Now().Unix()),
		Kind:      1, // Standard kind for text messages
		Tags:      make([]nostr.Tag, 0),
		Content:   content, // Plain text content
	}

	// Add NIP-10 compliant e tags for reply threading
	if replyTo != nil {
		// Find the root event in the thread
		rootEvent := findRootEvent(replyTo)

		// Add reply marker to the immediate parent
		evt.Tags = append(evt.Tags, []string{"e", replyTo.ID, "", "reply", replyTo.PubKey})

		// Add root marker if this is not the root
		if rootEvent.ID != replyTo.ID {
			evt.Tags = append(evt.Tags, []string{"e", rootEvent.ID, "", "root", rootEvent.PubKey})
		}

		// Add NIP-10 compliant p tags for participant tracking
		participants := getParticipantsFromEvent(replyTo)
		for participant := range participants {
			evt.Tags = append(evt.Tags, []string{"p", participant})
		}
	}

	// Add tags for better indexing and filtering

	// Type and category tags
	evt.Tags = append(evt.Tags, []string{"t", "message"}) // Type
	evt.Tags = append(evt.Tags, []string{"t", "text"})    // Content type
	evt.Tags = append(evt.Tags, []string{"t", "reply"})   // Reply marker

	// Group tag for filtering by group (NIP-29 compliant)
	if group != nil {
		evt.Tags = append(evt.Tags, []string{"h", *group}) // Group ID
	}

	return evt, nil
}

// findRootEvent finds the root event in a reply thread (NIP-10 compliant)
func findRootEvent(event *nostr.Event) *nostr.Event {
	// Look for root marker in e tags
	for _, tag := range event.Tags {
		if len(tag) >= 4 && tag[0] == "e" && tag[3] == "root" {
			// This event references a root, return the current event as it's part of the thread
			return event
		}
	}

	// If no root marker found, this might be the root or we need to traverse
	// For simplicity, return the current event
	// In a full implementation, you might want to traverse the thread
	return event
}

// getParticipantsFromEvent extracts all participants from an event's p tags (NIP-10 compliant)
func getParticipantsFromEvent(event *nostr.Event) map[string]bool {
	participants := make(map[string]bool)

	if event == nil {
		return participants
	}

	for _, tag := range event.Tags {
		if len(tag) >= 2 && tag[0] == "p" {
			participants[tag[1]] = true
		}
	}

	return participants
}

// GetReplyChainFromEvent extracts the reply chain from an event (NIP-10 compliant)
func GetReplyChainFromEvent(evt *nostr.Event) (string, string, error) {
	var replyTo, root string

	for _, tag := range evt.Tags {
		if len(tag) >= 4 && tag[0] == "e" {
			switch tag[3] {
			case "reply":
				replyTo = tag[1]
			case "root":
				root = tag[1]
			}
		}
	}

	return replyTo, root, nil
}

// GetParticipantsFromEvent extracts all participants from an event (NIP-10 compliant)
func GetParticipantsFromEvent(evt *nostr.Event) []string {
	var participants []string
	seen := make(map[string]bool)

	for _, tag := range evt.Tags {
		if len(tag) >= 2 && tag[0] == "p" && !seen[tag[1]] {
			participants = append(participants, tag[1])
			seen[tag[1]] = true
		}
	}

	return participants
}

// CreateQuoteRepostEvent creates a NIP-18 compliant quote event
func CreateQuoteRepostEvent(content string, group *string, repostedEvent *nostr.Event, relayURL string) (*nostr.Event, error) {
	// Create the Nostr event with plain text content
	evt := &nostr.Event{
		PubKey:    "", // Will be derived from private key
		CreatedAt: nostr.Timestamp(time.Now().Unix()),
		Kind:      1, // Standard kind for text messages
		Tags:      make([]nostr.Tag, 0),
		Content:   content, // Plain text content
	}

	// Add NIP-18 compliant q tag for quote
	if repostedEvent != nil {
		evt.Tags = append(evt.Tags, []string{"k", strconv.Itoa(repostedEvent.Kind)})

		evt.Tags = append(evt.Tags, []string{"q", repostedEvent.ID, relayURL, repostedEvent.PubKey})
		// Add p tag for mentioned event author (NIP-10 compliant)
		evt.Tags = append(evt.Tags, []string{"p", repostedEvent.PubKey})

		// Encode the reposted event ID using NIP-19 nevent format
		nevent, err := EncodeEventIDToNevent(repostedEvent.ID, relayURL, repostedEvent.PubKey, repostedEvent.Kind)
		if err != nil {
			return nil, fmt.Errorf("failed to encode event ID to nevent: %v", err)
		}
		evt.Content += "\n nostr:" + nevent
	}

	// Add tags for better indexing and filtering

	// Type and category tags
	evt.Tags = append(evt.Tags, []string{"t", "message"}) // Type
	evt.Tags = append(evt.Tags, []string{"t", "text"})    // Content type
	evt.Tags = append(evt.Tags, []string{"t", "mention"}) // Mention marker

	// Group tag for filtering by group (NIP-29 compliant)
	if group != nil {
		evt.Tags = append(evt.Tags, []string{"h", *group}) // Group ID
	}

	return evt, nil
}

// IsReplyEvent checks if an event is a reply (NIP-10 compliant)
func IsReplyEvent(evt *nostr.Event) bool {
	for _, tag := range evt.Tags {
		if len(tag) >= 4 && tag[0] == "e" && tag[3] == "reply" {
			return true
		}
	}
	return false
}

// IsMentionEvent checks if an event is a mention (NIP-10 compliant)
func IsMentionEvent(evt *nostr.Event) bool {
	for _, tag := range evt.Tags {
		if len(tag) >= 4 && tag[0] == "e" && tag[3] == "mention" {
			return true
		}
	}
	return false
}

// IsRootEvent checks if an event is a root event (NIP-10 compliant)
func IsRootEvent(evt *nostr.Event) bool {
	// A root event has no reply markers pointing to it
	// and may have root markers pointing to itself
	hasReplyMarker := false
	for _, tag := range evt.Tags {
		if len(tag) >= 4 && tag[0] == "e" && tag[3] == "reply" {
			hasReplyMarker = true
			break
		}
	}
	return !hasReplyMarker
}
