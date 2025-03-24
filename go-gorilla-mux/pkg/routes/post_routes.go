package routes

import (
	"github.com/gorilla/mux"
	controllers "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/controller"
)

func registerPostRoutes(r *mux.Router, postController *controllers.PostController) {
	r.HandleFunc("", postController.GetPosts).Methods("GET")
	r.HandleFunc("/{id}", postController.GetPostByID).Methods("GET")
	r.HandleFunc("", postController.CreatePost).Methods("POST")
	r.HandleFunc("/{id}", postController.UpdatePost).Methods("PUT")
	r.HandleFunc("/{id}", postController.DeletePost).Methods("DELETE")
}
