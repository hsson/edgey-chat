package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var endpoint = flag.String("endpoint", "ws://localhost:9999/", "server endpoint")
var port = flag.Int("port", 9999, "server port")

func main() {
	flag.Parse()
	fmt.Print("enter name: ")
	var name string
	fmt.Scanln(&name)
	name = strings.TrimSpace(name)

	url := fmt.Sprintf("%s:%d?name=%s", *endpoint, *port, name)
	ws, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		fmt.Printf("could not connect to server: %v\n", err)
		os.Exit(1)
		return
	}
	defer ws.Close()

	if resp.StatusCode >= 400 {
		fmt.Printf("got bad response from server: %d\n", resp.StatusCode)
		os.Exit(1)
		return
	}

	// start reading from server in a new goroutine
	go func() {
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				os.Exit(1)
				return
			}
			fmt.Printf("\r%s\n> ", message)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		<-interrupt
		fmt.Println("closing connection...")
		ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		time.Sleep(time.Second)
		os.Exit(1)
	}()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		msg, _ := reader.ReadString('\n')
		msg = strings.TrimSpace(msg)
		if msg != "" {
			ws.WriteMessage(websocket.TextMessage, []byte(msg))
		}
	}
}
