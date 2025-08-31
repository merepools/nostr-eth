package eth

// Re-export all functions from the log package
import (
	log "github.com/citizenwallet/nostr-eth/pkg/eth"
	"github.com/nbd-wtf/go-nostr"
)

// Re-export log package types
type TxLogEvent = log.TxLogEvent
type JSONOutputter = log.JSONOutputter

// Re-export log package constants
const (
	EventTypeTxLogCreated = log.EventTypeTxLogCreated
	EventTypeTxLogUpdated = log.EventTypeTxLogUpdated
)

// Re-export log package functions
func CreateTxLogEvent(logData log.JSONOutputter, privateKey string) (*nostr.Event, error) {
	return log.CreateTxLogEvent(logData, privateKey)
}

func UpdateTxLogEvent(logData map[string]interface{}, privateKey string, originalEventID ...string) (*nostr.Event, error) {
	return log.UpdateTxLogEvent(logData, privateKey, originalEventID...)
}

func ParseTxLogEvent(evt *nostr.Event) (*log.TxLogEvent, error) {
	return log.ParseTxLogEvent(evt)
}

func UpdateLogStatusEvent(logData map[string]interface{}, newStatus string, privateKey string, originalEventID ...string) (*nostr.Event, error) {
	return log.UpdateLogStatusEvent(logData, newStatus, privateKey, originalEventID...)
}

func GetTransferData(logData map[string]interface{}) (map[string]interface{}, error) {
	return log.GetTransferData(logData)
}

func NewGenericJSONOutputter(data map[string]interface{}) log.GenericJSONOutputter {
	return log.NewGenericJSONOutputter(data)
}
