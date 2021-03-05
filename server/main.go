package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/hsson/gows"
)

func main() {
	log.Println("starting server")

	chatHub := newHub()
	go chatHub.start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = fmt.Sprintf("user-%d", rand.Int31n(1000))
		}

		ws, err := gows.Upgrade(w, r)
		if err != nil {
			log.Printf("internal server error: %v", err)
			fmt.Fprintf(w, "internal server error: %v", err)
			w.WriteHeader(500)
			return
		}

		go startChat(chatHub, name, ws)
	})

	http.ListenAndServe(":9999", nil)
}

func startChat(chatHub *hub, name string, ws gows.Websocket) {
	client := &client{
		name: name,
		ws:   ws,
	}
	senderColor := randomColor()
	ws.WriteText([]byte(fmt.Sprintf("connected as %s", name)))
	chatHub.register <- client
	log.Printf("%s connected", name)

	for {
		select {
		case <-ws.OnClose():
			// connection was closed
			log.Printf("%s disconnected", name)
			chatHub.unregister <- client
			return
		case msg := <-ws.Read():
			// got message from client
			log.Printf("got msg from %s", name)
			chatMessage := message{
				from:        name,
				data:        string(msg.Data),
				senderColor: senderColor,
			}
			chatHub.broadcast <- chatMessage
		}
	}
}
