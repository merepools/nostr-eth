package event

import (
	"encoding/json"
	"math/big"

	"github.com/comunifi/nostr-eth/pkg/neth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nbd-wtf/go-nostr"
)

// NostrEventType represents the type of Nostr event for user operations
const (
	EventUserOpKind = 30001

	EventTypeUserOpRequested EventTypeUserOp = "user_op_requested"
	EventTypeUserOpSigned    EventTypeUserOp = "user_op_signed"
	EventTypeUserOpSubmitted EventTypeUserOp = "user_op_submitted"
	EventTypeUserOpExecuted  EventTypeUserOp = "user_op_executed"
	EventTypeUserOpConfirmed EventTypeUserOp = "user_op_confirmed"
	EventTypeUserOpExpired   EventTypeUserOp = "user_op_expired"
	EventTypeUserOpFailed    EventTypeUserOp = "user_op_failed"
)

type EventTypeUserOp string

// UserOpEvent represents a Nostr event for user operations
type UserOpEvent struct {
	UserOpData neth.UserOp     `json:"user_op_data"`
	Paymaster  *common.Address `json:"paymaster,omitempty"`
	EventType  EventTypeUserOp `json:"event_type"`
	Tags       []string        `json:"tags,omitempty"`
}

// CreateUserOpEvent creates a new Nostr event for a user operation
func RequestUserOpEvent(chainID *big.Int, paymaster *common.Address, userOp neth.UserOp) (*nostr.Event, error) {
	// Create the event data
	eventData := UserOpEvent{
		UserOpData: userOp,
		Paymaster:  paymaster,
		EventType:  EventTypeUserOpRequested,
		Tags:       []string{"user_op", "user_op_0_0_6", "ethereum", "account_abstraction"},
	}

	// Marshal the event data
	content, err := json.Marshal(eventData)
	if err != nil {
		return nil, err
	}

	// Create the Nostr event
	evt := &nostr.Event{
		PubKey:    "", // Will be derived from private key
		CreatedAt: nostr.Now(),
		Kind:      30001, // Custom kind for user operations
		Tags:      make([]nostr.Tag, 0),
		Content:   string(content),
	}

	// Add tags for better indexing and filtering
	evt.Tags = append(evt.Tags, []string{"d", userOp.GetHash(chainID)}) // Identifier using sender address

	// Type and category tags
	evt.Tags = append(evt.Tags, []string{"t", "user_op"})             // Type
	evt.Tags = append(evt.Tags, []string{"t", "user_op_0_0_6"})       // Version
	evt.Tags = append(evt.Tags, []string{"t", "ethereum"})            // Blockchain
	evt.Tags = append(evt.Tags, []string{"t", "account_abstraction"}) // AA specific

	// Chain-specific tag
	evt.Tags = append(evt.Tags, []string{"t", chainID.String()}) // Chain ID

	// Paymaster tag if present
	if paymaster != nil {
		evt.Tags = append(evt.Tags, []string{"paymaster", paymaster.Hex()})
	}

	// Sender address tag
	evt.Tags = append(evt.Tags, []string{"p", userOp.Sender.String()}) // Sender address

	// Nonce tag for ordering
	evt.Tags = append(evt.Tags, []string{"nonce", userOp.Nonce.String()})

	return evt, nil
}

// UpdateUserOpEvent creates a Nostr event for updating a user operation status
func UpdateUserOpEvent(chainID *big.Int, userOp neth.UserOp, eventType EventTypeUserOp, event *nostr.Event) (*nostr.Event, error) {

	userOpEvent, err := ParseUserOpEvent(event)
	if err != nil {
		return nil, err
	}

	// Create the event data
	eventData := UserOpEvent{
		UserOpData: userOp,
		Paymaster:  userOpEvent.Paymaster,
		EventType:  eventType,
		Tags:       []string{"user_op", "user_op_0_0_6", "ethereum", "account_abstraction", "update"},
	}

	// Marshal the event data
	content, err := json.Marshal(eventData)
	if err != nil {
		return nil, err
	}

	// Create the Nostr event
	evt := &nostr.Event{
		PubKey:    "", // Will be derived from private key
		CreatedAt: nostr.Now(),
		Kind:      30001, // Custom kind for user operations
		Tags:      make([]nostr.Tag, 0),
		Content:   string(content),
	}

	// Add tags for better indexing and filtering
	evt.Tags = append(evt.Tags, []string{"d", userOp.GetHash(chainID)}) // Identifier using sender address

	// Type and category tags
	evt.Tags = append(evt.Tags, []string{"t", "user_op"})             // Type
	evt.Tags = append(evt.Tags, []string{"t", "user_op_0_0_6"})       // Version
	evt.Tags = append(evt.Tags, []string{"t", "ethereum"})            // Blockchain
	evt.Tags = append(evt.Tags, []string{"t", "account_abstraction"}) // AA specific

	// Chain-specific tag
	evt.Tags = append(evt.Tags, []string{"t", chainID.String()}) // Chain ID

	// Paymaster tag if present
	if userOpEvent.Paymaster != nil {
		evt.Tags = append(evt.Tags, []string{"paymaster", userOpEvent.Paymaster.Hex()})
	}

	// Sender address tag
	evt.Tags = append(evt.Tags, []string{"p", userOp.Sender.String()}) // Sender address

	// Nonce tag for ordering
	evt.Tags = append(evt.Tags, []string{"nonce", userOp.Nonce.String()})

	return evt, nil
}

// ParseUserOpEvent parses a Nostr event back into a UserOpEvent
func ParseUserOpEvent(evt *nostr.Event) (*UserOpEvent, error) {
	var userOpEvent UserOpEvent
	err := json.Unmarshal([]byte(evt.Content), &userOpEvent)
	if err != nil {
		return nil, err
	}
	return &userOpEvent, nil
}

// isKnownFunctionSignature checks if the function signature is one of the known ones
func isKnownFunctionSignature(sig []byte) bool {
	knownSigs := [][]byte{
		neth.FuncSigSingle,
		neth.FuncSigBatch,
		neth.FuncSigSafeExecFromModule,
	}

	for _, knownSig := range knownSigs {
		if len(sig) == len(knownSig) {
			match := true
			for i, b := range sig {
				if b != knownSig[i] {
					match = false
					break
				}
			}
			if match {
				return true
			}
		}
	}
	return false
}
