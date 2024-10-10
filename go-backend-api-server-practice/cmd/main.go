package main

import (
	"log"
	"net/http"

	"github.com/lexyu-golang-project-collection/go-design-patterns/combined/factory_with_strategy/internal/handler"
	"github.com/lexyu-golang-project-collection/go-design-patterns/combined/factory_with_strategy/internal/middleware"
	"github.com/lexyu-golang-project-collection/go-design-patterns/combined/factory_with_strategy/internal/service"
	logger "github.com/lexyu-golang-project-collection/go-design-patterns/combined/factory_with_strategy/pkg"
	"github.com/lexyu-golang-project-collection/go-design-patterns/combined/factory_with_strategy/router"
)

func main() {
	carService := service.NewCarService()
	carHandler := handler.NewHandler(carService)

	// mux := http.NewServeMux()
	mux := router.NewCustomServeMux()

	// mux.HandleFunc("GET /api/car", middleware.Logging(middleware.Auth(carHandler.BuildCar)))
	// mux.HandleFunc("GET /api/car/sport", middleware.Logging(middleware.Auth(carHandler.BuildSportCar)))
	// mux.HandleFunc("GET /api/car/family", middleware.Logging(middleware.Auth(carHandler.BuildFamilyCar)))

	register := router.NewRouteRegister(mux)
	router.RegisterCarRoutes(register, carHandler)

	// 註冊全局錯誤消息
	mux.RegisterErrorMessage("/", http.StatusNotFound, "Global - endpoint not found")
	mux.RegisterErrorMessage("/", http.StatusInternalServerError, "Global - internal server error (non-panic)")
	mux.RegisterErrorMessage("/", http.StatusBadRequest, "Global - invalid request")

	handler := middleware.GlobalErrorHandler(mux)

	logger.Info("Server running on http://localhost:8888")
	err := http.ListenAndServe(":8888", handler)
	if err != nil {
		log.Fatal(err)
	}
}
