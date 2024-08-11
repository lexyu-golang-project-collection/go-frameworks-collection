package main

import (
	"log"
	"net/http"
)

func main() {
	setup()
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func setup() {
	wsManager := NewManager()
	http.Handle("/", http.FileServer(http.Dir("./fe")))
	http.HandleFunc("/ws", wsManager.serveWS)
}
