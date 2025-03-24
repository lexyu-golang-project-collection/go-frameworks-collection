package repo_mysql

import (
	"context"
	"errors"

	models "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/model"
)

type MovieRepository struct {
	movies []models.Movie
}

func NewMovieRepository() *MovieRepository {
	return &MovieRepository{
		movies: models.GetAllMovies(),
	}
}

func (r *MovieRepository) FindAll(ctx context.Context) ([]models.Movie, error) {
	moviesCopy := make([]models.Movie, len(r.movies))
	copy(moviesCopy, r.movies)
	return moviesCopy, nil
}

func (r *MovieRepository) FindByID(ctx context.Context, id string) (*models.Movie, error) {
	for _, movie := range r.movies {
		if movie.ID == id {
			return &movie, nil
		}
	}
	return nil, errors.New("電影未找到")
}

func (r *MovieRepository) Create(ctx context.Context, movie *models.Movie) error {
	r.movies = append(r.movies, *movie)
	return nil
}

func (r *MovieRepository) Update(ctx context.Context, id string, movie *models.Movie) error {
	for i, m := range r.movies {
		if m.ID == id {
			movie.ID = id
			r.movies[i] = *movie
			return nil
		}
	}
	return errors.New("電影未找到")
}

func (r *MovieRepository) Delete(ctx context.Context, id string) error {
	for i, movie := range r.movies {
		if movie.ID == id {
			r.movies = append(r.movies[:i], r.movies[i+1:]...)
			return nil
		}
	}
	return errors.New("電影未找到")
}
