package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hyeonmink/challenges-hyeonmink/apiserver/handlers"
	"github.com/hyeonmink/challenges-hyeonmink/apiserver/middleware"
	"github.com/hyeonmink/challenges-hyeonmink/apiserver/models/events"
	"github.com/hyeonmink/challenges-hyeonmink/apiserver/models/messages"
	"github.com/hyeonmink/challenges-hyeonmink/apiserver/models/users"
	"github.com/hyeonmink/challenges-hyeonmink/apiserver/sessions"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/redis.v5"
)

const defaultPort = "443"

const (
	apiRoot         = "/v1/"
	apiSummary      = apiRoot + "summary"
	apiUsers        = apiRoot + "users"
	apiSessions     = apiRoot + "sessions"
	apiSessionsMine = apiSessions + "/mine"
	apiUserMe       = apiUsers + "/me"
	apiChannel      = apiRoot + "channels"
	apiMessages     = apiRoot + "messages"
	apiWebSocket    = apiRoot + "websocket"
	apiChatbot      = apiRoot + "bot"
)

const (
	dbName      = "info344"
	colName     = "users"
	chanColName = "channels"
	msgColName  = "messages"
)

//main is the main entry point for this program
func main() {
	//read and use the following environment variables
	//when initializing and starting your web server
	// PORT - port number to listen on for HTTP requests (if not set, use defaultPort)
	// HOST - host address to respond to (if not set, leave empty, which means any host)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = defaultPort
	}
	host := os.Getenv("HOST")
	addr := host + ":" + port

	certPath := os.Getenv("TLSCERT")
	keyPath := os.Getenv("TLSKEY")
	sessKey := os.Getenv("SESSIONKEY")
	redisAddr := os.Getenv("REDISADDR")
	dbAddr := os.Getenv("DBADDR")
	svcAddr := os.Getenv("SVCADDR")

	if len(redisAddr) == 0 {
		redisAddr = "127.0.0.1:6379"
	}
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	if len(dbAddr) == 0 {
		dbAddr = "127.0.0.1:27017"
	}

	mgo, err := mgo.Dial(dbAddr)
	if err != nil {
		fmt.Printf("error dialing into Mongo: %v", err)
	}
	defer mgo.Close()

	rStore := sessions.NewRedisStore(client, sessions.DefaultSessionDuration)

	mongoStore := &users.MongoStore{
		Session:        mgo,
		DatabaseName:   dbName,
		CollectionName: colName,
	}

	msgStore := &messages.MongoStore{
		Session:               mgo,
		DatabaseName:          dbName,
		ChannelCollectionName: chanColName,
		MessageCollectionName: msgColName,
	}

	notifier := events.NewNotifier(300)
	go notifier.Start()

	ctx := handlers.Context{
		SessionKey:   sessKey,
		SessionStore: rStore,
		UserStore:    mongoStore,
		MessageStore: msgStore,
		Notifier:     notifier,
		SvcAddr:      svcAddr,
	}

	mux := http.NewServeMux()
	mux.HandleFunc(apiUsers, ctx.UsersHandler)
	mux.HandleFunc(apiSessions, ctx.SessionsHandler)
	mux.HandleFunc(apiSessionsMine, ctx.SessionsMineHandler)
	mux.HandleFunc(apiUserMe, ctx.UsersMeHanlder)
	mux.HandleFunc(apiSummary, handlers.SummaryHandler)
	mux.HandleFunc(apiChannel, ctx.ChannelsHandler)
	mux.HandleFunc(apiChannel+"/", ctx.SpecificChannelHandler)
	mux.HandleFunc(apiMessages, ctx.MessagesHandler)
	mux.HandleFunc(apiMessages+"/", ctx.SpecificMessageHandler)
	mux.HandleFunc(apiWebSocket, ctx.WebSocketUpgradeHandler)
	mux.HandleFunc(apiChatbot, ctx.ChatbotHandler)

	http.Handle(apiRoot, middleware.Adapt(mux,
		middleware.CORS("", "", "", ""),
	))

	//add your handlers.SummaryHandler function as a handler
	//for the apiSummary route
	//HINT: https://golang.org/pkg/net/http/#HandleFunc
	//http.HandleFunc(apiSummary, handlers.SummaryHandler)

	//start your web server and use log.Fatal() to log
	//any errors that occur if the server can't start
	//HINT: https://golang.org/pkg/net/http/#ListenAndServe
	fmt.Printf("server is listening at %s...\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, certPath, keyPath, nil))
}
