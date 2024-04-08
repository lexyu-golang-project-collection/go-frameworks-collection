package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// User Struct
type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Handler Interface
type welcome string

func (wc welcome) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcom to our Server!")
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Login!")
}

func getJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	switch r.Method {
	case "GET":
		w.Write([]byte(`"message":"GET endpoint is called"`))
	case "POST":
		w.Write([]byte(`"message":"POST endpoint is called"`))
	}
}

func checkUser(w http.ResponseWriter, r *http.Request) {
	var user User
	dbPassword := "secretes"

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal("error decoding into struct")
	}

	if user.Password != dbPassword {
		fmt.Println("Password is wrong!")
		return
	}

	fmt.Println("Success Log In!")
	fmt.Fprintf(w, "Response %v", user)
}

func main() {
	// Router
	router := http.NewServeMux()

	// Handler
	var wc welcome
	router.Handle("/", wc)

	// Handler Funcs
	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Logout Page!")
	})

	router.HandleFunc("/login", login)

	router.HandleFunc("/json", getJson)

	router.HandleFunc("/user", checkUser)

	// Server
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Run Server
	server.ListenAndServe()
}
