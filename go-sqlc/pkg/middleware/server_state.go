package custom_middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ServerState interface {
	IsClosing() bool
	GetTaskCount() int32
}

type StateChecker struct {
	state ServerState
}

func NewStateChecker(state ServerState) *StateChecker {
	return &StateChecker{
		state: state,
	}
}

// check server status
func (sc *StateChecker) CheckServerClosing(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if sc.state.IsClosing() {
			if c.Path() == "/health" {
				return c.JSON(http.StatusOK, map[string]string{
					"status": "shutting_down",
					"tasks":  string(sc.state.GetTaskCount()),
				})
			}

			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"message": "Server is shutting down, please try again later",
				"status":  "unavailable",
			})
		}
		return next(c)
	}
}
