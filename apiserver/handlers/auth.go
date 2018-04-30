package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/hyeonmink/challenges-hyeonmink/apiserver/models/events"
	"github.com/hyeonmink/challenges-hyeonmink/apiserver/models/users"
	"github.com/hyeonmink/challenges-hyeonmink/apiserver/sessions"
)

//UsersHandler allows new users to sign-up (POST) or returns all users
func (ctx *Context) UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		user := &users.NewUser{}
		if err := decoder.Decode(user); err != nil {
			http.Error(w, "error decoding JSON:", http.StatusBadRequest)
			return
		}

		if err := user.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//Ensure there isn't already a user in the UserStore with the same email address
		if _, err := ctx.UserStore.GetByEmail(user.Email); err == nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//Ensure there isn't already a user in the UserStore with the same user name
		if _, err := ctx.UserStore.GetByUserName(user.UserName); err == nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//Insert the new user into the UserStore
		u, err := ctx.UserStore.Insert(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//Begin a new session
		if _, err := sessions.BeginSession(ctx.SessionKey, ctx.SessionStore, u, w); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx.Notifier.AddEvent(events.NewUserEvent(events.NewUsercreated, u))
		//Respond to the client with the models.User struct encoded as a JSON object
		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		encoder := json.NewEncoder(w)
		encoder.Encode(user)

	case "GET":
		//Get all users from the UserStore and write them to the response as a JSON-encoded array
		users, err := ctx.UserStore.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		encoder := json.NewEncoder(w)
		encoder.Encode(users)

	default:
		http.Error(w, "request method should be either GET or POST", http.StatusMethodNotAllowed)
		return
	}

}

//SessionsHandler allows existing users to sign-in
func (ctx *Context) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	//The request method must be "POST"
	if r.Method == "POST" {
		// Decode the request body into a models.Credentials struct
		decoder := json.NewDecoder(r.Body)
		cred := &users.Credentials{}
		if err := decoder.Decode(cred); err != nil {
			http.Error(w, "error decoding JSON:"+err.Error(), http.StatusBadRequest)
			return
		}

		// Get the user with the provided email from the UserStore; if not found, respond with an http.StatusUnauthorized
		user, err := ctx.UserStore.GetByEmail(cred.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		state := &SessionState{
			BeganAt:    time.Now(),
			ClientAddr: r.RemoteAddr,
			User:       user,
		}

		// Authenticate the user using the provided password; if that fails, respond with an http.StatusUnauthorized
		if err := user.Authenticate(cred.Password); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}

		// Begin a new session
		_, err = sessions.BeginSession(ctx.SessionKey, ctx.SessionStore, state, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Respond to the client with the models.User struct encoded as a JSON object
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		encoder := json.NewEncoder(w)
		encoder.Encode(user)
	} else {
		http.Error(w, "request method should be POST", http.StatusMethodNotAllowed)
		return
	}
}

//SessionsMineHandler allows authenticated users to sign-out
func (ctx *Context) SessionsMineHandler(w http.ResponseWriter, r *http.Request) {
	// The request method must be "DELETE"
	if r.Method == "DELETE" {
		// End the session
		_, err := sessions.EndSession(r, ctx.SessionKey, ctx.SessionStore)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		// Respond to the client with a simple message saying that the user has been signed out
		io.WriteString(w, "user signed out")

	} else {
		http.Error(w, "request method should be DELETE", http.StatusMethodNotAllowed)
		return
	}
}

//UsersMeHanlder .
func (ctx *Context) UsersMeHanlder(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Get the session state
		sessionsState := &SessionState{}
		_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, sessionsState)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		// Respond to the client with the session state's User field, encoded as a JSON object
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		encoder := json.NewEncoder(w)
		encoder.Encode(sessionsState)
	} else if r.Method == "PATCH" {
		decoder := json.NewDecoder(r.Body)
		user := &users.UserUpdates{}
		if err := decoder.Decode(user); err != nil {
			http.Error(w, "error decoding JSON:", http.StatusBadRequest)
			return
		}

		sessionsState := &SessionState{}
		_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, sessionsState)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		err = ctx.UserStore.Update(user, sessionsState.User)
		if err != nil {
			http.Error(w, "error updating:"+err.Error(), http.StatusBadRequest)
			return
		}

		sid, _ := sessions.GetSessionID(r, ctx.SessionKey)
		sessionsState.User.FirstName = user.FirstName
		sessionsState.User.LastName = user.LastName
		if err := ctx.SessionStore.Save(sid, sessionsState); err != nil {
			http.Error(w, "unable to save session state: "+err.Error(), http.StatusInternalServerError)
			return
		}
		io.WriteString(w, "user information has been updated")

	} else {
		http.Error(w, "request method should be either GET or PATCH", http.StatusMethodNotAllowed)
		return
	}
}
