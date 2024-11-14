package model

import "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/db/sqlite"

type AuthorResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Bio  string `json:"bio,omitempty"`
}

func SQLiteAuthorToResponse(a sqlite.Author) AuthorResponse {
	return AuthorResponse{
		ID:   a.ID,
		Name: a.Name,
		Bio:  a.Bio.String,
	}
}

type CreateAuthorRequest struct {
	Name string `json:"name" validate:"required"`
	Bio  string `json:"bio"`
}

type UpdateAuthorRequest struct {
	Name string `json:"name" validate:"required"`
	Bio  string `json:"bio"`
}
