package handlers

import (
	"io"
	"net/http"

	"encoding/json"

	"path"

	"github.com/hyeonmink/challenges-hyeonmink/apiserver/models/events"
	"github.com/hyeonmink/challenges-hyeonmink/apiserver/models/messages"
	"github.com/hyeonmink/challenges-hyeonmink/apiserver/models/users"
	"github.com/hyeonmink/challenges-hyeonmink/apiserver/sessions"
)

//DefaultLimit of the message is set to be 500
const DefaultLimit = 500

//ChannelsHandler allows users to make new inserto
func (ctx *Context) ChannelsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		state := &SessionState{}
		_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, state)
		if err != nil {
			http.Error(w, "Error1 "+err.Error(), http.StatusUnauthorized)
			return
		}

		chs, err := ctx.MessageStore.GetAllChan(state.User)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		endocer := json.NewEncoder(w)
		endocer.Encode(chs)
	case "POST":
		state := &SessionState{}
		_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, state)
		if err != nil {
			http.Error(w, "Error1 "+err.Error(), http.StatusUnauthorized)
			return
		}

		decoder := json.NewDecoder(r.Body)
		ch := &messages.NewChannel{}
		if err := decoder.Decode(ch); err != nil {
			http.Error(w, "error decoding JSON:", http.StatusBadRequest)
			return
		}

		if err := ch.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		chs, err := ctx.MessageStore.Insert(ch, state.User)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		ctx.Notifier.AddEvent(events.NewChannelEvent(events.NewChannelcreated, chs, state.User))
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		endocer := json.NewEncoder(w)
		endocer.Encode(chs)
	default:
		io.WriteString(w, "Method has to be either POST or GET")
	}

}

// SpecificChannelHandler .
func (ctx *Context) SpecificChannelHandler(w http.ResponseWriter, r *http.Request) {
	state := &SessionState{}
	_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	_, id := path.Split(r.URL.Path)
	chanID := messages.ChannelID(id)

	switch r.Method {
	case "GET":
		msgs, err := ctx.MessageStore.GetRecentMessages(chanID, DefaultLimit, state.User)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		endocer := json.NewEncoder(w)
		endocer.Encode(msgs)

	case "PATCH":
		decoder := json.NewDecoder(r.Body)
		msgs := &messages.ChannelUpdates{}
		if err := decoder.Decode(msgs); err != nil {
			http.Error(w, "error decoding JSON:", http.StatusBadRequest)
			return
		}
		ch, err := ctx.MessageStore.UpdateChan(chanID, msgs, state.User)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx.Notifier.AddEvent(events.NewChannelEvent(events.Channelupdated, ch, state.User))
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		endocer := json.NewEncoder(w)
		endocer.Encode(ch)

	case "DELETE":
		ch, err := ctx.MessageStore.DeleteChan(chanID, state.User)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx.Notifier.AddEvent(events.NewChannelEvent(events.Channeldeleted, ch, state.User))
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		endocer := json.NewEncoder(w)
		endocer.Encode(ch)
	case "LINK":
		anotherUserID := r.Header.Get("Link")
		if len(anotherUserID) != 0 {
			newUser, err := ctx.UserStore.GetByID(users.UserID(anotherUserID))
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ch, err := ctx.MessageStore.AddUser(chanID, state.User, newUser)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx.Notifier.AddEvent(events.NewChannelEvent(events.UserJoinedChannel, ch, newUser))
			io.WriteString(w, "New Member has been added!")
		} else {
			ch, err := ctx.MessageStore.AddUser(chanID, state.User, state.User)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			ctx.Notifier.AddEvent(events.NewChannelEvent(events.UserJoinedChannel, ch, state.User))
			io.WriteString(w, "New Member has been added!")
		}

	case "UNLINK":
		anotherUserID := r.Header.Get("Link")
		if len(anotherUserID) != 0 {
			remove, err := ctx.UserStore.GetByID(users.UserID(anotherUserID))
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ch, err := ctx.MessageStore.RemoveUser(chanID, state.User, remove)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx.Notifier.AddEvent(events.NewChannelEvent(events.UserLeftChannel, ch, state.User))
			io.WriteString(w, "Member has been removed!")
		} else {
			ch, err := ctx.MessageStore.RemoveUser(chanID, state.User, state.User)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			ctx.Notifier.AddEvent(events.NewChannelEvent(events.UserLeftChannel, ch, state.User))
			io.WriteString(w, "Member has been removed!")
		}
	default:
		io.WriteString(w, "You have wrong method")
	}
}

//MessagesHandler .
func (ctx *Context) MessagesHandler(w http.ResponseWriter, r *http.Request) {
	state := &SessionState{}
	_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		newMsg := &messages.NewMessage{}
		if err := decoder.Decode(newMsg); err != nil {
			http.Error(w, "error decoding JSON:", http.StatusBadRequest)
			return
		}
		msg, err := ctx.MessageStore.InsertMessage(newMsg, state.User)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx.Notifier.AddEvent(events.NewMessageEvent(events.NewMessagePosted, msg))
		io.WriteString(w, "New message has been added!")
	default:
		io.WriteString(w, "Method should be POST")
	}
}

//SpecificMessageHandler .
func (ctx *Context) SpecificMessageHandler(w http.ResponseWriter, r *http.Request) {
	state := &SessionState{}
	_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	_, id := path.Split(r.URL.Path)
	MsgID := messages.MessageID(id)

	switch r.Method {
	case "PATCH":
		decoder := json.NewDecoder(r.Body)
		msg := &messages.MessageUpdates{}
		if err := decoder.Decode(msg); err != nil {
			http.Error(w, "error decoding JSON:", http.StatusBadRequest)
			return
		}
		m, err := ctx.MessageStore.UpdateMessage(MsgID, msg, state.User)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx.Notifier.AddEvent(events.NewMessageEvent(events.MessageUpdated, m))
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		endocer := json.NewEncoder(w)
		endocer.Encode(m)
	case "DELETE":
		msg, err := ctx.MessageStore.DeleteMessage(MsgID, state.User)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx.Notifier.AddEvent(events.NewMessageEvent(events.MessageDeleted, msg))
		io.WriteString(w, "Message has been deleted!")
	default:
		io.WriteString(w, "Method should be either PATCH or DELETE!")
	}
}
