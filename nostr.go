package eth

// Re-export all functions from the log package
import log "github.com/citizenwallet/nostr-eth/pkg/eth"

// Re-export log package types
type TxLogEvent = log.TxLogEvent
type NostrEvent = log.NostrEvent
type DataOutputter = log.DataOutputter
type MapDataOutputter = log.MapDataOutputter

// Re-export log package constants
const (
	EventTypeTxLogCreated = log.EventTypeTxLogCreated
	EventTypeTxLogUpdated = log.EventTypeTxLogUpdated
)

// Re-export log package functions
func CreateTxLogEvent(logData log.DataOutputter, privateKey string) (*log.NostrEvent, error) {
	return log.CreateTxLogEvent(logData, privateKey)
}

func UpdateTxLogEvent(logData map[string]interface{}, privateKey string, originalEventID ...string) (*log.NostrEvent, error) {
	return log.UpdateTxLogEvent(logData, privateKey, originalEventID...)
}

func ParseTxLogEvent(evt *log.NostrEvent) (*log.TxLogEvent, error) {
	return log.ParseTxLogEvent(evt)
}

func UpdateLogStatusEvent(logData map[string]interface{}, newStatus string, privateKey string, originalEventID ...string) (*log.NostrEvent, error) {
	return log.UpdateLogStatusEvent(logData, newStatus, privateKey, originalEventID...)
}

func GetTransferData(logData map[string]interface{}) (map[string]interface{}, error) {
	return log.GetTransferData(logData)
}

func NewMapDataOutputter(data map[string]interface{}) *log.MapDataOutputter {
	return log.NewMapDataOutputter(data)
}
