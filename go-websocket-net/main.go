package main

import (
	"fmt"
	"log"
	"net/http"
)

func Echo() {

}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Setting up the server!")
	})

	log.Fatal(http.ListenAndServe(":8888", nil))
}
