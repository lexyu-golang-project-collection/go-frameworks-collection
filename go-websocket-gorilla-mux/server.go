package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type webSocketHandler struct {
	upgrader websocket.Upgrader
}

func (wsh webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error %s when upgrading connection to websocket", err)
		return
	}
	defer func() {
		log.Println("closing connection")
		c.Close()
	}()

	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			return
		}
		if mt == websocket.BinaryMessage {
			err = c.WriteMessage(websocket.TextMessage, []byte("server doesn't support binary messages"))
			if err != nil {
				log.Printf("Error %s when sending message to client", err)
			}
			return
		}

		log.Printf("Receive message %s", string(msg))
		if strings.TrimSpace(string(msg)) != "open sesame" {
			err = c.WriteMessage(websocket.TextMessage, []byte("You need say the password"))
			if err != nil {
				log.Printf("Error %s when sending message to client", err)
				return
			}
			continue
		}

		log.Println("start responding to client...")
		i := 1
		for {
			response := fmt.Sprintf("Notification %d", i)
			err = c.WriteMessage(websocket.TextMessage, []byte(response))
			if err != nil {
				return
			}
			i += 1
			time.Sleep(3 * time.Second)
		}
	}
}

func main() {
	webSocketHandler := webSocketHandler{
		upgrader: websocket.Upgrader{},
	}

	http.Handle("/", webSocketHandler)
	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe("localhost:8888", nil))
}

/*
curl -i --header "Upgrade: websocket" \
--header "Connection: Upgrade" \
--header "Sec-WebSocket-Key: YSBzYW1wbGUgMTYgYnl0ZQ==" \
--header "Sec-Websocket-Version: 13" \
localhost:8888
*/
