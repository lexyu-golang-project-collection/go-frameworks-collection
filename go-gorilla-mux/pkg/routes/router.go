package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	controllers "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/controller"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(
	movieController *controllers.MovieController,
	bookController *controllers.BookController,
	postController *controllers.PostController,
) *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
	))

	r.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	api := r.PathPrefix("/api").Subrouter()

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	registerMovieRoutes(api.PathPrefix("/movies").Subrouter(), movieController)
	registerBookRoutes(api.PathPrefix("/books").Subrouter(), bookController)
	registerPostRoutes(api.PathPrefix("/posts").Subrouter(), postController)

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "資源未找到"}`))
	})

	return r
}
