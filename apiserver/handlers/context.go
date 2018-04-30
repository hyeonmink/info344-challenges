package handlers

import (
	"github.com/hyeonmink/challenges-hyeonmink/apiserver/models/events"
	"github.com/hyeonmink/challenges-hyeonmink/apiserver/models/messages"
	"github.com/hyeonmink/challenges-hyeonmink/apiserver/models/users"
	"github.com/hyeonmink/challenges-hyeonmink/apiserver/sessions"
)

// Context contains the stores for the server
type Context struct {
	SessionKey   string
	SessionStore sessions.Store
	UserStore    users.Store
	MessageStore messages.Store
	Notifier     *events.Notifier
	SvcAddr      string
}
