package routes

import (
	"github.com/gorilla/mux"
	controllers "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/controller"
)

func registerMovieRoutes(r *mux.Router, movieController *controllers.MovieController) {
	r.HandleFunc("", movieController.GetMovies).Methods("GET")
	r.HandleFunc("/{id}", movieController.GetMovie).Methods("GET")
	r.HandleFunc("", movieController.CreateMovie).Methods("POST")
	r.HandleFunc("/{id}", movieController.UpdateMovie).Methods("PUT")
	r.HandleFunc("/{id}", movieController.DeleteMovie).Methods("DELETE")
}
