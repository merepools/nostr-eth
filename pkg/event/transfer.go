package event

import (
	"encoding/json"
	"fmt"

	"github.com/comunifi/nostr-eth/pkg/neth"
	"github.com/nbd-wtf/go-nostr"
)

// NostrEventType represents the type of Nostr event for transaction logs
const (
	KindTxTransfer = 9735

	EventTypeTxTransferCreated EventTypeTxTransfer = "tx_transfer_created"
)

type EventTypeTxTransfer string

// TxLogEvent represents a Nostr event for transaction logs
type TxTransferEvent struct {
	LogData   neth.Log            `json:"log_data"`
	EventType EventTypeTxTransfer `json:"event_type"`
	Tags      []string            `json:"tags,omitempty"`
}

// CreateTxTransferEvent creates a new Nostr event for a transfer
func CreateTxTransferEvent(log neth.Log) (*nostr.Event, error) {
	if log.Topic != neth.TopicERC20Transfer {
		return nil, fmt.Errorf("topic is not an ERC20 transfer")
	}

	// Create the event data
	eventData := TxTransferEvent{
		LogData:   log,
		EventType: EventTypeTxTransferCreated,
		Tags:      []string{"tx_transfer", "evm", log.ChainID},
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
		Kind:      KindTxTransfer, // Custom kind for transaction logs
		Tags:      make([]nostr.Tag, 0),
		Content:   string(content),
	}

	// Add tags for better indexing and filtering
	evt.Tags = append(evt.Tags, []string{"d", log.Hash}) // Identifier

	// Type and category tags
	evt.Tags = append(evt.Tags, []string{"t", "tx_transfer"}) // Type
	evt.Tags = append(evt.Tags, []string{"network", "evm"})   // Blockchain

	// Chain-specific tag
	evt.Tags = append(evt.Tags, []string{"layer", log.ChainID}) // Chain ID

	// Reference tags for transaction hash
	evt.Tags = append(evt.Tags, []string{"r", log.TxHash}) // Transaction hash as reference

	if log.Data != nil {
		var data map[string]interface{}
		err := json.Unmarshal(*log.Data, &data)
		if err != nil {
			return nil, err
		}

		sender, ok := data[neth.DataKeyFrom].(string)
		if !ok {
			return nil, fmt.Errorf("from is not a string")
		}

		to, ok := data[neth.DataKeyTo].(string)
		if !ok {
			return nil, fmt.Errorf("to is not a string")
		}

		amount, ok := data[neth.DataKeyValue].(string)
		if !ok {
			return nil, fmt.Errorf("amount is not a string")
		}

		evt.Tags = append(evt.Tags, []string{"P", sender}) // Sender address
		evt.Tags = append(evt.Tags, []string{"p", to})     // Recipient/Contract address

		evt.Tags = append(evt.Tags, []string{"amount", amount}) // Amount
	}

	// Topic tag
	evt.Tags = append(evt.Tags, []string{"t", log.Topic})

	// Contract address tag
	evt.Tags = append(evt.Tags, []string{"t", log.To})

	// Flatten data into tags
	dataTags := []nostr.Tag{}
	if log.Data != nil {
		dataTags = flattenDataToTags(*log.Data)
		evt.Tags = append(evt.Tags, dataTags...)
	}

	// Alt tag
	alt := fmt.Sprintf("This is an evm transaction log for topic %s on chain %s", log.Topic, log.ChainID)
	if len(dataTags) > 0 {
		alt += "\n Data:"
	}
	for _, tag := range dataTags {
		alt += fmt.Sprintf("\n %s: %s", tag[0], tag[1])
	}

	evt.Tags = append(evt.Tags, []string{"alt", alt})

	return evt, nil
}

// ParseTxTransferEvent parses a Nostr event back into a TxTransferEvent
func ParseTxTransferEvent(evt *nostr.Event) (*TxTransferEvent, error) {
	var txTransferEvent TxTransferEvent
	err := json.Unmarshal([]byte(evt.Content), &txTransferEvent)
	if err != nil {
		return nil, err
	}
	return &txTransferEvent, nil
}
