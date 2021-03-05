package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/hsson/gows"
)

type message struct {
	from        string
	data        string
	senderColor *color.Color
}

type client struct {
	name string
	ws   gows.Websocket
}

type hub struct {
	clients map[*client]struct{}

	register   chan *client
	unregister chan *client
	broadcast  chan message
}

func (h *hub) start() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = struct{}{}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				// don't send to self
				if client.name != message.from {
					err := client.ws.WriteText([]byte(fmt.Sprintf("%s: %s", message.senderColor.Sprint(message.from), message.data)))
					if err != nil {
						fmt.Printf("failed to send message to %s: %v", client.name, err)
					}
				}
			}
		}
	}
}

func newHub() *hub {
	return &hub{
		clients:    make(map[*client]struct{}),
		broadcast:  make(chan message),
		register:   make(chan *client),
		unregister: make(chan *client),
	}
}
