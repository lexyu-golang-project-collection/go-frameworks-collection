package repo_postgres

import (
	"context"

	models "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/model"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/repository/base"
	"gorm.io/gorm"
)

type PostRepository struct {
	*base.BaseRepository[models.Post]
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{
		BaseRepository: base.NewBaseRepository[models.Post](db),
	}
}

func (r *PostRepository) FindByCategory(ctx context.Context, category string) ([]models.Post, error) {
	var posts []models.Post
	result := r.GetDB(ctx).Where("category = ?", category).Find(&posts)
	return posts, result.Error
}
