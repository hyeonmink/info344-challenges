package events

import (
	"github.com/hyeonmink/challenges-hyeonmink/apiserver/models/messages"
	"github.com/hyeonmink/challenges-hyeonmink/apiserver/models/users"
)

type UserEventType string
type MessageEventType string
type ChannelEventType string

const (
	NewUsercreated    UserEventType    = "NewUsercreated"
	NewChannelcreated ChannelEventType = "newChannelcreated"
	Channelupdated    ChannelEventType = "Channelupdated"
	Channeldeleted    ChannelEventType = "Channeldeleted"
	UserJoinedChannel ChannelEventType = "UserJoinedChannel"
	UserLeftChannel   ChannelEventType = "UserLeftChannel"
	NewMessagePosted  MessageEventType = "NewMessagePosted"
	MessageUpdated    MessageEventType = "MessageUpdated"
	MessageDeleted    MessageEventType = "MessageDeleted"
)

//UserEvent UserEvent type
type UserEvent struct {
	Type UserEventType `json:"type"`
	User *users.User   `json:"user"`
}

//MessageEvent MessageEvent type
type MessageEvent struct {
	Type    MessageEventType  `json:"type"`
	Message *messages.Message `json:"message"`
}

//ChannelEvent ChannelEvent type
type ChannelEvent struct {
	Type    ChannelEventType  `json:"type"`
	Channel *messages.Channel `json:"channel"`
	User    *users.User       `json:"user"`
}

//NewMessageEvent creates new messageEvent and return it
func NewMessageEvent(eventType MessageEventType, message *messages.Message) *MessageEvent {
	messageEvent := &MessageEvent{
		Type:    eventType,
		Message: message,
	}
	return messageEvent
}

//NewChannelEvent creates new channelEvent and return it
func NewChannelEvent(eventType ChannelEventType, channel *messages.Channel, user *users.User) *ChannelEvent {
	channelEvent := &ChannelEvent{
		Type:    eventType,
		Channel: channel,
		User:    user,
	}
	return channelEvent
}

//NewUserEvent creates new userEvent and return it
func NewUserEvent(eventType UserEventType, user *users.User) *UserEvent {
	userEvent := &UserEvent{
		Type: eventType,
		User: user,
	}
	return userEvent
}
