package nostreth

// Re-export all functions from the log package
import (
	"math/big"

	"github.com/comunifi/nostr-eth/pkg/event"
	"github.com/comunifi/nostr-eth/pkg/neth"
	"github.com/nbd-wtf/go-nostr"
)

// Re-export log package types
type TxLogEvent = event.TxLogEvent
type Log = neth.Log
type UserOpEvent = event.UserOpEvent
type UserOp = neth.UserOp

// Re-export log package constants
const (
	EventTypeTxLogCreated    = event.EventTypeTxLogCreated
	EventTypeTxLogUpdated    = event.EventTypeTxLogUpdated
	EventTypeUserOpRequested = event.EventTypeUserOpRequested
	EventTypeUserOpExecuted  = event.EventTypeUserOpExecuted
	EventTypeUserOpConfirmed = event.EventTypeUserOpConfirmed
	EventTypeUserOpExpired   = event.EventTypeUserOpExpired
	EventTypeUserOpFailed    = event.EventTypeUserOpFailed
)

// Re-export log package functions
func CreateTxLogEvent(log neth.Log) (*nostr.Event, error) {
	return event.CreateTxLogEvent(log)
}

func UpdateTxLogEvent(log neth.Log, originalEventID ...string) (*nostr.Event, error) {
	return event.UpdateTxLogEvent(log, originalEventID...)
}

func ParseTxLogEvent(evt *nostr.Event) (*event.TxLogEvent, error) {
	return event.ParseTxLogEvent(evt)
}

func GetEventData(log neth.Log) (map[string]interface{}, error) {
	return log.GetEventData()
}

func RequestUserOpEvent(chainID *big.Int, userOp neth.UserOp) (*nostr.Event, error) {
	return event.RequestUserOpEvent(chainID, userOp)
}

func UpdateUserOpEvent(chainID *big.Int, userOp neth.UserOp, eventType event.EventTypeUserOp, originalEventID ...string) (*nostr.Event, error) {
	return event.UpdateUserOpEvent(chainID, userOp, eventType, originalEventID...)
}

func ParseUserOpEvent(evt *nostr.Event) (*event.UserOpEvent, error) {
	return event.ParseUserOpEvent(evt)
}
