package service

import (
	"context"
	"database/sql"

	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/db/sqlite"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/model"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/repository"
)

type AuthorService struct {
	dbManager *repository.DBManager
}

func NewAuthorService(dbManager *repository.DBManager) *AuthorService {
	return &AuthorService{
		dbManager: dbManager,
	}
}

type CreateAuthorRequest struct {
	Name string `json:"name" validate:"required"`
	Bio  string `json:"bio"`
}

func (s *AuthorService) CreateAuthor(ctx context.Context, req CreateAuthorRequest) (sqlite.Author, error) {
	params := sqlite.CreateAuthorParams{
		Name: req.Name,
		Bio:  sql.NullString{String: req.Bio, Valid: req.Bio != ""},
	}

	return s.dbManager.SQLiteQueries().CreateAuthor(ctx, params)
}

func (s *AuthorService) GetAuthor(ctx context.Context, id int64) (model.AuthorResponse, error) {
	author, err := s.dbManager.SQLiteQueries().GetAuthor(ctx, id)
	if err != nil {
		return model.AuthorResponse{}, err
	}
	return model.SQLiteAuthorToResponse(author), nil
}

func (s *AuthorService) ListAuthors(ctx context.Context) ([]model.AuthorResponse, error) {
	authors, err := s.dbManager.SQLiteQueries().ListAuthors(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]model.AuthorResponse, len(authors))
	for i, author := range authors {
		response[i] = model.SQLiteAuthorToResponse(author)
	}
	return response, nil
}

func (s *AuthorService) UpdateAuthor(ctx context.Context, id int64, req model.UpdateAuthorRequest) (model.AuthorResponse, error) {
	params := sqlite.UpdateAuthorParams{
		ID:   id,
		Name: req.Name,
		Bio:  sql.NullString{String: req.Bio, Valid: req.Bio != ""},
	}

	author, err := s.dbManager.SQLiteQueries().UpdateAuthor(ctx, params)
	if err != nil {
		return model.AuthorResponse{}, err
	}
	return model.SQLiteAuthorToResponse(author), nil
}

func (s *AuthorService) DeleteAuthor(ctx context.Context, id int64) error {
	return s.dbManager.SQLiteQueries().DeleteAuthor(ctx, id)
}

// func TestInMemory(ec echo.Context) error {

// 	author, err := s.CreateAuthor(ec.Request().Context(), sqlite.CreateAuthorParams{
// 		Name: "Brian Kernighan",
// 		Bio:  sql.NullString{String: "Co-author of The C Programming Language", Valid: true},
// 	})
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
// 	}

// 	authors, err := queries.ListAuthors(c.Request().Context())
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
// 	}

// 	return c.JSON(http.StatusOK, map[string]interface{}{
// 		"created_author": author,
// 		"all_authors":    authors,
// 	})
// }
