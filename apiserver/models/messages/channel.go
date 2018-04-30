package messages

import (
	"errors"
	"time"

	"github.com/hyeonmink/challenges-hyeonmink/apiserver/models/users"
)

//ChannelID defines the type for user IDs
type ChannelID string

//Channel represents a Channel in the database
type Channel struct {
	ID        ChannelID      `json:"id" bson:"_id"`
	Name      string         `json:"name"`
	Descr     string         `json:"descr"`
	CreatedAt time.Time      `json:"createdAt"`
	CreatorID users.UserID   `json:"creatorID"`
	Members   []users.UserID `json:"members"`
	Private   bool           `json:"private"`
}

//NewChannel represents a new NewChannel for an website
type NewChannel struct {
	Name    string         `json:"name"`
	Descr   string         `json:"descr"`
	Members []users.UserID `json:"members"`
	Private bool           `json:"private"`
}

//ChannelUpdates represents updates one can make to a channel
type ChannelUpdates struct {
	Name    string `json:"name"`
	Descr   string `json:"descr"`
	Private bool   `json:"private"`
}

//ToChannel converts the newChannel to a Channel
func (newChan *NewChannel) ToChannel(currentUser *users.User) *Channel {
	channel := &Channel{
		Name:      newChan.Name,
		Descr:     newChan.Descr,
		CreatedAt: time.Now(),
		CreatorID: currentUser.ID,
		Members:   newChan.Members,
		Private:   newChan.Private,
	}
	return channel
}

//Validate new chan
func (newChan *NewChannel) Validate() error {
	if len(newChan.Name) == 0 {
		return errors.New("Channel Name has to be non-zero length")
	}

	if len(newChan.Descr) == 0 {
		return errors.New("Channel Description has to be non-zero length")
	}
	return nil
}

func (ch *Channel) Contains(user *users.User) bool {
	for _, id := range ch.Members {
		if id == user.ID {
			return true
		}
	}
	return false
}

func (ch *Channel) AddMember(user *users.User) {
	if !ch.Contains(user) {
		ch.Members = append(ch.Members, user.ID)
	}
}

func (ch *Channel) RemoveUser(user *users.User) {
	for i, id := range ch.Members {
		if id == user.ID {
			ch.Members = ch.Members[:i+copy(ch.Members[i:], ch.Members[i+1:])]
		}
	}
}
