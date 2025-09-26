package event

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/nbd-wtf/go-nostr"
)

// NIP-29 Group Event Kinds
const (
	// Group Moderation Events (9000s)
	KindGroupAddUser      = 9000 // Add User
	KindGroupRemoveUser   = 9001 // Remove User
	KindGroupEditMetadata = 9002 // Edit Metadata
	KindGroupAddAdmin     = 9003 // Add Admin Permission
	KindGroupRemoveAdmin  = 9004 // Remove Admin Permission
	KindGroupDeleteEvent  = 9005 // Delete Event
	KindGroupUpdateStatus = 9006 // Update Group Status
	KindGroupCreate       = 9007 // Create Group
	KindGroupDelete       = 9008 // Delete Group
	KindGroupJoinRequest  = 9021 // Join Request

	// Group Metadata Events (39000s)
	KindGroupMetadata   = 39000 // Group metadata
	KindGroupName       = 39001 // Group name
	KindGroupAbout      = 39002 // Group about/description
	KindGroupPicture    = 39003 // Group picture
	KindGroupAdmins     = 39004 // Group admins
	KindGroupModerators = 39005 // Group moderators
	KindGroupPrivate    = 39006 // Group privacy setting
	KindGroupClosed     = 39007 // Group closed setting
	KindGroupCreated    = 39008 // Group creation timestamp
	KindGroupUpdated    = 39009 // Group update timestamp
)

// GroupMetadata represents the metadata for a group
type GroupMetadata struct {
	Name       string   `json:"name"`
	About      string   `json:"about,omitempty"`
	Picture    string   `json:"picture,omitempty"`
	Admins     []string `json:"admins,omitempty"`
	Moderators []string `json:"moderators,omitempty"`
	Private    bool     `json:"private,omitempty"`
	Closed     bool     `json:"closed,omitempty"`
	CreatedAt  int64    `json:"created_at"`
	UpdatedAt  int64    `json:"updated_at"`
}

// GroupMessage represents a message sent to a group
type GroupMessage struct {
	Content   string   `json:"content"`
	ReplyTo   string   `json:"reply_to,omitempty"`
	Mentions  []string `json:"mentions,omitempty"`
	CreatedAt int64    `json:"created_at"`
}

// GroupJoin represents a user joining a group
type GroupJoin struct {
	User     string `json:"user"`
	JoinedAt int64  `json:"joined_at"`
	Role     string `json:"role,omitempty"` // "admin", "moderator", "member"
}

// GroupLeave represents a user leaving a group
type GroupLeave struct {
	User   string `json:"user"`
	LeftAt int64  `json:"left_at"`
	Reason string `json:"reason,omitempty"`
}

// GroupModeration represents moderation actions in a group
type GroupModeration struct {
	Action    string `json:"action"` // "ban", "unban", "mute", "unmute", "promote", "demote"
	Target    string `json:"target"` // User being moderated
	Reason    string `json:"reason,omitempty"`
	Duration  int64  `json:"duration,omitempty"` // Duration in seconds for temporary actions
	CreatedAt int64  `json:"created_at"`
}

// GroupMetadataEvent represents a group metadata event (kind 39000)
type GroupMetadataEvent struct {
	GroupID   string        `json:"group_id"`
	Metadata  GroupMetadata `json:"metadata"`
	CreatedAt int64         `json:"created_at"`
}

