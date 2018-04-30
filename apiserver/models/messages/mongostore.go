package messages

import (
	"errors"
	"time"

	"github.com/hyeonmink/challenges-hyeonmink/apiserver/models/users"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//MongoStore is database of messages
type MongoStore struct {
	Session               *mgo.Session
	DatabaseName          string
	ChannelCollectionName string
	MessageCollectionName string
}

const defaultAddr = "127.0.0.1:27017"
const defaultDBName = "defaultDatabase"
const defaultChanColName = "defaultChanColName"
const defaultMsgColName = "defaultMsgColName"

//NewMongoStore creates newNewMongoStoreMon
func NewMongoStore(session *mgo.Session, databaseName string, chanColName string, msgColName string) (*MongoStore, error) {
	if session == nil {
		var err error
		session, err = mgo.Dial(defaultAddr)
		if err != nil {
			return nil, err
		}
	}

	if len(databaseName) == 0 {
		databaseName = defaultDBName
	}

	if len(chanColName) == 0 {
		chanColName = defaultChanColName
	}

	if len(msgColName) == 0 {
		msgColName = defaultMsgColName
	}

	return &MongoStore{
		Session:               session,
		DatabaseName:          databaseName,
		ChannelCollectionName: chanColName,
		MessageCollectionName: msgColName,
	}, nil
}

func (ms *MongoStore) GetChannelByID(id ChannelID, currentUser *users.User) (*Channel, error) {
	ch := &Channel{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.ChannelCollectionName).FindId(string(id)).One(&ch)
	if err != nil {
		return ch, err
	}

	if ch.Private && !ch.Contains(currentUser) {
		return nil, errors.New("User is not authorized")
	}
	return ch, nil
}

func (ms *MongoStore) GetMessageByID(id MessageID, currentUser *users.User) (*Message, error) {
	msg := &Message{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.MessageCollectionName).FindId(string(id)).One(msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

//GetAllChan returns all channels
func (ms *MongoStore) GetAllChan(currentUser *users.User) ([]*Channel, error) {
	channels := []*Channel{}
	query := bson.M{
		"$or": []map[string]interface{}{
			{"private": false},
			{"members": currentUser.ID},
		},
	}

	err := ms.Session.DB(ms.DatabaseName).C(ms.ChannelCollectionName).Find(query).All(&channels)
	if err != nil {
		return nil, err
	}
	return channels, nil
}

//Insert inserts a new NewChannel into the database
//and return a NewChannel with new ID, or an error
func (ms *MongoStore) Insert(newChan *NewChannel, user *users.User) (*Channel, error) {
	ch := newChan.ToChannel(user)
	ch.ID = ChannelID(bson.NewObjectId().Hex())
	ch.AddMember(user)
	err := ms.Session.DB(ms.DatabaseName).C(ms.ChannelCollectionName).Insert(ch)
	return ch, err
}

//GetRecentMessages returns slice of messages or error
func (ms *MongoStore) GetRecentMessages(chanID ChannelID, num int, currentUser *users.User) ([]*Message, error) {
	if num == 0 {
		return nil, errors.New("Number must be bigger than 0")
	}

	ch, err := ms.GetChannelByID(chanID, currentUser)
	if err != nil {
		return nil, err
	}
	msgs := []*Message{}
	err = ms.Session.DB(ms.DatabaseName).C(ms.MessageCollectionName).Find(bson.M{"channelid": string(ch.ID)}).Sort("_id").Limit(num).All(&msgs)
	return msgs, err
}

//UpdateChan updates an existing channel
func (ms *MongoStore) UpdateChan(chanID ChannelID, updates *ChannelUpdates, currentUser *users.User) (*Channel, error) {
	ch, err := ms.GetChannelByID(chanID, currentUser)
	if err != nil {
		return nil, err
	}

	if ch.CreatorID != currentUser.ID {
		return nil, errors.New("User is not authorized")
	}
	ch.Name = updates.Name
	ch.Descr = updates.Descr
	ch.Private = updates.Private

	err = ms.Session.DB(ms.DatabaseName).C(ms.ChannelCollectionName).UpdateId(chanID, ch)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

//DeleteChan deletes channel from the database
func (ms *MongoStore) DeleteChan(chanID ChannelID, currentUser *users.User) (*Channel, error) {
	ch, err := ms.GetChannelByID(chanID, currentUser)
	if err != nil {
		return nil, err
	}

	if ch.CreatorID != currentUser.ID {
		return nil, errors.New("User is not authorized")
	}

	_, err = ms.Session.DB(ms.DatabaseName).C(ms.MessageCollectionName).RemoveAll(bson.M{"channelid": string(chanID)}) //
	if err != nil {
		return nil, err
	}

	return ch, ms.Session.DB(ms.DatabaseName).C(ms.ChannelCollectionName).RemoveId(chanID)
}

//AddUser adds the user to the members list
func (ms *MongoStore) AddUser(chanID ChannelID, user *users.User, newUser *users.User) (*Channel, error) {
	ch, err := ms.GetChannelByID(chanID, user)
	if err != nil {
		return nil, err
	}

	if ch.Private && ch.CreatorID != user.ID {
		return nil, errors.New("User is not authorized")
	}
	ch.AddMember(newUser)
	err = ms.Session.DB(ms.DatabaseName).C(ms.ChannelCollectionName).UpdateId(ch.ID, bson.M{"$addToSet": bson.M{"members": newUser.ID}})

	return ch, err
}

//RemoveUser removes the user from the members list
func (ms *MongoStore) RemoveUser(chanID ChannelID, user *users.User, remove *users.User) (*Channel, error) {
	ch, err := ms.GetChannelByID(chanID, user)
	if err != nil {
		return nil, err
	}

	if ch.Private && ch.CreatorID != user.ID {
		return nil, errors.New("User is not authorized")
	}

	ch.RemoveUser(remove)
	err = ms.Session.DB(ms.DatabaseName).C(ms.ChannelCollectionName).UpdateId(ch.ID, bson.M{"$pull": bson.M{"members": remove.ID}})

	return ch, err
}

//InsertMessage inserts new message to the channel
func (ms *MongoStore) InsertMessage(newMessage *NewMessage, user *users.User) (*Message, error) {
	ch, err := ms.GetChannelByID(newMessage.ChannelID, user)
	if err != nil {
		return nil, err
	}

	if !(ch.Private || ch.Contains(user)) {
		return nil, errors.New("User is not authorized")
	}

	msg := newMessage.toMessage(user)
	msg.ID = MessageID(bson.NewObjectId().Hex())
	err = ms.Session.DB(ms.DatabaseName).C(ms.MessageCollectionName).Insert(msg)
	return msg, err
}

//UpdateMessage updates the message in the channel
func (ms *MongoStore) UpdateMessage(messageID MessageID, updates *MessageUpdates, currentUser *users.User) (*Message, error) {
	msg, err := ms.GetMessageByID(messageID, currentUser)
	if err != nil {
		return msg, err
	}

	if msg.CreatorID != currentUser.ID {
		return nil, errors.New("User is not authorized")
	}

	msg.Body = updates.Body
	now := time.Now()
	msg.EditedAt = &now
	err = ms.Session.DB(ms.DatabaseName).C(ms.MessageCollectionName).UpdateId(messageID, bson.M{"$set": msg})
	if err != nil {
		return nil, err
	}
	return msg, err
}

//DeleteMessage deletes the message from the channel
func (ms *MongoStore) DeleteMessage(messageID MessageID, currentUser *users.User) (*Message, error) {
	msg, err := ms.GetMessageByID(messageID, currentUser)
	if err != nil {
		return msg, err
	}

	if msg.CreatorID != currentUser.ID {
		return nil, errors.New("User is not authorized")
	}
	return msg, ms.Session.DB(ms.DatabaseName).C(ms.MessageCollectionName).RemoveId(messageID)
}
