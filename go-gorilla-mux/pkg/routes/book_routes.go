package routes

import (
	"github.com/gorilla/mux"
	controllers "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/controller"
)

func registerBookRoutes(r *mux.Router, bookController *controllers.BookController) {
	r.HandleFunc("", bookController.GetBooks).Methods("GET")
	r.HandleFunc("/{id}", bookController.GetBookByID).Methods("GET")
	r.HandleFunc("", bookController.CreateBook).Methods("POST")
	r.HandleFunc("/{id}", bookController.UpdateBook).Methods("PUT")
	r.HandleFunc("/{id}", bookController.DeleteBook).Methods("DELETE")
}
