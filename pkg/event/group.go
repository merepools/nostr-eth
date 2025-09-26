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
	KindGroupMetadata   = 39000 // Group metadata (creation, updates)
	KindGroupMessage    = 39001 // Group messages
	KindGroupJoin       = 39002 // User joins group
	KindGroupLeave      = 39003 // User leaves group
	KindGroupModeration = 39004 // Group moderation actions
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

// CreateGroupMetadataEvent creates a group metadata event (kind 39000)
func CreateGroupMetadataEvent(groupID, name, about, picture string, admins, moderators []string, private, closed bool) (*nostr.Event, error) {
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
		Kind:      KindGroupMetadata,
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

// CreateGroupMessageEvent creates a group message event (kind 39001)
func CreateGroupMessageEvent(groupID, content, replyTo string, mentions []string) (*nostr.Event, error) {
	now := time.Now().Unix()

	message := GroupMessage{
		Content:   content,
		ReplyTo:   replyTo,
		Mentions:  mentions,
		CreatedAt: now,
	}

	messageContent, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal group message: %w", err)
	}

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupMessage,
		Tags:      make([]nostr.Tag, 0),
		Content:   string(messageContent),
	}

	// Add group identifier tag (h tag with group ID)
	evt.Tags = append(evt.Tags, []string{"h", groupID})

	// Add reply reference if this is a reply
	if replyTo != "" {
		evt.Tags = append(evt.Tags, []string{"e", replyTo, "reply"})
	}

	// Add mention tags
	for _, mention := range mentions {
		evt.Tags = append(evt.Tags, []string{"p", mention, "mention"})
	}

	// Add message type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "message"})

	return evt, nil
}

// CreateGroupJoinEvent creates a group join event (kind 39002)
func CreateGroupJoinEvent(groupID, user, role string) (*nostr.Event, error) {
	now := time.Now().Unix()

	join := GroupJoin{
		User:     user,
		JoinedAt: now,
		Role:     role,
	}

	content, err := json.Marshal(join)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal group join: %w", err)
	}

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupJoin,
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

	// Add join type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "join"})

	return evt, nil
}

// CreateGroupLeaveEvent creates a group leave event (kind 39003)
func CreateGroupLeaveEvent(groupID, user, reason string) (*nostr.Event, error) {
	now := time.Now().Unix()

	leave := GroupLeave{
		User:   user,
		LeftAt: now,
		Reason: reason,
	}

	content, err := json.Marshal(leave)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal group leave: %w", err)
	}

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupLeave,
		Tags:      make([]nostr.Tag, 0),
		Content:   string(content),
	}

	// Add group identifier tag (h tag with group ID)
	evt.Tags = append(evt.Tags, []string{"h", groupID})

	// Add user tag
	evt.Tags = append(evt.Tags, []string{"p", user, "former_member"})

	// Add leave type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "leave"})

	return evt, nil
}

// CreateGroupModerationEvent creates a group moderation event (kind 39004)
func CreateGroupModerationEvent(groupID, action, target, reason string, duration int64) (*nostr.Event, error) {
	now := time.Now().Unix()

	moderation := GroupModeration{
		Action:    action,
		Target:    target,
		Reason:    reason,
		Duration:  duration,
		CreatedAt: now,
	}

	content, err := json.Marshal(moderation)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal group moderation: %w", err)
	}

	evt := &nostr.Event{
		PubKey:    "", // Will be set by the client
		CreatedAt: nostr.Timestamp(now),
		Kind:      KindGroupModeration,
		Tags:      make([]nostr.Tag, 0),
		Content:   string(content),
	}

	// Add group identifier tag (h tag with group ID)
	evt.Tags = append(evt.Tags, []string{"h", groupID})

	// Add target user tag
	evt.Tags = append(evt.Tags, []string{"p", target, "target"})

	// Add action tag
	evt.Tags = append(evt.Tags, []string{"t", action})

	// Add moderation type tags
	evt.Tags = append(evt.Tags, []string{"t", "group"})
	evt.Tags = append(evt.Tags, []string{"t", "moderation"})

	return evt, nil
}

// ParseGroupMetadataEvent parses a group metadata event
func ParseGroupMetadataEvent(evt *nostr.Event) (*GroupMetadata, error) {
	if evt.Kind != KindGroupMetadata {
		return nil, fmt.Errorf("event is not a group metadata event (kind %d)", evt.Kind)
	}

	var metadata GroupMetadata
	err := json.Unmarshal([]byte(evt.Content), &metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal group metadata: %w", err)
	}

	return &metadata, nil
}

// ParseGroupMessageEvent parses a group message event
func ParseGroupMessageEvent(evt *nostr.Event) (*GroupMessage, error) {
	if evt.Kind != KindGroupMessage {
		return nil, fmt.Errorf("event is not a group message event (kind %d)", evt.Kind)
	}

	var message GroupMessage
	err := json.Unmarshal([]byte(evt.Content), &message)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal group message: %w", err)
	}

	return &message, nil
}

// ParseGroupJoinEvent parses a group join event
func ParseGroupJoinEvent(evt *nostr.Event) (*GroupJoin, error) {
	if evt.Kind != KindGroupJoin {
		return nil, fmt.Errorf("event is not a group join event (kind %d)", evt.Kind)
	}

	var join GroupJoin
	err := json.Unmarshal([]byte(evt.Content), &join)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal group join: %w", err)
	}

	return &join, nil
}

// ParseGroupLeaveEvent parses a group leave event
func ParseGroupLeaveEvent(evt *nostr.Event) (*GroupLeave, error) {
	if evt.Kind != KindGroupLeave {
		return nil, fmt.Errorf("event is not a group leave event (kind %d)", evt.Kind)
	}

	var leave GroupLeave
	err := json.Unmarshal([]byte(evt.Content), &leave)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal group leave: %w", err)
	}

	return &leave, nil
}

// ParseGroupModerationEvent parses a group moderation event
func ParseGroupModerationEvent(evt *nostr.Event) (*GroupModeration, error) {
	if evt.Kind != KindGroupModeration {
		return nil, fmt.Errorf("event is not a group moderation event (kind %d)", evt.Kind)
	}

	var moderation GroupModeration
	err := json.Unmarshal([]byte(evt.Content), &moderation)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal group moderation: %w", err)
	}

	return &moderation, nil
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
	return evt.Kind == KindGroupMetadata ||
		evt.Kind == KindGroupMessage ||
		evt.Kind == KindGroupJoin ||
		evt.Kind == KindGroupLeave ||
		evt.Kind == KindGroupModeration
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
	case KindGroupMetadata:
		return "metadata"
	case KindGroupMessage:
		return "message"
	case KindGroupJoin:
		return "join"
	case KindGroupLeave:
		return "leave"
	case KindGroupModeration:
		return "moderation"
	default:
		return "unknown"
	}
}
