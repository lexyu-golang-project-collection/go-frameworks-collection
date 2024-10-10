package middleware

import (
	"fmt"
	"net/http"

	"github.com/lexyu-golang-project-collection/go-design-patterns/combined/factory_with_strategy/internal/utils"
	logger "github.com/lexyu-golang-project-collection/go-design-patterns/combined/factory_with_strategy/pkg"
)

func GlobalErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Middleware GlobalErrorHandler------------------------")
		defer func() {
			if err := recover(); err != nil {
				logger.Error("occur error: %v", err)
				utils.JSONResponse(w, http.StatusInternalServerError, "internal error", err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
