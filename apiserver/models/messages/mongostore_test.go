package messages

import (
	"testing"

	"github.com/hyeonmink/challenges-hyeonmink/apiserver/models/users"
	mgo "gopkg.in/mgo.v2"
)

func Test(t *testing.T) {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		t.Fatalf("error %v\n", err)
	}

	_, err = session.DB("TestDatabase").C("MessageCollection").RemoveAll(nil)
	_, err = session.DB("TestDatabase").C("ChannelCollection").RemoveAll(nil)
	_, err = session.DB("test").C("users").RemoveAll(nil)

	MongoStore := &MongoStore{
		Session:               session,
		DatabaseName:          "344api",
		ChannelCollectionName: "channels",
		MessageCollectionName: "messages",
	}

	userStore := &users.MongoStore{
		Session:        session,
		DatabaseName:   "344api",
		CollectionName: "users",
	}

	pubNewCh := &NewChannel{
		Name:    "PubChan",
		Descr:   "PubChanDescr",
		Members: []users.UserID{},
		Private: false,
	}

	privNewCh := &NewChannel{
		Name:    "privNewCh",
		Descr:   "privNewChDescr",
		Members: []users.UserID{},
		Private: true,
	}

	pubNewUser := &users.NewUser{
		Email:        "test@test.com",
		Password:     "Password",
		PasswordConf: "Password",
		UserName:     "pubUser",
		FirstName:    "PubFirst",
		LastName:     "PubLast",
	}
	pubUser, err := userStore.Insert(pubNewUser)
	if err != nil {
		t.Errorf("error inserting pub user: %v/n", err)
	}
	pubCh, err := MongoStore.Insert(pubNewCh, pubUser)

	privNewUser := &users.NewUser{
		Email:        "test@test.com",
		Password:     "Password",
		PasswordConf: "Password",
		UserName:     "privUser",
		FirstName:    "privFirst",
		LastName:     "privLast",
	}
	privUser, err := userStore.Insert(privNewUser)
	if err != nil {
		t.Errorf("error inserting priv user: %v\n", err)
	}
	privCh, err := MongoStore.Insert(privNewCh, privUser)

	_, err = MongoStore.GetChannelByID(pubCh.ID, pubUser)
	if err != nil {
		t.Errorf("error finding pub channel: %v\n", err)
	}

	_, err = MongoStore.GetChannelByID(privCh.ID, privUser)
	if err != nil {
		t.Errorf("error finding priv channel: %v\n", err)
	}

	chs, err := MongoStore.GetAllChan(pubUser)
	if err != nil {
		t.Errorf("error get all channels: %v\n", err)
	}
	if len(chs) == 0 {
		t.Errorf("number of channel should be 1: %v\n", err)
	}

	chs, err = MongoStore.GetAllChan(privUser)
	if err != nil {
		t.Errorf("error get all channels: %v\n", err)
	}
	if len(chs) == 0 {
		t.Errorf("number of channel should be 1: %v\n", err)
	}

	NewPubMsg := &NewMessage{
		ChannelID: pubCh.ID,
		Body:      "Test Message",
	}

	NewPrivMsg := &NewMessage{
		ChannelID: privCh.ID,
		Body:      "Test Message",
	}

	PubMsg, err := MongoStore.InsertMessage(NewPubMsg, pubUser)
	if err != nil {
		t.Errorf("error inserting new message %v:\n", err)
	}

	PrivMsg, err := MongoStore.InsertMessage(NewPrivMsg, privUser)
	if err != nil {
		t.Errorf("error inserting new message %v:\n", err)
	}

	UpdateMsg := &MessageUpdates{
		Body: "Update Msg",
	}

	_, err = MongoStore.UpdateMessage(PubMsg.ID, UpdateMsg, pubUser)
	if err != nil {
		t.Errorf("error updating msg %v\n", err)
	}

	_, err = MongoStore.UpdateMessage(PrivMsg.ID, UpdateMsg, privUser)
	if err != nil {
		t.Errorf("error updating msg %v\n", err)
	}

	// _, err = MongoStore.DeleteMessage(PrivMsg.ID, pubUser)
	// if err == nil {
	// 	t.Errorf("pub user shouldn't have permission to delete priv : %v\n", err)
	// }

	// _, err = MongoStore.DeleteMessage(PubMsg.ID, pubUser)
	// if err != nil {
	// 	t.Errorf("err deleting pub message %v\n", err)
	// }

	// _, err = MongoStore.DeleteMessage(PrivMsg.ID, privUser)
	// if err != nil {
	// 	t.Errorf("err deleting priv message %v\n", err)
	// }

	updateChan := &ChannelUpdates{
		Name:    "Channel Updates",
		Descr:   "Updated Description",
		Private: true,
	}

	updateChan2, err := MongoStore.UpdateChan(privCh.ID, updateChan, privUser)
	if err != nil {
		t.Error(privCh.ID)
		t.Errorf("error updating channel : %v\n", err)
	}
	if updateChan2.Name != updateChan.Name {
		t.Errorf("error updating Name : %v\n", err)
	}
	if updateChan2.Descr != updateChan.Descr {
		t.Errorf("error updating Descr : %v\n", err)
	}
	if updateChan2.Private != updateChan.Private {
		t.Errorf("error updating Private : %v\n", err)
	}

	ch2, err := MongoStore.AddUser(pubCh.ID, privUser, privUser)
	if err != nil {
		t.Errorf("error adding a user to public channel: %v\n", err) //v
	}
	if len(ch2.Members) != 2 {
		t.Errorf("Member should be 2 but : %v\n", len(ch2.Members)) //v
	}

	// ch2, err = MongoStore.RemoveUser(pubCh.ID, privUser, privUser)
	// if err != nil {
	// 	t.Errorf("error deleting a user to public channel: %v\n", err) //v
	// }
	// if len(ch2.Members) != 1 {
	// 	t.Errorf("Member should be 1 but : %v\n", len(ch2.Members)) //v
	// }

	NewPrivMsg = &NewMessage{
		ChannelID: privCh.ID,
		Body:      "Test Message",
	}

	PrivMsg, err = MongoStore.InsertMessage(NewPrivMsg, privUser)
	if err != nil {
		t.Errorf("error inserting new message %v:\n", err)
	}

	PrivMsg, err = MongoStore.InsertMessage(NewPrivMsg, privUser)
	if err != nil {
		t.Errorf("error inserting new message %v:\n", err)
	}

	all, err := MongoStore.GetRecentMessages(privCh.ID, 2, privUser)
	if err != nil {
		t.Errorf("error getting 2 messages : %v\n", err)
	}
	if len(all) != 2 {
		t.Errorf("error getting 2 messages : %v\n", err)
	}

	_, err = MongoStore.DeleteMessage(PrivMsg.ID, privUser)
	if err != nil {
		t.Errorf("error deleting a messages : %v\n", err)
	}

	_, err = MongoStore.DeleteChan(privCh.ID, privUser)
	if err != nil {
		t.Errorf("error deleting a messages : %v\n", err)
	}

	// _, err = MongoStore.RemoveUser(pubCh.ID, pubUser)
	// if err != nil {
	// 	t.Errorf("error removing user : %v\n", err)
	// }
}
