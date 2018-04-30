package messages

import (
	"fmt"
	"time"

	"github.com/hyeonmink/challenges-hyeonmink/apiserver/models/users"
)

//MessageID defines the type for user IDs
type MessageID string

//Message represents a Message in the database
type Message struct {
	ID        MessageID    `json:"id" bson:"_id"`
	ChannelID ChannelID    `json:"ChannelID"`
	Body      string       `json:"body"`
	CreatorID users.UserID `json:"creatorID"`
	CreatedAt time.Time    `json:"createdAt"`
	EditedAt  *time.Time   `json:"editedAt"`
}

type NewMessage struct {
	ChannelID ChannelID `json:"ChannelID"`
	Body      string    `json:"body"`
}

type MessageUpdates struct {
	Body string `json:"body"`
}

//Validate make sure NewMessage has all fields
func (nm *NewMessage) Validate() error {
	if len(nm.ChannelID) == 0 {
		return fmt.Errorf("error: no channelID")
	}
	if len(nm.Body) == 0 {
		return fmt.Errorf("error: no body")
	}
	return nil
}

//Validate make sure MessageUpdates has all fields
func (mu *MessageUpdates) Validate() error {
	if len(mu.Body) == 0 {
		return fmt.Errorf("error: no body")
	}
	return nil
}

//toMessage converts the NewMessage to a Message
func (nm *NewMessage) toMessage(currentUser *users.User) *Message {
	msg := &Message{
		ChannelID: nm.ChannelID,
		Body:      nm.Body,
		CreatorID: currentUser.ID,
		CreatedAt: time.Now(),
	}
	return msg
}
