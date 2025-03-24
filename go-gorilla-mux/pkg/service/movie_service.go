package services

import (
	"context"
	"math/rand"
	"strconv"
	"sync"

	models "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/model"
	repository "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/repository/mysql"
)

type MovieService struct {
	repo       *repository.MovieRepository
	movieMutex sync.RWMutex
}

func NewMovieService(repo *repository.MovieRepository) *MovieService {
	return &MovieService{
		repo: repo,
	}
}

func (s *MovieService) GetAllMovies() []models.Movie {
	s.movieMutex.RLock()
	defer s.movieMutex.RUnlock()

	ctx := context.Background()
	movies, _ := s.repo.FindAll(ctx)
	return movies
}

func (s *MovieService) GetMovieByID(id string) (models.Movie, error) {
	s.movieMutex.RLock()
	defer s.movieMutex.RUnlock()

	ctx := context.Background()
	movie, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return models.Movie{}, err
	}
	return *movie, nil
}

func (s *MovieService) CreateMovie(movie models.Movie) models.Movie {
	s.movieMutex.Lock()
	defer s.movieMutex.Unlock()

	movie.ID = strconv.Itoa(rand.Intn(100000000000))

	ctx := context.Background()
	s.repo.Create(ctx, &movie)
	return movie
}

func (s *MovieService) UpdateMovie(id string, updatedMovie models.Movie) (models.Movie, error) {
	s.movieMutex.Lock()
	defer s.movieMutex.Unlock()

	ctx := context.Background()
	err := s.repo.Update(ctx, id, &updatedMovie)
	if err != nil {
		return models.Movie{}, err
	}

	updatedMovie.ID = id
	return updatedMovie, nil
}

func (s *MovieService) DeleteMovie(id string) error {
	s.movieMutex.Lock()
	defer s.movieMutex.Unlock()

	ctx := context.Background()
	return s.repo.Delete(ctx, id)
}
