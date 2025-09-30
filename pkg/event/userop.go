package event

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/comunifi/nostr-eth/pkg/neth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nbd-wtf/go-nostr"
)

// NostrEventType represents the type of Nostr event for user operations
const (
	EventUserOpKind = 111001

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
	UserOpData neth.UserOp      `json:"user_op_data"`
	Paymaster  *common.Address  `json:"paymaster,omitempty"`
	EntryPoint *common.Address  `json:"entry_point,omitempty"`
	Data       *json.RawMessage `json:"data,omitempty"`
	TxHash     *string          `json:"tx_hash,omitempty"`
	EventType  EventTypeUserOp  `json:"event_type"`
	RetryCount int              `json:"retry_count,omitempty"`
	Tags       []string         `json:"tags,omitempty"`
}

func (u *UserOpEvent) MarshalJSON() ([]byte, error) {
	// Marshal UserOpData using its custom marshaling
	userOpDataBytes, err := json.Marshal(u.UserOpData)
	if err != nil {
		return nil, err
	}

	// Create a temporary struct that embeds all fields but overrides UserOpData
	type Alias UserOpEvent
	return json.Marshal(&struct {
		UserOpData json.RawMessage `json:"user_op_data"`
		*Alias
	}{
		UserOpData: userOpDataBytes,
		Alias:      (*Alias)(u),
	})
}

func (u *UserOpEvent) UnmarshalJSON(data []byte) error {
	// Create a temporary struct to handle the custom unmarshaling
	type Alias UserOpEvent
	aux := &struct {
		UserOpData json.RawMessage `json:"user_op_data"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Unmarshal the UserOpData using the custom unmarshaling from neth.UserOp
	if err := json.Unmarshal(aux.UserOpData, &u.UserOpData); err != nil {
		return err
	}

	return nil
}

// CreateUserOpEvent creates a new Nostr event for a user operation
func CreateUserOpEvent(chainID *big.Int, paymaster, entryPoint *common.Address, data *json.RawMessage, txHash *string, retryCount int, userOp neth.UserOp, eventType EventTypeUserOp) (*nostr.Event, error) {
	// Create the event data
	eventData := UserOpEvent{
		UserOpData: userOp,
		Paymaster:  paymaster,
		EntryPoint: entryPoint,
		Data:       data,
		TxHash:     txHash,
		EventType:  eventType,
		RetryCount: retryCount,
		Tags:       []string{"user_op", "user_op_0_0_6", "evm", chainID.String(), "account_abstraction"},
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
		Kind:      EventUserOpKind, // Custom kind for user operations
		Tags:      make([]nostr.Tag, 0),
		Content:   string(content),
	}

	// Add tags for better indexing and filtering
	evt.Tags = append(evt.Tags, []string{"d", userOp.GetHash(chainID)}) // Identifier using sender address

	// Type and category tags
	evt.Tags = append(evt.Tags, []string{"t", "user_op"})             // Type
	evt.Tags = append(evt.Tags, []string{"t", "user_op_0_0_6"})       // Version
	evt.Tags = append(evt.Tags, []string{"network", "evm"})           // Blockchain
	evt.Tags = append(evt.Tags, []string{"t", "account_abstraction"}) // AA specific

	// Chain-specific tag
	evt.Tags = append(evt.Tags, []string{"layer", chainID.String()}) // Chain ID

	// Paymaster tag if present
	if paymaster != nil {
		evt.Tags = append(evt.Tags, []string{"paymaster", paymaster.Hex()})
	}

	// Entry point tag if present
	if entryPoint != nil {
		evt.Tags = append(evt.Tags, []string{"entry_point", entryPoint.Hex()})
	}

	// Tx hash tag if present
	if txHash != nil {
		evt.Tags = append(evt.Tags, []string{"t", *txHash})
	}

	// Sender address tag
	evt.Tags = append(evt.Tags, []string{"p", userOp.Sender.String()}) // Sender address

	// Nonce tag for ordering
	evt.Tags = append(evt.Tags, []string{"nonce", userOp.Nonce.String()})

	// Alt tag
	alt := fmt.Sprintf("This is a new user operation request on chain %s", chainID.String())
	if paymaster != nil {
		alt += fmt.Sprintf("\n this is intended for processing by paymaster: %s", paymaster.Hex())
	}

	evt.Tags = append(evt.Tags, []string{"alt", alt})

	return evt, nil
}

// UpdateUserOpEvent creates a Nostr event for updating a user operation status
func UpdateUserOpEvent(chainID *big.Int, userOp neth.UserOp, txHash *string, retryCount int, eventType EventTypeUserOp, event *nostr.Event) (*nostr.Event, error) {

	userOpEvent, err := ParseUserOpEvent(event)
	if err != nil {
		return nil, err
	}

	// Create the event data
	eventData := UserOpEvent{
		UserOpData: userOp,
		Paymaster:  userOpEvent.Paymaster,
		EntryPoint: userOpEvent.EntryPoint,
		Data:       userOpEvent.Data,
		TxHash:     txHash,
		EventType:  eventType,
		RetryCount: retryCount,
		Tags:       []string{"user_op", "user_op_0_0_6", "evm", chainID.String(), "account_abstraction", "update"},
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
		Kind:      EventUserOpKind, // Custom kind for user operations
		Tags:      make([]nostr.Tag, 0),
		Content:   string(content),
	}

	// Add tags for better indexing and filtering
	evt.Tags = append(evt.Tags, []string{"d", userOp.GetHash(chainID)}) // Identifier using sender address

	// Type and category tags
	evt.Tags = append(evt.Tags, []string{"t", "user_op"})             // Type
	evt.Tags = append(evt.Tags, []string{"t", "user_op_0_0_6"})       // Version
	evt.Tags = append(evt.Tags, []string{"network", "evm"})           // Blockchain
	evt.Tags = append(evt.Tags, []string{"t", "account_abstraction"}) // AA specific

	// Chain-specific tag
	evt.Tags = append(evt.Tags, []string{"layer", chainID.String()}) // Chain ID

	// Paymaster tag if present
	if userOpEvent.Paymaster != nil {
		evt.Tags = append(evt.Tags, []string{"paymaster", userOpEvent.Paymaster.Hex()})
	}

	// Entry point tag if present
	if userOpEvent.EntryPoint != nil {
		evt.Tags = append(evt.Tags, []string{"entry_point", userOpEvent.EntryPoint.Hex()})
	}

	// Tx hash tag if present
	if userOpEvent.TxHash != nil {
		evt.Tags = append(evt.Tags, []string{"t", *userOpEvent.TxHash})
	}

	// Sender address tag
	evt.Tags = append(evt.Tags, []string{"p", userOp.Sender.String()}) // Sender address

	// Nonce tag for ordering
	evt.Tags = append(evt.Tags, []string{"nonce", userOp.Nonce.String()})

	// Alt tag
	alt := fmt.Sprintf("This is a user operation update of type %s on chain %s", eventType, chainID.String())
	if userOpEvent.Paymaster != nil {
		alt += fmt.Sprintf("\n this is intended for processing by paymaster: %s", userOpEvent.Paymaster.Hex())
	}

	evt.Tags = append(evt.Tags, []string{"alt", alt})

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
