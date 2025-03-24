package repo_mysql

import (
	"context"

	models "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/model"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/repository/base"
	"gorm.io/gorm"
)

type BookRepository struct {
	*base.BaseRepository[models.Book]
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{
		BaseRepository: base.NewBaseRepository[models.Book](db),
	}
}

func (r *BookRepository) FindByAuthor(ctx context.Context, author string) ([]models.Book, error) {
	var books []models.Book
	result := r.GetDB(ctx).Where("author = ?", author).Find(&books)
	return books, result.Error
}
