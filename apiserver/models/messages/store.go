package messages

import "github.com/hyeonmink/challenges-hyeonmink/apiserver/models/users"

//Store represents an abstract store for messages.channel objects.
type Store interface {
	//GetAllChan returns al chan
	GetAllChan(user *users.User) ([]*Channel, error)

	Insert(newChan *NewChannel, currentUser *users.User) (*Channel, error)

	GetRecentMessages(chanID ChannelID, num int, currentUser *users.User) ([]*Message, error)

	UpdateChan(chanID ChannelID, updatedChan *ChannelUpdates, currentUser *users.User) (*Channel, error)

	DeleteChan(chanID ChannelID, currentUser *users.User) (*Channel, error)

	AddUser(chanID ChannelID, user *users.User, newUser *users.User) (*Channel, error)

	RemoveUser(chanID ChannelID, user *users.User, removeUser *users.User) (*Channel, error)

	InsertMessage(newMessage *NewMessage, currentUser *users.User) (*Message, error)

	UpdateMessage(messageID MessageID, updateMessage *MessageUpdates, currentUser *users.User) (*Message, error)

	DeleteMessage(messageID MessageID, currentUser *users.User) (*Message, error)
}
