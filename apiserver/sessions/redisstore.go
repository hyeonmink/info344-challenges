package sessions

import (
	"time"

	"encoding/json"

	"gopkg.in/redis.v5"
)

//redisKeyPrefix is the prefix we will use for keys
//related to session IDs. This keeps session ID keys
//separate from other keys in the shared redis key
//namespace.
const redisKeyPrefix = "sid:"

//RedisStore represents a session.Store backed by redis.
type RedisStore struct {
	//Redis client used to talk to redis server.
	Client *redis.Client
	//Used for key expiry time on redis.
	SessionDuration time.Duration
}

//NewRedisStore constructs a new RedisStore, using the provided client and
//session duration. If the `client`` is nil, it will be set to redis.NewClient()
//pointing at a local redis instance. If `sessionDuration`` is negative, it will
//be set to `DefaultSessionDuration`.
func NewRedisStore(client *redis.Client, sessionDuration time.Duration) *RedisStore {
	//set defaults for parameters
	//if `client` is nil, set it to a redis.NewClient()
	//pointing at a redis instance on the same machine
	//i.e., Addr is "127.0.0.1"
	ropts := redis.Options{
		Addr: "localhost:6379",
	}
	if client == nil {
		client = redis.NewClient(&ropts)
	}

	//if `sessionDuration` is < 0
	//set it to DefaultSessionDuration
	if sessionDuration < 0 {
		sessionDuration = DefaultSessionDuration
	}
	//return a new RedisStore with the Client field set to `client`
	//and the SessionDuration field set to `sessionDuration`
	return &RedisStore{client, sessionDuration}
}

//Store implementation

//Save associates the provided `state` data with the provided `sid` in the store.
func (rs *RedisStore) Save(sid SessionID, state interface{}) error {
	//encode the `state` into JSON
	data, _ := json.Marshal(state)
	//use the redis client's Set() method, using `sid.getRedisKey()`
	//as the key, the JSON as the data, and the store's session duration
	//as the expiration
	StatusCmd := rs.Client.Set(sid.getRedisKey(), data, rs.SessionDuration)
	//Set() returns a StatusCmd, which has an .Err() method that will
	//report any error that occurred; return the result of that method
	return StatusCmd.Err()
}

//Get retrieves the previously saved data for the session id,
//and populates the `state` parameter with it. This will also
//reset the data's time to live in the store.
func (rs *RedisStore) Get(sid SessionID, state interface{}) error {
	//use the .Get() method to get the data associated
	//with the key `sid.getRedisKey()`
	stringCmd := rs.Client.Get(sid.getRedisKey())
	//if the Get command returned an error,
	//return ErrStateNotFound if the error == redis.Nil
	//otherwise return the error
	if stringCmd.Err() != nil {
		if stringCmd.Err() == redis.Nil {
			return ErrStateNotFound
		}
		return stringCmd.Err()
	}

	//get the returned bytes and Unmarshal them into
	//the `state` parameter
	//if you get an error, return it
	bytes, _ := stringCmd.Bytes()
	err := json.Unmarshal(bytes, state)
	if err != nil {
		return err
	}
	//use the .Expire() command to reset the expiry duration
	//to the store's session duration
	rs.Client.Expire(sid.getRedisKey(), rs.SessionDuration)

	//for EXTRA CREDIT use the Pipeline feature
	//to do the .Get() and .Expire() commands
	//in just one round-trip!

	return nil
}

//Delete deletes all data associated with the session id from the store.
func (rs *RedisStore) Delete(sid SessionID) error {
	//use the .Del() method to delete the data associated
	//with the key `sid.getRedisKey()`, and use .Err()
	//to report any errors that occurred
	state := rs.Client.Del(sid.getRedisKey())
	return state.Err()
}

//returns the key to use in redis
func (sid SessionID) getRedisKey() string {
	return redisKeyPrefix + sid.String()
}
