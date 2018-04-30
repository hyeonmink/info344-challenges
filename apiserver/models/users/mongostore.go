package users

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//MongoStore is an implementation of UserStore
//backed by an in-memory slice. This should only
//be used for automated testing. This is not
//safe for concurrent access
type MongoStore struct {
	Session        *mgo.Session
	DatabaseName   string
	CollectionName string
}

const defaultAddr = "127.0.0.1:27017"
const databaseNAme = "defaultDatabase"

//NewMongoStore creates newNewMongoStoreMon
func NewMongoStore(session *mgo.Session, databaseName string) (*MongoStore, error) {
	if session == nil {
		var err error
		session, err = mgo.Dial(defaultAddr)
		if err != nil {
			return nil, err
		}
	}

	if len(databaseName) <= 0 {
		databaseName = databaseNAme
	}

	return &MongoStore{
		Session:        session,
		DatabaseName:   databaseName,
		CollectionName: "users",
	}, nil

}

//GetAll returns all users
func (ms *MongoStore) GetAll() ([]*User, error) {
	users := []*User{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).Find(nil).All(&users)
	if err != nil {
		return nil, err
	}
	return users, err
}

//GetByID returns the User with the given ID
func (ms *MongoStore) GetByID(id UserID) (*User, error) {
	user := &User{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).FindId(id).One(user)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

//GetByEmail returns the User with the given email
func (ms *MongoStore) GetByEmail(email string) (*User, error) {
	user := &User{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).Find(bson.M{"email": email}).One(user)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

//GetByUserName returns the User with the given user name
func (ms *MongoStore) GetByUserName(name string) (*User, error) {
	user := &User{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).Find(bson.M{"username": name}).One(user)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

//Insert inserts a new NewUser into the database
//and return a User with new ID, or an error
func (ms *MongoStore) Insert(newUser *NewUser) (*User, error) {
	u, err := newUser.ToUser()
	if err != nil {
		return nil, err
	}
	u.ID = UserID(bson.NewObjectId().Hex())
	err = ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).Insert(u)
	return u, err
}

//Update applies UserUpdates to the currentUser
func (ms *MongoStore) Update(updates *UserUpdates, currentuser *User) error {
	col := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName)
	err := col.UpdateId(currentuser.ID, bson.M{"$set": updates})
	if err != nil {
		return err
	}
	currentuser.FirstName = updates.FirstName
	currentuser.LastName = updates.LastName

	return nil
}
