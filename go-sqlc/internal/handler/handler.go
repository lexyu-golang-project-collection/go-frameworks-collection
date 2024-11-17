package handler

import (
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/service"
)

type Handler struct {
	Author *AuthorHandler
	Task   *TaskHandler
}

func NewHandler(services *service.Services, tracker service.Tracker) *Handler {
	return &Handler{
		Author: NewAuthorHandler(services.Author),
		Task:   NewTaskHandler(tracker),
	}
}
