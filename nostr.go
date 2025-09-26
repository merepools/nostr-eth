package nostreth

// Re-export all functions from the log package
import (
	"math/big"

	"github.com/comunifi/nostr-eth/pkg/event"
	"github.com/comunifi/nostr-eth/pkg/neth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nbd-wtf/go-nostr"
)

// Re-export log package types
type TxLogEvent = event.TxLogEvent
type Log = neth.Log
type UserOpEvent = event.UserOpEvent
type UserOp = neth.UserOp

// Re-export group package types
type GroupMetadata = event.GroupMetadata
type GroupMessage = event.GroupMessage
type GroupJoin = event.GroupJoin
type GroupLeave = event.GroupLeave
type GroupModeration = event.GroupModeration

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

// Re-export group package constants
const (
	KindGroupMetadata   = event.KindGroupMetadata
	KindGroupMessage    = event.KindGroupMessage
	KindGroupJoin       = event.KindGroupJoin
	KindGroupLeave      = event.KindGroupLeave
	KindGroupModeration = event.KindGroupModeration
)

// Re-export log package functions
func CreateTxLogEvent(log neth.Log) (*nostr.Event, error) {
	return event.CreateTxLogEvent(log)
}

func UpdateTxLogEvent(log neth.Log, ev *nostr.Event) (*nostr.Event, error) {
	return event.UpdateTxLogEvent(log, ev)
}

func ParseTxLogEvent(evt *nostr.Event) (*event.TxLogEvent, error) {
	return event.ParseTxLogEvent(evt)
}

func GetEventData(log neth.Log) (map[string]interface{}, error) {
	return log.GetEventData()
}

func RequestUserOpEvent(chainID *big.Int, paymaster *common.Address, userOp neth.UserOp) (*nostr.Event, error) {
	return event.RequestUserOpEvent(chainID, paymaster, userOp)
}

func UpdateUserOpEvent(chainID *big.Int, userOp neth.UserOp, eventType event.EventTypeUserOp, ev *nostr.Event) (*nostr.Event, error) {
	return event.UpdateUserOpEvent(chainID, userOp, eventType, ev)
}

func ParseUserOpEvent(evt *nostr.Event) (*event.UserOpEvent, error) {
	return event.ParseUserOpEvent(evt)
}

// Re-export group package functions
func CreateGroupMetadataEvent(groupID, name, about, picture string, admins, moderators []string, private, closed bool) (*nostr.Event, error) {
	return event.CreateGroupMetadataEvent(groupID, name, about, picture, admins, moderators, private, closed)
}

func CreateGroupMessageEvent(groupID, content, replyTo string, mentions []string) (*nostr.Event, error) {
	return event.CreateGroupMessageEvent(groupID, content, replyTo, mentions)
}

func CreateGroupJoinEvent(groupID, user, role string) (*nostr.Event, error) {
	return event.CreateGroupJoinEvent(groupID, user, role)
}

func CreateGroupLeaveEvent(groupID, user, reason string) (*nostr.Event, error) {
	return event.CreateGroupLeaveEvent(groupID, user, reason)
}

func CreateGroupModerationEvent(groupID, action, target, reason string, duration int64) (*nostr.Event, error) {
	return event.CreateGroupModerationEvent(groupID, action, target, reason, duration)
}

func ParseGroupMetadataEvent(evt *nostr.Event) (*event.GroupMetadata, error) {
	return event.ParseGroupMetadataEvent(evt)
}

func ParseGroupMessageEvent(evt *nostr.Event) (*event.GroupMessage, error) {
	return event.ParseGroupMessageEvent(evt)
}

func ParseGroupJoinEvent(evt *nostr.Event) (*event.GroupJoin, error) {
	return event.ParseGroupJoinEvent(evt)
}

func ParseGroupLeaveEvent(evt *nostr.Event) (*event.GroupLeave, error) {
	return event.ParseGroupLeaveEvent(evt)
}

func ParseGroupModerationEvent(evt *nostr.Event) (*event.GroupModeration, error) {
	return event.ParseGroupModerationEvent(evt)
}

func GetGroupIDFromEvent(evt *nostr.Event) (string, error) {
	return event.GetGroupIDFromEvent(evt)
}

func IsGroupEvent(evt *nostr.Event) bool {
	return event.IsGroupEvent(evt)
}

func FilterGroupEventsByGroupID(events []*nostr.Event, groupID string) []*nostr.Event {
	return event.FilterGroupEventsByGroupID(events, groupID)
}

func ParseGroupIdentifier(groupIdentifier string) (host, groupID string, err error) {
	return event.ParseGroupIdentifier(groupIdentifier)
}

func FormatGroupIdentifier(host, groupID string) string {
	return event.FormatGroupIdentifier(host, groupID)
}

func ValidateGroupID(groupID string) error {
	return event.ValidateGroupID(groupID)
}

func GetEventTypeFromGroupEvent(evt *nostr.Event) string {
	return event.GetEventTypeFromGroupEvent(evt)
}
