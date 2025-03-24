package models

import (
	"sync"
)

// Director 代表電影導演
type Director struct {
	Firstname string `json:"firstname" example:"Steven"`
	Lastname  string `json:"lastname" example:"Spielberg"`
}

// Movie 代表電影資料
type Movie struct {
	ID       string    `json:"id" example:"1"`
	ISBN     string    `json:"isbn" example:"1234567890"`
	Title    string    `json:"title" example:"電影標題"`
	Director *Director `json:"director"`
}

var (
	movies     []Movie
	movieMutex sync.RWMutex
)

func GetAllMovies() []Movie {
	movieMutex.RLock()
	defer movieMutex.RUnlock()

	moviesCopy := make([]Movie, len(movies))
	copy(moviesCopy, movies)
	return moviesCopy
}

func InitMovies() {
	movieMutex.Lock()
	defer movieMutex.Unlock()

	if len(movies) == 0 {
		movies = append(movies, Movie{ID: "1", ISBN: "1234567", Title: "Movie One", Director: &Director{Firstname: "Lex", Lastname: "Yu"}})
		movies = append(movies, Movie{ID: "2", ISBN: "7654321", Title: "Movie Two", Director: &Director{Firstname: "Test", Lastname: "test"}})
	}
}
