package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

var errNameRequired = errors.New("name is required")
var errProjectIDRequired = errors.New("project id is required")
var errUserIDRequired = errors.New("user id is required")

type TaskService struct {
	store Repository
}

func NewTasksService(s Repository) *TaskService {
	return &TaskService{store: s}
}

func (s *TaskService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/tasks", WithJWTAuth(s.handleCreateTask)).Methods("POST")
	r.HandleFunc("/tasks/{id}", s.handleGetTask).Methods("GET")
}

func (s *TaskService) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	defer r.Body.Close()

	var task *Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	if err := validateTaskPayload(task); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	t, err := s.store.CreateTask(task)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error createing task"})
		return
	}

	WriteJSON(w, http.StatusCreated, t)
}

func (s *TaskService) handleGetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	/*
		if id == "" {
			WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "id is required"})
			return
		}
	*/

	t, err := s.store.GetTask(id)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "task not found"})
		return
	}

	WriteJSON(w, http.StatusOK, t)
}

func validateTaskPayload(task *Task) error {
	if task.Name == "" {
		return errNameRequired
	}

	if task.ProjectID == 0 {
		return errProjectIDRequired
	}

	if task.AssignedToID == 0 {
		return errUserIDRequired
	}

	return nil
}
