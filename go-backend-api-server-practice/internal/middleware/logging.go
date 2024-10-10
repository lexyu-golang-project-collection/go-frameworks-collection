package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func Logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Middleware Logging------------------------")

		start := time.Now()
		fmt.Printf("開始處理請求: %s %s\n", r.Method, r.URL.Path)
		next(w, r)
		fmt.Printf("請求處理完成，耗時: %v\n", time.Since(start))
	}
}
