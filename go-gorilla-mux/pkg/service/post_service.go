package services

import (
	"context"

	models "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/model"
	base "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/repository/interfaces"
	repository "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/repository/postgres"
)

type PostService struct {
	repo      *repository.PostRepository
	txManager *base.TxManager
}

func NewPostService(repo *repository.PostRepository, txManager *base.TxManager) *PostService {
	return &PostService{
		repo:      repo,
		txManager: txManager,
	}
}

func (s *PostService) GetAllPosts() ([]models.Post, error) {
	ctx := context.Background()
	return s.repo.FindAll(ctx)
}

func (s *PostService) GetPostByID(id uint) (*models.Post, error) {
	ctx := context.Background()
	return s.repo.FindByID(ctx, id)
}

func (s *PostService) CreatePost(post *models.Post) (*models.Post, error) {
	ctx := context.Background()
	err := s.repo.Create(ctx, post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) UpdatePost(id uint, updateData *models.Post) (*models.Post, error) {
	ctx := context.Background()

	post, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if updateData.Title != "" {
		post.Title = updateData.Title
	}
	if updateData.Body != "" {
		post.Body = updateData.Body
	}

	err = s.repo.Update(ctx, post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) DeletePost(id uint) error {
	ctx := context.Background()
	return s.repo.Delete(ctx, id)
}

/*
func (s *PostService) PublishPostWithTags(ctx context.Context, post *model.Post, tags []string) error {
	return s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		post.Status = "published"
		post.PublishedAt = time.Now()

		if err := s.repo.Create(txCtx, post); err != nil {
			return err
		}

		return nil
	})
}
*/
