package router

import (
	"net/http"

	"github.com/lexyu-golang-project-collection/go-design-patterns/combined/factory_with_strategy/internal/handler"
	"github.com/lexyu-golang-project-collection/go-design-patterns/combined/factory_with_strategy/internal/middleware"
)

func RegisterCarRoutes(rr *RouterRegister, carHandler *handler.CarHandler) {

	rr.Register("GET", "/api/car", middleware.Logging(middleware.Auth(carHandler.BuildCar)),
		map[int]string{
			http.StatusNotFound:            "Car resource not found",
			http.StatusBadRequest:          "Invalid car request",
			http.StatusInternalServerError: "Error processing car request",
		})

	rr.Register("GET", "/api/car/sport", middleware.Logging(middleware.Auth(carHandler.BuildSportCar)), map[int]string{
		http.StatusNotFound:   "Sport car resource not found",
		http.StatusBadRequest: "Invalid sport car request",
	})

	rr.Register("GET", "/api/car/family", middleware.Logging(middleware.Auth(carHandler.BuildFamilyCar)), map[int]string{
		http.StatusNotFound:   "Family car resource not found",
		http.StatusBadRequest: "Invalid family car request",
	})

	rr.Register("GET", "/api/simulate/internal-error", middleware.Logging(middleware.Auth(carHandler.SimulateInternalError)),
		map[int]string{})

	rr.Register("GET", "/api/simulate/non-panic-internal-error", middleware.Logging(middleware.Auth(carHandler.SimulateNonPanicInternalError)),
		map[int]string{})

	rr.Register("GET", "/api/simulate/bad-request", middleware.Logging(middleware.Auth(carHandler.SimulateBadRequest)),
		map[int]string{})

}
