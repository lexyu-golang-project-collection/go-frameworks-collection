package main

import (
	"fmt"
	// "slices"
	// "log/slog"
	// "math/rand/v2"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	// After 1.22 version
	/* HTTP Methods
	mux.HandleFunc("GET /demo", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`Hello World!`))
	})

	*/

	// /* Parameter, WildCards
	// Demo 1
	// rand.N[int64](int64(10))
	// slog.SetLogLoggerLevel()
	// slices.Concat[]()

	mux.HandleFunc("GET /demo/{name...}", func(w http.ResponseWriter, r *http.Request) {

		name := r.PathValue("name")

		println("{}", name)

		w.Write([]byte(fmt.Sprintf("Hello %s!", name)))
	})

	mux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the home page!"))
	})
	// */

	/* Parameter, WildCards
	// Demo 2
	mux.HandleFunc("GET /orders", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Get all orders"))
	})
	mux.HandleFunc("POST /orders", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Create an order"))
	})
	mux.HandleFunc("GET /orders/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "Get order with id: %s", id)
	})
	mux.HandleFunc("PUT /orders/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "Update order with id: %s", id)
	})
	*/

	// Before 1.22 version
	/* HTTP Methods
	mux.HandleFunc("/demo", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(""))
			return
		}

		w.Write([]byte(`Hello World!`))
	})
	*/

	/* Parameter
		mux.HandleFunc("/demo/", func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path

		parts := strings.Split(path, "/")

		if len(parts) < 3 {
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}

		name := parts[2]

		w.Write([]byte(fmt.Sprintf("Hello %s!", name)))
	})
	*/

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}

}
