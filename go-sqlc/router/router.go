package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/handler"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/service"
	custom_middleware "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/pkg/middleware"
)

type Server struct {
	echo       *echo.Echo
	handler    *handler.Handler
	closing    atomic.Bool   // track server status
	shutdownCh chan struct{} // notify close channel
	tracker    *service.TaskTracker
}

func NewServer(services *service.Services) *Server {
	e := echo.New()
	// Create tracker
	tracker := service.NewTaskTracker()
	// init handler
	h := handler.NewHandler(services, tracker)

	server := &Server{
		echo:       e,
		handler:    h,
		shutdownCh: make(chan struct{}),
		tracker:    tracker,
	}

	// Global
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	stateChecker := custom_middleware.NewStateChecker(server)
	e.Use(stateChecker.CheckServerClosing)

	// register routes
	server.setupRoutes()

	return server
}

func (s *Server) setupRoutes() {
	e := s.echo

	// health check
	e.GET("/health", func(c echo.Context) error {
		status := "ok"
		if s.closing.Load() {
			status = "shutting_down"
		}
		return c.JSON(http.StatusOK, map[string]string{
			"status": status,
			"tasks":  fmt.Sprintf("%d", s.tracker.TaskCount()),
		})
	})

	// manually shutdown
	e.POST("/api/shutdown", s.handleShutdown)

	// API v1
	v1 := e.Group("/api/v1")

	// Authors routes
	authors := v1.Group("/authors")
	authors.POST("", s.handler.Author.Create)
	authors.GET("", s.handler.Author.List)
	authors.GET("/:id", s.handler.Author.Get)
	authors.PUT("/:id", s.handler.Author.Update)
	authors.DELETE("/:id", s.handler.Author.Delete)

	v1.POST("/long-task", s.handler.Task.LongRunningTask)
}

func (s *Server) Start() {
	fmt.Println("→ Initializing server...")

	// HTTP Server
	srv := &http.Server{
		Addr:         ":8888",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start
	go func() {
		fmt.Println("→ Server is starting on port 8888")

		if err := s.echo.StartServer(srv); err != nil {
			if err == http.ErrServerClosed {
				fmt.Println("← Server closed under request")
			} else {
				s.echo.Logger.Fatal("shutting down the server")
			}
		}
	}()

	// wait signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		fmt.Printf("→ System received shutdown signal: %v\n", sig)
	case <-s.shutdownCh:
		fmt.Println("→ Received shutdown request through API")
	}

	s.performGracefulShutdown()
}

// gracefully shutdown
func (s *Server) performGracefulShutdown() {
	fmt.Println("→ Starting graceful shutdown...")

	// Step. 1 close HTTP server
	if err := s.shutdownHTTPServer(); err != nil {
		fmt.Printf("← Warning: HTTP server shutdown error: %v\n", err)
	}

	// Step. 2
	s.waitForTasksCompletion()

	fmt.Println("← Server gracefully stopped")
}

func (s *Server) shutdownHTTPServer() error {
	fmt.Println("→ Phase 1: Closing HTTP connections...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return s.echo.Shutdown(ctx)
}

func (s *Server) waitForTasksCompletion() {
	fmt.Println("→ Phase 2: Waiting for running tasks to complete...")

	attemptCount := 1
	for {
		// create this countdown context
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

		// monitor this countdown
		go s.monitorTaskProgress(ctx, attemptCount)

		if completed := s.waitForTasksOrTimeout(ctx); completed {
			cancel()
			return
		}

		// next countdown
		cancel()
		attemptCount++
	}
}

func (s *Server) monitorTaskProgress(ctx context.Context, attempt int) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	timeLeft := 15
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			timeLeft--
			if timeLeft > 0 {
				remainingTasks := s.tracker.TaskCount()
				fmt.Printf("→ Attempt %d: %d seconds left (remaining tasks: %d)\n",
					attempt, timeLeft, remainingTasks)
			}
		}
	}
}

func (s *Server) waitForTasksOrTimeout(ctx context.Context) bool {
	// create task complete done chan
	taskDone := make(chan struct{})
	go func() {
		s.tracker.WaitForTasks()
		close(taskDone)
	}()

	select {
	case <-taskDone:
		fmt.Println("← All tasks completed successfully")
		return true
	case <-ctx.Done():
		remainingTasks := s.tracker.TaskCount()
		if remainingTasks > 0 {
			fmt.Printf("→ Tasks still running (remaining: %d). Starting next attempt...\n",
				remainingTasks)
			return false
		}
		return true
	}
}

// close request handler
func (s *Server) handleShutdown(c echo.Context) error {
	// check server closed or not
	if s.closing.Load() {
		return c.JSON(http.StatusConflict, map[string]string{
			"message": "Server is already shutting down",
		})
	}

	// set close true
	s.closing.Store(true)

	// in new goroutine exec close，return signal
	go func() {
		time.Sleep(500 * time.Millisecond)
		close(s.shutdownCh)
	}()

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Server shutdown initiated",
		"note":    "No new requests will be accepted",
	})
}

func (s *Server) IsClosing() bool {
	return s.closing.Load()
}

func (s *Server) GetTaskCount() int32 {
	return s.tracker.TaskCount()
}