// GroupNameEvent represents a group name event (kind 39001)
type GroupNameEvent struct {
	GroupID   string `json:"group_id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
}

// GroupAboutEvent represents a group about/description event (kind 39002)
type GroupAboutEvent struct {
	GroupID   string `json:"group_id"`
	About     string `json:"about"`
	CreatedAt int64  `json:"created_at"`
}

// GroupPictureEvent represents a group picture event (kind 39003)
type GroupPictureEvent struct {
	GroupID   string `json:"group_id"`
	Picture   string `json:"picture"`
	CreatedAt int64  `json:"created_at"`
}

// GroupAdminsEvent represents a group admins event (kind 39004)
type GroupAdminsEvent struct {
	GroupID   string   `json:"group_id"`
	Admins    []string `json:"admins"`
	CreatedAt int64    `json:"created_at"`
}

// GroupModeratorsEvent represents a group moderators event (kind 39005)
type GroupModeratorsEvent struct {
	GroupID    string   `json:"group_id"`
	Moderators []string `json:"moderators"`
	CreatedAt  int64    `json:"created_at"`
}

// GroupPrivateEvent represents a group privacy setting event (kind 39006)
type GroupPrivateEvent struct {
	GroupID   string `json:"group_id"`
	Private   bool   `json:"private"`
	CreatedAt int64  `json:"created_at"`
}

// GroupClosedEvent represents a group closed setting event (kind 39007)
type GroupClosedEvent struct {
	GroupID   string `json:"group_id"`
	Closed    bool   `json:"closed"`
	CreatedAt int64  `json:"created_at"`
}

// GroupCreatedEvent represents a group creation timestamp event (kind 39008)
type GroupCreatedEvent struct {
	GroupID   string `json:"group_id"`
	CreatedAt int64  `json:"created_at"`
}

// GroupUpdatedEvent represents a group update timestamp event (kind 39009)
type GroupUpdatedEvent struct {
	GroupID   string `json:"group_id"`
	UpdatedAt int64  `json:"updated_at"`
	CreatedAt int64  `json:"created_at"`
}

// CreateGroupEvent creates a group event (kind 9007)
func CreateGroupEvent(groupID, name, about, picture string, admins, moderators []string, private, closed bool) (*nostr.Event, error) {
	now := time.Now().Unix()

	metadata := GroupMetadata{
		Name:       name,
		About:      about,
		Picture:    picture,
		Admins:     admins,
		Moderators: moderators,
		Private:    private,
		Closed:     closed,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	content, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal group metadata: %w", err)
	}

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupCreate,
		Tags:      make([]nostr.Tag, 0),
		Content:   string(content),
	}

	// Add group identifier tag (h tag with group ID)
	evt.Tags = append(evt.Tags, []string{"h", groupID})

	// Add admin tags
	for _, admin := range admins {
		evt.Tags = append(evt.Tags, []string{"p", admin, "admin"})
	}

	// Add moderator tags
	for _, moderator := range moderators {
		evt.Tags = append(evt.Tags, []string{"p", moderator, "moderator"})
	}

	// Add group type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "metadata"})

	if private {
		evt.Tags = append(evt.Tags, []string{"t", "private"})
	}

	if closed {
		evt.Tags = append(evt.Tags, []string{"t", "closed"})
	}

	return evt, nil
}

// CreateAddUserEvent creates an add user event (kind 9000)
func CreateAddUserEvent(groupID, user, role string) (*nostr.Event, error) {
	now := time.Now().Unix()

	join := GroupJoin{
		User:     user,
		JoinedAt: now,
		Role:     role,
	}

	content, err := json.Marshal(join)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal add user: %w", err)
	}

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupAddUser,
		Tags:      make([]nostr.Tag, 0),
		Content:   string(content),
	}

	// Add group identifier tag (h tag with group ID)
	evt.Tags = append(evt.Tags, []string{"h", groupID})

	// Add user tag
	evt.Tags = append(evt.Tags, []string{"p", user, "member"})

	// Add role tag if specified
	if role != "" {
		evt.Tags = append(evt.Tags, []string{"t", role})
	}

	// Add add user type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "add_user"})

	return evt, nil
}

// CreateRemoveUserEvent creates a remove user event (kind 9001)
func CreateRemoveUserEvent(groupID, user, reason string) (*nostr.Event, error) {
	now := time.Now().Unix()

	leave := GroupLeave{
		User:   user,
		LeftAt: now,
		Reason: reason,
	}

	content, err := json.Marshal(leave)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal remove user: %w", err)
	}

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupRemoveUser,
		Tags:      make([]nostr.Tag, 0),
		Content:   string(content),
	}

	// Add group identifier tag (h tag with group ID)
	evt.Tags = append(evt.Tags, []string{"h", groupID})

	// Add user tag
	evt.Tags = append(evt.Tags, []string{"p", user, "former_member"})

	// Add remove user type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "remove_user"})

	return evt, nil
}

// CreateEditMetadataEvent creates an edit metadata event (kind 9002)
func CreateEditMetadataEvent(groupID, name, about, picture string, admins, moderators []string, private, closed bool) (*nostr.Event, error) {
	now := time.Now().Unix()

	metadata := GroupMetadata{
		Name:       name,
		About:      about,
		Picture:    picture,
		Admins:     admins,
		Moderators: moderators,
		Private:    private,
		Closed:     closed,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	content, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal group metadata: %w", err)
	}

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupEditMetadata,
		Tags:      make([]nostr.Tag, 0),
		Content:   string(content),
	}

	// Add group identifier tag (h tag with group ID)
	evt.Tags = append(evt.Tags, []string{"h", groupID})

	// Add admin tags
	for _, admin := range admins {
		evt.Tags = append(evt.Tags, []string{"p", admin, "admin"})
	}

	// Add moderator tags
	for _, moderator := range moderators {
		evt.Tags = append(evt.Tags, []string{"p", moderator, "moderator"})
	}

	// Add group type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "edit_metadata"})

	if private {
		evt.Tags = append(evt.Tags, []string{"t", "private"})
	}

	if closed {
		evt.Tags = append(evt.Tags, []string{"t", "closed"})
	}

	return evt, nil
}

// CreateAddAdminEvent creates an add admin event (kind 9003)
func CreateAddAdminEvent(groupID, user string) (*nostr.Event, error) {
	now := time.Now().Unix()

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupAddAdmin,
		Tags:      make([]nostr.Tag, 0),
		Content:   "",
	}

	// Add group identifier tag (h tag with group ID)
	evt.Tags = append(evt.Tags, []string{"h", groupID})

	// Add user tag
	evt.Tags = append(evt.Tags, []string{"p", user, "admin"})

	// Add add admin type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "add_admin"})

	return evt, nil
}

// CreateRemoveAdminEvent creates a remove admin event (kind 9004)
func CreateRemoveAdminEvent(groupID, user string) (*nostr.Event, error) {
	now := time.Now().Unix()

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupRemoveAdmin,
		Tags:      make([]nostr.Tag, 0),
		Content:   "",
	}

	// Add group identifier tag (h tag with group ID)
	evt.Tags = append(evt.Tags, []string{"h", groupID})

	// Add user tag
	evt.Tags = append(evt.Tags, []string{"p", user, "former_admin"})

	// Add remove admin type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "remove_admin"})

	return evt, nil
}

// CreateDeleteEventEvent creates a delete event event (kind 9005)
func CreateDeleteEventEvent(groupID, eventID string) (*nostr.Event, error) {
	now := time.Now().Unix()

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupDeleteEvent,
		Tags:      make([]nostr.Tag, 0),
		Content:   "",
	}

	// Add group identifier tag (h tag with group ID)
	evt.Tags = append(evt.Tags, []string{"h", groupID})

	// Add event tag
	evt.Tags = append(evt.Tags, []string{"e", eventID, "delete"})

	// Add delete event type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "delete_event"})

	return evt, nil
}

// CreateUpdateGroupStatusEvent creates an update group status event (kind 9006)
func CreateUpdateGroupStatusEvent(groupID, status string) (*nostr.Event, error) {
	now := time.Now().Unix()

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupUpdateStatus,
		Tags:      make([]nostr.Tag, 0),
		Content:   status,
	}

	// Add group identifier tag (h tag with group ID)
	evt.Tags = append(evt.Tags, []string{"h", groupID})

	// Add status tag
	evt.Tags = append(evt.Tags, []string{"t", status})

	// Add update status type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "update_status"})

	return evt, nil
}

// CreateDeleteGroupEvent creates a delete group event (kind 9008)
func CreateDeleteGroupEvent(groupID string) (*nostr.Event, error) {
	now := time.Now().Unix()

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupDelete,
		Tags:      make([]nostr.Tag, 0),
		Content:   "",
	}

	// Add group identifier tag (h tag with group ID)
	evt.Tags = append(evt.Tags, []string{"h", groupID})

	// Add delete group type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "delete_group"})

	return evt, nil
}

// CreateJoinRequestEvent creates a join request event (kind 9021)
func CreateJoinRequestEvent(groupID, message string) (*nostr.Event, error) {
	now := time.Now().Unix()

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupJoinRequest,
		Tags:      make([]nostr.Tag, 0),
		Content:   message,
	}

	// Add group identifier tag (h tag with group ID)
	evt.Tags = append(evt.Tags, []string{"h", groupID})

	// Add join request type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "join_request"})

	return evt, nil
}

// CreateGroupMetadataEvent creates a group metadata event (kind 39000)
func CreateGroupMetadataEvent(groupID string, metadata GroupMetadata) (*nostr.Event, error) {
	now := time.Now().Unix()

	eventData := GroupMetadataEvent{
		GroupID:   groupID,
		Metadata:  metadata,
		CreatedAt: now,
	}

	content, err := json.Marshal(eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal group metadata event: %w", err)
	}

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupMetadata,
		Tags:      make([]nostr.Tag, 0),
		Content:   string(content),
	}

	// Add group identifier tag (h tag with group ID)
	evt.Tags = append(evt.Tags, []string{"h", groupID})

	// Add group metadata type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "metadata"})

	return evt, nil
}

// CreateGroupNameEvent creates a group name event (kind 39001)
func CreateGroupNameEvent(groupID, name string) (*nostr.Event, error) {
	now := time.Now().Unix()

	eventData := GroupNameEvent{
		GroupID:   groupID,
		Name:      name,
		CreatedAt: now,
	}

	content, err := json.Marshal(eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal group name event: %w", err)
	}

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupName,
		Tags:      make([]nostr.Tag, 0),
		Content:   string(content),
	}

	// Add group identifier tag (h tag with group ID)
	evt.Tags = append(evt.Tags, []string{"h", groupID})

	// Add group name type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "name"})

	return evt, nil
}

// CreateGroupAboutEvent creates a group about event (kind 39002)
func CreateGroupAboutEvent(groupID, about string) (*nostr.Event, error) {
	now := time.Now().Unix()

	eventData := GroupAboutEvent{
		GroupID:   groupID,
		About:     about,
		CreatedAt: now,
	}

	content, err := json.Marshal(eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal group about event: %w", err)
	}

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupAbout,
		Tags:      make([]nostr.Tag, 0),
		Content:   string(content),
	}

	// Add group identifier tag (h tag with group ID)
	evt.Tags = append(evt.Tags, []string{"h", groupID})

	// Add group about type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "about"})

	return evt, nil
}

// CreateGroupPictureEvent creates a group picture event (kind 39003)
func CreateGroupPictureEvent(groupID, picture string) (*nostr.Event, error) {
	now := time.Now().Unix()

	eventData := GroupPictureEvent{
		GroupID:   groupID,
		Picture:   picture,
		CreatedAt: now,
	}

	content, err := json.Marshal(eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal group picture event: %w", err)
	}

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupPicture,
		Tags:      make([]nostr.Tag, 0),
		Content:   string(content),
	}

	// Add group identifier tag (h tag with group ID)
	evt.Tags = append(evt.Tags, []string{"h", groupID})

	// Add group picture type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "picture"})

	return evt, nil
}

// CreateGroupAdminsEvent creates a group admins event (kind 39004)
func CreateGroupAdminsEvent(groupID string, admins []string) (*nostr.Event, error) {
	now := time.Now().Unix()

	eventData := GroupAdminsEvent{
		GroupID:   groupID,
		Admins:    admins,
		CreatedAt: now,
	}

	content, err := json.Marshal(eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal group admins event: %w", err)
	}

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupAdmins,
		Tags:      make([]nostr.Tag, 0),
		Content:   string(content),
	}

	// Add group identifier tag (h tag with group ID)
	evt.Tags = append(evt.Tags, []string{"h", groupID})

	// Add admin tags
	for _, admin := range admins {
		evt.Tags = append(evt.Tags, []string{"p", admin, "admin"})
	}

	// Add group admins type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "admins"})

	return evt, nil
}

// CreateGroupModeratorsEvent creates a group moderators event (kind 39005)
func CreateGroupModeratorsEvent(groupID string, moderators []string) (*nostr.Event, error) {
	now := time.Now().Unix()

	eventData := GroupModeratorsEvent{
		GroupID:    groupID,
		Moderators: moderators,
		CreatedAt:  now,
	}

	content, err := json.Marshal(eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal group moderators event: %w", err)
	}

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupModerators,
		Tags:      make([]nostr.Tag, 0),
		Content:   string(content),
	}

	// Add group identifier tag (h tag with group ID)
	evt.Tags = append(evt.Tags, []string{"h", groupID})

	// Add moderator tags
	for _, moderator := range moderators {
		evt.Tags = append(evt.Tags, []string{"p", moderator, "moderator"})
	}

	// Add group moderators type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "moderators"})

	return evt, nil
}

// CreateGroupPrivateEvent creates a group private event (kind 39006)
func CreateGroupPrivateEvent(groupID string, private bool) (*nostr.Event, error) {
	now := time.Now().Unix()

	eventData := GroupPrivateEvent{
		GroupID:   groupID,
		Private:   private,
		CreatedAt: now,
	}

	content, err := json.Marshal(eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal group private event: %w", err)
	}

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupPrivate,
		Tags:      make([]nostr.Tag, 0),
		Content:   string(content),
	}

	// Add group identifier tag (h tag with group ID)
	evt.Tags = append(evt.Tags, []string{"h", groupID})

	// Add group private type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "private"})

	if private {
		evt.Tags = append(evt.Tags, []string{"t", "private"})
	}

	return evt, nil
}

// CreateGroupClosedEvent creates a group closed event (kind 39007)
func CreateGroupClosedEvent(groupID string, closed bool) (*nostr.Event, error) {
	now := time.Now().Unix()

	eventData := GroupClosedEvent{
		GroupID:   groupID,
		Closed:    closed,
		CreatedAt: now,
	}

	content, err := json.Marshal(eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal group closed event: %w", err)
	}

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupClosed,
		Tags:      make([]nostr.Tag, 0),
		Content:   string(content),
	}

	// Add group identifier tag (h tag with group ID)
	evt.Tags = append(evt.Tags, []string{"h", groupID})

	// Add group closed type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "closed"})

	if closed {
		evt.Tags = append(evt.Tags, []string{"t", "closed"})
	}

	return evt, nil
}

// CreateGroupCreatedEvent creates a group created event (kind 39008)
func CreateGroupCreatedEvent(groupID string, createdAt int64) (*nostr.Event, error) {
	now := time.Now().Unix()

	eventData := GroupCreatedEvent{
		GroupID:   groupID,
		CreatedAt: createdAt,
	}

	content, err := json.Marshal(eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal group created event: %w", err)
	}

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupCreated,
		Tags:      make([]nostr.Tag, 0),
		Content:   string(content),
	}

	// Add group identifier tag (h tag with group ID)
	evt.Tags = append(evt.Tags, []string{"h", groupID})

	// Add group created type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "created"})

	return evt, nil
}

// CreateGroupUpdatedEvent creates a group updated event (kind 39009)
func CreateGroupUpdatedEvent(groupID string, updatedAt int64) (*nostr.Event, error) {
	now := time.Now().Unix()

	eventData := GroupUpdatedEvent{
		GroupID:   groupID,
		UpdatedAt: updatedAt,
		CreatedAt: now,
	}

	content, err := json.Marshal(eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal group updated event: %w", err)
	}

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupUpdated,
		Tags:      make([]nostr.Tag, 0),
		Content:   string(content),
	}

	// Add group identifier tag (h tag with group ID)
	evt.Tags = append(evt.Tags, []string{"h", groupID})

	// Add group updated type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "updated"})

	return evt, nil
}

// ParseGroupEvent parses a group creation event (kind 9007)
func ParseGroupEvent(evt *nostr.Event) (*GroupMetadata, error) {
	if evt.Kind != KindGroupCreate {
		return nil, fmt.Errorf("event is not a group creation event (kind %d)", evt.Kind)
	}

	var metadata GroupMetadata
	err := json.Unmarshal([]byte(evt.Content), &metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal group metadata: %w", err)
	}

	return &metadata, nil
}

// ParseEditMetadataEvent parses an edit metadata event (kind 9002)
func ParseEditMetadataEvent(evt *nostr.Event) (*GroupMetadata, error) {
	if evt.Kind != KindGroupEditMetadata {
		return nil, fmt.Errorf("event is not an edit metadata event (kind %d)", evt.Kind)
	}

	var metadata GroupMetadata
	err := json.Unmarshal([]byte(evt.Content), &metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal group metadata: %w", err)
	}

	return &metadata, nil
}

// ParseAddUserEvent parses an add user event (kind 9000)
func ParseAddUserEvent(evt *nostr.Event) (*GroupJoin, error) {
	if evt.Kind != KindGroupAddUser {
		return nil, fmt.Errorf("event is not an add user event (kind %d)", evt.Kind)
	}

	var join GroupJoin
	err := json.Unmarshal([]byte(evt.Content), &join)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal add user: %w", err)
	}

	return &join, nil
}

// ParseRemoveUserEvent parses a remove user event (kind 9001)
func ParseRemoveUserEvent(evt *nostr.Event) (*GroupLeave, error) {
	if evt.Kind != KindGroupRemoveUser {
		return nil, fmt.Errorf("event is not a remove user event (kind %d)", evt.Kind)
	}

	var leave GroupLeave
	err := json.Unmarshal([]byte(evt.Content), &leave)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal remove user: %w", err)
	}

	return &leave, nil
}

// ParseGroupMetadataEvent parses a group metadata event (kind 39000)
func ParseGroupMetadataEvent(evt *nostr.Event) (*GroupMetadataEvent, error) {
	if evt.Kind != KindGroupMetadata {
		return nil, fmt.Errorf("event is not a group metadata event (kind %d)", evt.Kind)
	}

	var eventData GroupMetadataEvent
	err := json.Unmarshal([]byte(evt.Content), &eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal group metadata event: %w", err)
	}

	return &eventData, nil
}

// ParseGroupNameEvent parses a group name event (kind 39001)
func ParseGroupNameEvent(evt *nostr.Event) (*GroupNameEvent, error) {
	if evt.Kind != KindGroupName {
		return nil, fmt.Errorf("event is not a group name event (kind %d)", evt.Kind)
	}

	var eventData GroupNameEvent
	err := json.Unmarshal([]byte(evt.Content), &eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal group name event: %w", err)
	}

	return &eventData, nil
}

// ParseGroupAboutEvent parses a group about event (kind 39002)
func ParseGroupAboutEvent(evt *nostr.Event) (*GroupAboutEvent, error) {
	if evt.Kind != KindGroupAbout {
		return nil, fmt.Errorf("event is not a group about event (kind %d)", evt.Kind)
	}

	var eventData GroupAboutEvent
	err := json.Unmarshal([]byte(evt.Content), &eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal group about event: %w", err)
	}

	return &eventData, nil
}

// ParseGroupPictureEvent parses a group picture event (kind 39003)
func ParseGroupPictureEvent(evt *nostr.Event) (*GroupPictureEvent, error) {
	if evt.Kind != KindGroupPicture {
		return nil, fmt.Errorf("event is not a group picture event (kind %d)", evt.Kind)
	}

	var eventData GroupPictureEvent
	err := json.Unmarshal([]byte(evt.Content), &eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal group picture event: %w", err)
	}

	return &eventData, nil
}

// ParseGroupAdminsEvent parses a group admins event (kind 39004)
func ParseGroupAdminsEvent(evt *nostr.Event) (*GroupAdminsEvent, error) {
	if evt.Kind != KindGroupAdmins {
		return nil, fmt.Errorf("event is not a group admins event (kind %d)", evt.Kind)
	}

	var eventData GroupAdminsEvent
	err := json.Unmarshal([]byte(evt.Content), &eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal group admins event: %w", err)
	}

	return &eventData, nil
}

// ParseGroupModeratorsEvent parses a group moderators event (kind 39005)
func ParseGroupModeratorsEvent(evt *nostr.Event) (*GroupModeratorsEvent, error) {
	if evt.Kind != KindGroupModerators {
		return nil, fmt.Errorf("event is not a group moderators event (kind %d)", evt.Kind)
	}

	var eventData GroupModeratorsEvent
	err := json.Unmarshal([]byte(evt.Content), &eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal group moderators event: %w", err)
	}

	return &eventData, nil
}

// ParseGroupPrivateEvent parses a group private event (kind 39006)
func ParseGroupPrivateEvent(evt *nostr.Event) (*GroupPrivateEvent, error) {
	if evt.Kind != KindGroupPrivate {
		return nil, fmt.Errorf("event is not a group private event (kind %d)", evt.Kind)
	}

	var eventData GroupPrivateEvent
	err := json.Unmarshal([]byte(evt.Content), &eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal group private event: %w", err)
	}

	return &eventData, nil
}

// ParseGroupClosedEvent parses a group closed event (kind 39007)
func ParseGroupClosedEvent(evt *nostr.Event) (*GroupClosedEvent, error) {
	if evt.Kind != KindGroupClosed {
		return nil, fmt.Errorf("event is not a group closed event (kind %d)", evt.Kind)
	}

	var eventData GroupClosedEvent
	err := json.Unmarshal([]byte(evt.Content), &eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal group closed event: %w", err)
	}

	return &eventData, nil
}

// ParseGroupCreatedEvent parses a group created event (kind 39008)
func ParseGroupCreatedEvent(evt *nostr.Event) (*GroupCreatedEvent, error) {
	if evt.Kind != KindGroupCreated {
		return nil, fmt.Errorf("event is not a group created event (kind %d)", evt.Kind)
	}

	var eventData GroupCreatedEvent
	err := json.Unmarshal([]byte(evt.Content), &eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal group created event: %w", err)
	}

	return &eventData, nil
}

// ParseGroupUpdatedEvent parses a group updated event (kind 39009)
func ParseGroupUpdatedEvent(evt *nostr.Event) (*GroupUpdatedEvent, error) {
	if evt.Kind != KindGroupUpdated {
		return nil, fmt.Errorf("event is not a group updated event (kind %d)", evt.Kind)
	}

	var eventData GroupUpdatedEvent
	err := json.Unmarshal([]byte(evt.Content), &eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal group updated event: %w", err)
	}

	return &eventData, nil
}

// GetGroupIDFromEvent extracts the group ID from a Nostr event
func GetGroupIDFromEvent(evt *nostr.Event) (string, error) {
	for _, tag := range evt.Tags {
		if len(tag) >= 2 && tag[0] == "h" {
			return tag[1], nil
		}
	}
	return "", fmt.Errorf("group ID tag (h) not found in event")
}

// IsGroupEvent checks if a Nostr event is a group-related event
func IsGroupEvent(evt *nostr.Event) bool {
	// Group Moderation Events (9000s)
	if evt.Kind == KindGroupAddUser ||
		evt.Kind == KindGroupRemoveUser ||
		evt.Kind == KindGroupEditMetadata ||
		evt.Kind == KindGroupAddAdmin ||
		evt.Kind == KindGroupRemoveAdmin ||
		evt.Kind == KindGroupDeleteEvent ||
		evt.Kind == KindGroupUpdateStatus ||
		evt.Kind == KindGroupCreate ||
		evt.Kind == KindGroupDelete ||
		evt.Kind == KindGroupJoinRequest {
		return true
	}

	// Group Metadata Events (39000s)
	if evt.Kind == KindGroupMetadata ||
		evt.Kind == KindGroupName ||
		evt.Kind == KindGroupAbout ||
		evt.Kind == KindGroupPicture ||
		evt.Kind == KindGroupAdmins ||
		evt.Kind == KindGroupModerators ||
		evt.Kind == KindGroupPrivate ||
		evt.Kind == KindGroupClosed ||
		evt.Kind == KindGroupCreated ||
		evt.Kind == KindGroupUpdated {
		return true
	}

	return false
}

// FilterGroupEventsByGroupID filters a list of events by group ID
func FilterGroupEventsByGroupID(events []*nostr.Event, groupID string) []*nostr.Event {
	var filtered []*nostr.Event
	for _, evt := range events {
		if groupTag, err := GetGroupIDFromEvent(evt); err == nil && groupTag == groupID {
			filtered = append(filtered, evt)
		}
	}
	return filtered
}

// ParseGroupIdentifier parses a group identifier in the format "host'group-id"
func ParseGroupIdentifier(groupIdentifier string) (host, groupID string, err error) {
	parts := strings.Split(groupIdentifier, "'")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid group identifier format: expected 'host'group-id', got '%s'", groupIdentifier)
	}
	return parts[0], parts[1], nil
}

// FormatGroupIdentifier formats a host and group ID into a group identifier
func FormatGroupIdentifier(host, groupID string) string {
	return fmt.Sprintf("%s'%s", host, groupID)
}

// ValidateGroupID validates that a group ID is properly formatted
func ValidateGroupID(groupID string) error {
	if groupID == "" {
		return fmt.Errorf("group ID cannot be empty")
	}

	// Group IDs should be alphanumeric and can contain hyphens and underscores
	for _, char := range groupID {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_') {
			return fmt.Errorf("group ID contains invalid character: %c", char)
		}
	}

	return nil
}

// GetEventTypeFromGroupEvent determines the type of group event
func GetEventTypeFromGroupEvent(evt *nostr.Event) string {
	switch evt.Kind {
	// Group Moderation Events (9000s)
	case KindGroupAddUser:
		return "add_user"
	case KindGroupRemoveUser:
		return "remove_user"
	case KindGroupEditMetadata:
		return "edit_metadata"
	case KindGroupAddAdmin:
		return "add_admin"
	case KindGroupRemoveAdmin:
		return "remove_admin"
	case KindGroupDeleteEvent:
		return "delete_event"
	case KindGroupUpdateStatus:
		return "update_group_status"
	case KindGroupCreate:
		return "create_group"
	case KindGroupDelete:
		return "delete_group"
	case KindGroupJoinRequest:
		return "join_request"

	// Group Metadata Events (39000s)
	case KindGroupMetadata:
		return "group_metadata"
	case KindGroupName:
		return "group_name"
	case KindGroupAbout:
		return "group_about"
	case KindGroupPicture:
		return "group_picture"
	case KindGroupAdmins:
		return "group_admins"
	case KindGroupModerators:
		return "group_moderators"
	case KindGroupPrivate:
		return "group_private"
	case KindGroupClosed:
		return "group_closed"
	case KindGroupCreated:
		return "group_created"
	case KindGroupUpdated:
		return "group_updated"

	default:
		return "unknown"
	}
}
