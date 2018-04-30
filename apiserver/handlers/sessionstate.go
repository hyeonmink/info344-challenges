package handlers

import (
	"time"

	"github.com/hyeonmink/challenges-hyeonmink/apiserver/models/users"
)

// SessionState shows a user and current session
type SessionState struct {
	BeganAt    time.Time   `json:"BeganAt"`
	ClientAddr string      `json:"ClientAddr"`
	User       *users.User `json:"User"`
}
