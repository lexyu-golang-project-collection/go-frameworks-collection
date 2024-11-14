package handler

import "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/service"

type Handler struct {
	Author *AuthorHandler
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		Author: NewAuthorHandler(services.Author),
	}
}
