package router

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	embedSQL "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/db/embed"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/db/sqlite"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/handler"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/service"
)

type Server struct {
	echo    *echo.Echo
	handler *handler.Handler
}

func NewServer(services *service.Services) *Server {
	e := echo.New()

	// Global
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// init handler
	h := handler.NewHandler(services)

	server := &Server{
		echo:    e,
		handler: h,
	}

	// register routes
	server.setupRoutes()

	return server
}

func (s *Server) setupRoutes() {
	// health check
	s.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// API v1
	v1 := s.echo.Group("/api/v1")

	// Authors routes
	authors := v1.Group("/authors")
	authors.POST("", s.handler.Author.Create)
	authors.GET("", s.handler.Author.List)
	authors.GET("/:id", s.handler.Author.Get)
	authors.PUT("/:id", s.handler.Author.Update)
	authors.DELETE("/:id", s.handler.Author.Delete)
	// Test route for in-memory SQLite
	v1.GET("/test", s.TestInMemory)
}

func (s *Server) TestInMemory(c echo.Context) error {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer db.Close()

	// exec create table
	if _, err := db.ExecContext(c.Request().Context(), embedSQL.GetSQLiteDDL()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	queries := sqlite.New(db)

	author, err := queries.CreateAuthor(c.Request().Context(), sqlite.CreateAuthorParams{
		Name: "Brian Kernighan",
		Bio:  sql.NullString{String: "Co-author of The C Programming Language", Valid: true},
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	authors, err := queries.ListAuthors(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"created_author": author,
		"all_authors":    authors,
	})
}

func (s *Server) Start() {
	// HTTP Server
	srv := &http.Server{
		Addr:         ":8888",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start
	go func() {
		if err := s.echo.StartServer(srv); err != nil && err != http.ErrServerClosed {
			s.echo.Logger.Fatal("shutting down the server")
		}
	}()

	// wait signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// gracefully shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.echo.Shutdown(ctx); err != nil {
		s.echo.Logger.Fatal(err)
	}
}
