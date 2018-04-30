package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"

	"github.com/hyeonmink/challenges-hyeonmink/apiserver/sessions"
)

//ChatbotHandler handles chatbot
func (ctx *Context) ChatbotHandler(w http.ResponseWriter, r *http.Request) {
	state := &SessionState{}

	_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, &state)
	if err != nil {
		http.Error(w, "error getting session state "+err.Error(), http.StatusForbidden)
		return
	}
	proxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			//reset the scheme and Host of
			//the request URL
			r.URL.Scheme = "http"
			r.URL.Host = ctx.SvcAddr
			j, _ := json.Marshal(state.User)
			r.Header.Add("X-User", string(j))
		},
	}

	proxy.ServeHTTP(w, r)
}
