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

// Group Metadata Event Types (39000s)
type GroupMetadataEvent = event.GroupMetadataEvent
type GroupNameEvent = event.GroupNameEvent
type GroupAboutEvent = event.GroupAboutEvent
type GroupPictureEvent = event.GroupPictureEvent
type GroupAdminsEvent = event.GroupAdminsEvent
type GroupModeratorsEvent = event.GroupModeratorsEvent
type GroupPrivateEvent = event.GroupPrivateEvent
type GroupClosedEvent = event.GroupClosedEvent
type GroupCreatedEvent = event.GroupCreatedEvent
type GroupUpdatedEvent = event.GroupUpdatedEvent

// Re-export log package constants
const (
	KindTxLog       = event.KindTxLog
	EventUserOpKind = event.EventUserOpKind

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
	// Group Moderation Events (9000s)
	KindGroupAddUser      = event.KindGroupAddUser
	KindGroupRemoveUser   = event.KindGroupRemoveUser
	KindGroupEditMetadata = event.KindGroupEditMetadata
	KindGroupAddAdmin     = event.KindGroupAddAdmin
	KindGroupRemoveAdmin  = event.KindGroupRemoveAdmin
	KindGroupDeleteEvent  = event.KindGroupDeleteEvent
	KindGroupUpdateStatus = event.KindGroupUpdateStatus
	KindGroupCreate       = event.KindGroupCreate
	KindGroupDelete       = event.KindGroupDelete
	KindGroupJoinRequest  = event.KindGroupJoinRequest

	// Group Metadata Events (39000s)
	KindGroupMetadata   = event.KindGroupMetadata
	KindGroupName       = event.KindGroupName
	KindGroupAbout      = event.KindGroupAbout
	KindGroupPicture    = event.KindGroupPicture
	KindGroupAdmins     = event.KindGroupAdmins
	KindGroupModerators = event.KindGroupModerators
	KindGroupPrivate    = event.KindGroupPrivate
	KindGroupClosed     = event.KindGroupClosed
	KindGroupCreated    = event.KindGroupCreated
	KindGroupUpdated    = event.KindGroupUpdated
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
// Group Moderation Events (9000s)
func CreateGroupEvent(groupID, name, about, picture string, admins, moderators []string, private, closed bool) (*nostr.Event, error) {
	return event.CreateGroupEvent(groupID, name, about, picture, admins, moderators, private, closed)
}

func CreateAddUserEvent(groupID, user, role string) (*nostr.Event, error) {
	return event.CreateAddUserEvent(groupID, user, role)
}

func CreateRemoveUserEvent(groupID, user, reason string) (*nostr.Event, error) {
	return event.CreateRemoveUserEvent(groupID, user, reason)
}

func CreateEditMetadataEvent(groupID, name, about, picture string, admins, moderators []string, private, closed bool) (*nostr.Event, error) {
	return event.CreateEditMetadataEvent(groupID, name, about, picture, admins, moderators, private, closed)
}

func CreateAddAdminEvent(groupID, user string) (*nostr.Event, error) {
	return event.CreateAddAdminEvent(groupID, user)
}

func CreateRemoveAdminEvent(groupID, user string) (*nostr.Event, error) {
	return event.CreateRemoveAdminEvent(groupID, user)
}

func CreateDeleteEventEvent(groupID, eventID string) (*nostr.Event, error) {
	return event.CreateDeleteEventEvent(groupID, eventID)
}

func CreateUpdateGroupStatusEvent(groupID, status string) (*nostr.Event, error) {
	return event.CreateUpdateGroupStatusEvent(groupID, status)
}

func CreateDeleteGroupEvent(groupID string) (*nostr.Event, error) {
	return event.CreateDeleteGroupEvent(groupID)
}

func CreateJoinRequestEvent(groupID, message string) (*nostr.Event, error) {
	return event.CreateJoinRequestEvent(groupID, message)
}

// Group Metadata Events (39000s)
func CreateGroupMetadataEvent(groupID string, metadata event.GroupMetadata) (*nostr.Event, error) {
	return event.CreateGroupMetadataEvent(groupID, metadata)
}

func CreateGroupNameEvent(groupID, name string) (*nostr.Event, error) {
	return event.CreateGroupNameEvent(groupID, name)
}

func CreateGroupAboutEvent(groupID, about string) (*nostr.Event, error) {
	return event.CreateGroupAboutEvent(groupID, about)
}

func CreateGroupPictureEvent(groupID, picture string) (*nostr.Event, error) {
	return event.CreateGroupPictureEvent(groupID, picture)
}

func CreateGroupAdminsEvent(groupID string, admins []string) (*nostr.Event, error) {
	return event.CreateGroupAdminsEvent(groupID, admins)
}

func CreateGroupModeratorsEvent(groupID string, moderators []string) (*nostr.Event, error) {
	return event.CreateGroupModeratorsEvent(groupID, moderators)
}

func CreateGroupPrivateEvent(groupID string, private bool) (*nostr.Event, error) {
	return event.CreateGroupPrivateEvent(groupID, private)
}

func CreateGroupClosedEvent(groupID string, closed bool) (*nostr.Event, error) {
	return event.CreateGroupClosedEvent(groupID, closed)
}

func CreateGroupCreatedEvent(groupID string, createdAt int64) (*nostr.Event, error) {
	return event.CreateGroupCreatedEvent(groupID, createdAt)
}

func CreateGroupUpdatedEvent(groupID string, updatedAt int64) (*nostr.Event, error) {
	return event.CreateGroupUpdatedEvent(groupID, updatedAt)
}

// Parse functions
func ParseGroupEvent(evt *nostr.Event) (*event.GroupMetadata, error) {
	return event.ParseGroupEvent(evt)
}

func ParseAddUserEvent(evt *nostr.Event) (*event.GroupJoin, error) {
	return event.ParseAddUserEvent(evt)
}

func ParseRemoveUserEvent(evt *nostr.Event) (*event.GroupLeave, error) {
	return event.ParseRemoveUserEvent(evt)
}

func ParseEditMetadataEvent(evt *nostr.Event) (*event.GroupMetadata, error) {
	return event.ParseEditMetadataEvent(evt)
}

func ParseGroupMetadataEvent(evt *nostr.Event) (*event.GroupMetadataEvent, error) {
	return event.ParseGroupMetadataEvent(evt)
}

func ParseGroupNameEvent(evt *nostr.Event) (*event.GroupNameEvent, error) {
	return event.ParseGroupNameEvent(evt)
}

func ParseGroupAboutEvent(evt *nostr.Event) (*event.GroupAboutEvent, error) {
	return event.ParseGroupAboutEvent(evt)
}

func ParseGroupPictureEvent(evt *nostr.Event) (*event.GroupPictureEvent, error) {
	return event.ParseGroupPictureEvent(evt)
}

func ParseGroupAdminsEvent(evt *nostr.Event) (*event.GroupAdminsEvent, error) {
	return event.ParseGroupAdminsEvent(evt)
}

func ParseGroupModeratorsEvent(evt *nostr.Event) (*event.GroupModeratorsEvent, error) {
	return event.ParseGroupModeratorsEvent(evt)
}

func ParseGroupPrivateEvent(evt *nostr.Event) (*event.GroupPrivateEvent, error) {
	return event.ParseGroupPrivateEvent(evt)
}

func ParseGroupClosedEvent(evt *nostr.Event) (*event.GroupClosedEvent, error) {
	return event.ParseGroupClosedEvent(evt)
}

func ParseGroupCreatedEvent(evt *nostr.Event) (*event.GroupCreatedEvent, error) {
	return event.ParseGroupCreatedEvent(evt)
}

func ParseGroupUpdatedEvent(evt *nostr.Event) (*event.GroupUpdatedEvent, error) {
	return event.ParseGroupUpdatedEvent(evt)
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
