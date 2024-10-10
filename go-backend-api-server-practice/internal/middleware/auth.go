package middleware

import (
	"fmt"
	"net/http"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Middleware Auth------------------------")
		fmt.Println("Step.2 - 執行身份驗證")
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "未提供授權令牌", http.StatusUnauthorized)
			return
		}
		fmt.Println("Token: ", token)
		fmt.Println("身份驗證成功")
		next(w, r)
	}
}
