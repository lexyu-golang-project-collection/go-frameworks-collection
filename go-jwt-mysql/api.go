package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	addr       string
	repository Repository
}

func NewApiServer(addr string, store Repository) *ApiServer {
	return &ApiServer{addr: addr, repository: store}
}

func (s *ApiServer) Serve() {
	router := mux.NewRouter()

	subRouter := router.PathPrefix("/api/v1").Subrouter()

	// Register services

	taskService := NewTasksService(s.repository)
	taskService.RegisterRoutes(subRouter)

	log.Println("Starting the API Server at", s.addr)

	log.Fatal(http.ListenAndServe(s.addr, subRouter))
}
