package router

import (
	"net/http"

	"github.com/lexyu-golang-project-collection/go-design-patterns/combined/factory_with_strategy/internal/utils"
)

type CustomServerMux struct {
	*http.ServeMux
	errorMessages map[string]map[int]string // pattern -> statusCode -> message
}

func NewCustomServeMux() *CustomServerMux {
	return &CustomServerMux{
		ServeMux:      http.NewServeMux(),
		errorMessages: make(map[string]map[int]string),
	}
}

func (csm *CustomServerMux) RegisterErrorMessage(pattern string, statusCode int, message string) {
	if _, exists := csm.errorMessages[pattern]; !exists {
		csm.errorMessages[pattern] = make(map[int]string)
	}
	csm.errorMessages[pattern][statusCode] = message
}

func (csm *CustomServerMux) getErrorMessage(pattern string, statusCode int) string {
	if messages, exists := csm.errorMessages[pattern]; exists {
		if msg, ok := messages[statusCode]; ok {
			return msg
		}
	}
	if msg, ok := csm.errorMessages["/"][statusCode]; ok {
		return msg
	}
	return http.StatusText(statusCode)
}

func (csm *CustomServerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, _ := csm.Handler(r)
	if handler == nil {
		// 如果沒有找到處理程序，這是一個 404 錯誤
		message := csm.getErrorMessage(r.URL.Path, http.StatusNotFound)
		utils.JSONResponse(w, http.StatusNotFound, message, nil)
		return
	}

	// // 使用 http.ResponseRecorder 來捕獲響應
	// rec := httptest.NewRecorder()
	// handler.ServeHTTP(rec, r)

	// // 檢查狀態碼
	// if rec.Code >= 400 {
	// 	message := csm.getErrorMessage(r.URL.Path, rec.Code)
	// 	utils.JSONResponse(w, rec.Code, message, nil)
	// } else {
	// 	// 如果不是錯誤，將記錄的響應複製到原始 ResponseWriter
	// 	for k, v := range rec.Header() {
	// 		w.Header()[k] = v
	// 	}
	// 	w.WriteHeader(rec.Code)
	// 	w.Write(rec.Body.Bytes())
	// }

	wrw := &wrappedResponseWriter{ResponseWriter: w, mux: csm}
	handler.ServeHTTP(wrw, r)
	wrw.Flush()
}

type wrappedResponseWriter struct {
	http.ResponseWriter
	mux     *CustomServerMux
	status  int
	written bool
	body    []byte
}

func (w *wrappedResponseWriter) WriteHeader(status int) {
	if !w.written {
		w.status = status
		w.written = true
	}
}

func (w *wrappedResponseWriter) Write(b []byte) (int, error) {
	if !w.written {
		w.WriteHeader(http.StatusOK)
	}
	w.body = append(w.body, b...)
	return len(b), nil
}

func (w *wrappedResponseWriter) Flush() {
	if w.status >= 400 {
		message := w.mux.getErrorMessage("/", w.status)
		utils.JSONResponse(w.ResponseWriter, w.status, message, nil)
	} else {
		w.ResponseWriter.WriteHeader(w.status)
		w.ResponseWriter.Write(w.body)
	}
	w.written = true
}
