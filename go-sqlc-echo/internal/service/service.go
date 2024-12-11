package service

import (
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/repository"
)

type Services struct {
	Author *AuthorService
}

func NewServices(dbManager *repository.DBManager) *Services {
	return &Services{
		Author: NewAuthorService(dbManager),
	}
}
