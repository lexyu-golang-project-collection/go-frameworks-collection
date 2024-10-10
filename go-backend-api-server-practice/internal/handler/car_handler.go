package handler

import (
	"net/http"

	"github.com/lexyu-golang-project-collection/go-design-patterns/combined/factory_with_strategy/internal/service"
	"github.com/lexyu-golang-project-collection/go-design-patterns/combined/factory_with_strategy/internal/utils"
	logger "github.com/lexyu-golang-project-collection/go-design-patterns/combined/factory_with_strategy/pkg"
	types "github.com/lexyu-golang-project-collection/go-design-patterns/combined/factory_with_strategy/types"
)

type CarHandler struct {
	carService *service.CarService
}

func NewHandler(cs *service.CarService) *CarHandler {
	return &CarHandler{carService: cs}
}

func (handler *CarHandler) BuildCar(w http.ResponseWriter, r *http.Request) {
	carTypeStr := r.URL.Query().Get("type")

	err := utils.ValidatorCarType(carTypeStr)
	if err != nil {
		logger.Error("invalid car type: %s", carTypeStr)
		utils.JSONResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	carType := types.CarType(carTypeStr)
	car, err := handler.carService.BuildCar(carType)
	if err != nil {
		logger.Error("build car error: %v", err)
		return
	}

	logger.Info("build car success: %s %s", car.Brand, car.Model)
	utils.JSONResponse(w, http.StatusCreated, "build car success", car)

}

func (handler *CarHandler) BuildSportCar(w http.ResponseWriter, r *http.Request) {
	car, err := handler.carService.BuildCar(types.SPORT_CAR)
	if err != nil {
		logger.Error(string(types.BUILD_SPORT_CAR_ERR)+" %v", err)
		utils.JSONResponse(w, http.StatusInternalServerError, string(types.BUILD_SPORT_CAR_ERR), nil)
		return
	}

	logger.Info("build success: %s %s", car.Brand, car.Model)
	utils.JSONResponse(w, http.StatusCreated, "build success", car)
}

func (handler *CarHandler) BuildFamilyCar(w http.ResponseWriter, r *http.Request) {
	car, err := handler.carService.BuildCar(types.FAMILY_CAR)
	if err != nil {
		logger.Error(string(types.BUILD_FAMILY_CAR_ERR)+" %v", err)
		utils.JSONResponse(w, http.StatusInternalServerError, string(types.BUILD_FAMILY_CAR_ERR), nil)
		return
	}

	logger.Info("build success: %s %s", car.Brand, car.Model)
	utils.JSONResponse(w, http.StatusCreated, "build success", car)
}

// SimulateInternalError 模擬一個內部服務器錯誤
func (handler *CarHandler) SimulateInternalError(w http.ResponseWriter, r *http.Request) {
	// 故意觸發 panic 來模擬未處理的錯誤
	panic("Simulated internal server error")
}

// SimulateNonPanicInternalError 模擬一個非 panic 的內部服務器錯誤
func (handler *CarHandler) SimulateNonPanicInternalError(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Simulated internal server error", http.StatusInternalServerError)
}

// SimulateBadRequest 模擬一個無效請求錯誤
func (handler *CarHandler) SimulateBadRequest(w http.ResponseWriter, r *http.Request) {
	// 直接返回一個 400 錯誤
	http.Error(w, "Simulated bad request", http.StatusBadRequest)
}
