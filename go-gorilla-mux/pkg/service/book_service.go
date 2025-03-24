package services

import (
	"context"

	models "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/model"
	base "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/repository/interfaces"
	repository "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/repository/mysql"
)

type BookService struct {
	repo      *repository.BookRepository
	txManager *base.TxManager
}

func NewBookService(repo *repository.BookRepository, txManager *base.TxManager) *BookService {
	return &BookService{
		repo:      repo,
		txManager: txManager,
	}
}

func (s *BookService) GetAllBooks() ([]models.Book, error) {
	ctx := context.Background()
	return s.repo.FindAll(ctx)
}

func (s *BookService) GetBookByID(id uint) (*models.Book, error) {
	ctx := context.Background()
	return s.repo.FindByID(ctx, id)
}

func (s *BookService) CreateBook(book *models.Book) (*models.Book, error) {
	ctx := context.Background()
	err := s.repo.Create(ctx, book)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (s *BookService) UpdateBook(id uint, updateData *models.Book) (*models.Book, error) {
	ctx := context.Background()

	book, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if updateData.Name != "" {
		book.Name = updateData.Name
	}
	if updateData.Author != "" {
		book.Author = updateData.Author
	}
	if updateData.Publication != "" {
		book.Publication = updateData.Publication
	}

	err = s.repo.Update(ctx, book)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (s *BookService) DeleteBook(id uint) error {
	ctx := context.Background()
	return s.repo.Delete(ctx, id)
}

/*
func (s *BookService) CreateBookWithAuthors(ctx context.Context, book *models.Book, authors []models.Author) error {
	return s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		if err := s.repo.Create(txCtx, book); err != nil {
			return err
		}

		return nil
	})
}
*/
