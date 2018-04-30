package events

import (
	"fmt"
	"sync"

	"encoding/json"

	"github.com/gorilla/websocket"
)

//Notifier struct
type Notifier struct {
	eventq  chan interface{}
	clients map[*websocket.Conn]bool
	mu      sync.RWMutex
}

//NewNotifier creates new notifier
func NewNotifier(size int) *Notifier {
	if size == 0 {
		size = 300
	}
	notifier := &Notifier{
		eventq:  make(chan interface{}, size),
		clients: make(map[*websocket.Conn]bool),
		mu:      sync.RWMutex{},
	}
	return notifier
}

//Start Begins a loot that checks for new events
func (n *Notifier) Start() {
	for {
		select {
		case e := <-n.eventq:
			n.broadcast(e)
		}
	}
}

func (n *Notifier) broadcast(event interface{}) {
	n.mu.Lock()
	defer n.mu.Unlock()

	eJSON, err := json.Marshal(event)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	pm, err := websocket.NewPreparedMessage(websocket.TextMessage, eJSON)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	for c := range n.clients {
		err := c.WritePreparedMessage(pm)
		if err != nil {
			delete(n.clients, c)
			c.Close()
		}
	}
}

//AddEvent add a new event into event queue
func (n *Notifier) AddEvent(event interface{}) {
	n.eventq <- event
}

//AddClient add a new client into client map
func (n *Notifier) AddClient(client *websocket.Conn) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.clients[client] = true
	go n.readPump(client)
}

func (n *Notifier) readPump(client *websocket.Conn) {
	for {
		if _, _, err := client.NextReader(); err != nil {
			client.Close()
			break
		}
	}
}
