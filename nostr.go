package nostreth

// Re-export all functions from the log package
import (
	"github.com/citizenapp2/nostr-eth/pkg/event"
	"github.com/citizenapp2/nostr-eth/pkg/neth"
	"github.com/nbd-wtf/go-nostr"
)

// Re-export log package types
type TxLogEvent = event.TxLogEvent
type Log = neth.Log

// Re-export log package constants
const (
	EventTypeTxLogCreated = event.EventTypeTxLogCreated
	EventTypeTxLogUpdated = event.EventTypeTxLogUpdated
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
