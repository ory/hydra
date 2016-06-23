package slack

// OutgoingMessage is used for the realtime API, and seems incomplete.
type OutgoingMessage struct {
	ID      int    `json:"id"`
	Channel string `json:"channel,omitempty"`
	Text    string `json:"text,omitempty"`
	Type    string `json:"type,omitempty"`
}

// Message is an auxiliary type to allow us to have a message containing sub messages
type Message struct {
	Msg
	SubMessage *Msg `json:"message,omitempty"`
}

// Msg contains information about a slack message
type Msg struct {
	// Basic Message
	Type        string       `json:"type,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	User        string       `json:"user,omitempty"`
	Text        string       `json:"text,omitempty"`
	Timestamp   string       `json:"ts,omitempty"`
	IsStarred   bool         `json:"is_starred,omitempty"`
	PinnedTo    []string     `json:"pinned_to, omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	Edited      *Edited      `json:"edited,omitempty"`

	// Message Subtypes
	SubType string `json:"subtype,omitempty"`

	// Hidden Subtypes
	Hidden           bool   `json:"hidden,omitempty"`     // message_changed, message_deleted, unpinned_item
	DeletedTimestamp string `json:"deleted_ts,omitempty"` // message_deleted
	EventTimestamp   string `json:"event_ts,omitempty"`

	// bot_message (https://api.slack.com/events/message/bot_message)
	BotID    string `json:"bot_id,omitempty"`
	Username string `json:"username,omitempty"`
	Icons    *Icon  `json:"icons,omitempty"`

	// channel_join, group_join
	Inviter string `json:"inviter,omitempty"`

	// channel_topic, group_topic
	Topic string `json:"topic,omitempty"`

	// channel_purpose, group_purpose
	Purpose string `json:"purpose,omitempty"`

	// channel_name, group_name
	Name    string `json:"name,omitempty"`
	OldName string `json:"old_name,omitempty"`

	// channel_archive, group_archive
	Members []string `json:"members,omitempty"`

	// file_share, file_comment, file_mention
	File *File `json:"file,omitempty"`

	// file_share
	Upload bool `json:"upload,omitempty"`

	// file_comment
	Comment *Comment `json:"comment,omitempty"`

	// pinned_item
	ItemType string `json:"item_type,omitempty"`

	// https://api.slack.com/rtm
	ReplyTo int    `json:"reply_to,omitempty"`
	Team    string `json:"team,omitempty"`

	// reactions
	Reactions []ItemReaction `json:"reactions,omitempty"`
}

// Icon is used for bot messages
type Icon struct {
	IconURL   string `json:"icon_url,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
}

// Edited indicates that a message has been edited.
type Edited struct {
	User      string `json:"user,omitempty"`
	Timestamp string `json:"ts,omitempty"`
}

// Event contains the event type
type Event struct {
	Type string `json:"type,omitempty"`
}

// Ping contains information about a Ping Event
type Ping struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

// Pong contains information about a Pong Event
type Pong struct {
	Type    string `json:"type"`
	ReplyTo int    `json:"reply_to"`
}

// NewOutgoingMessage prepares an OutgoingMessage that the user can
// use to send a message. Use this function to properly set the
// messageID.
func (rtm *RTM) NewOutgoingMessage(text string, channel string) *OutgoingMessage {
	id := rtm.idGen.Next()
	return &OutgoingMessage{
		ID:      id,
		Type:    "message",
		Channel: channel,
		Text:    text,
	}
}

// NewTypingMessage prepares an OutgoingMessage that the user can
// use to send as a typing indicator. Use this function to properly set the
// messageID.
func (rtm *RTM) NewTypingMessage(channel string) *OutgoingMessage {
	id := rtm.idGen.Next()
	return &OutgoingMessage{
		ID:      id,
		Type:    "typing",
		Channel: channel,
	}
}
