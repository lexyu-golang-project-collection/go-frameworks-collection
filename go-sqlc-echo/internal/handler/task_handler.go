package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/service"
)

type TaskHandler struct {
	BaseHandler
}

func NewTaskHandler(tracker service.Tracker) *TaskHandler {
	return &TaskHandler{
		BaseHandler: BaseHandler{tracker: tracker},
	}
}

func (h *TaskHandler) LongRunningTask(c echo.Context) error {
	done := h.tracker.TrackTask()
	defer done()

	go func() {
		defer h.tracker.TrackTask()()

		fmt.Println("→ Long task started")
		time.Sleep(25 * time.Second)
		fmt.Println("→ Long task completed")
	}()

	return c.JSON(http.StatusAccepted, map[string]string{
		"message": "Task started",
		"status":  "processing",
	})
}
