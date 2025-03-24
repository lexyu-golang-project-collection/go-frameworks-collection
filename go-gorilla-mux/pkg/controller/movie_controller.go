package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	models "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/model"
	services "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/service"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/utils"
)

// MovieController 處理電影相關的 HTTP 請求
type MovieController struct {
	movieService *services.MovieService
}

// NewMovieController 創建新的電影控制器
func NewMovieController(movieService *services.MovieService) *MovieController {
	return &MovieController{
		movieService: movieService,
	}
}

// GetMovies 返回所有電影
// @Summary 獲取所有電影
// @Description 獲取系統中所有電影的列表
// @Tags movies
// @Accept json
// @Produce json
// @Success 200 {array} models.Movie
// @Router /movies [get]
func (c *MovieController) GetMovies(w http.ResponseWriter, r *http.Request) {
	movies := c.movieService.GetAllMovies()
	utils.RespondWithJSON(w, http.StatusOK, movies)
}

// GetMovie 返回特定電影
// @Summary 獲取特定電影
// @Description 根據 ID 獲取特定電影的詳細信息
// @Tags movies
// @Accept json
// @Produce json
// @Param id path string true "電影 ID"
// @Success 200 {object} models.Movie
// @Failure 404 {object} utils.ErrorResponse "電影未找到"
// @Router /movies/{id} [get]
func (c *MovieController) GetMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	movie, err := c.movieService.GetMovieByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "電影未找到")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, movie)
}

// CreateMovie 創建新電影
// @Summary 創建電影
// @Description 創建新的電影記錄
// @Tags movies
// @Accept json
// @Produce json
// @Param movie body models.Movie true "電影資料"
// @Success 201 {object} models.Movie
// @Failure 400 {object} utils.ErrorResponse "無效的請求體"
// @Router /movies [post]
func (c *MovieController) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var movie models.Movie

	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "無效的請求體")
		return
	}

	createdMovie := c.movieService.CreateMovie(movie)

	utils.RespondWithJSON(w, http.StatusCreated, createdMovie)
}

// UpdateMovie 更新電影
// @Summary 更新電影
// @Description 根據 ID 更新特定電影的信息
// @Tags movies
// @Accept json
// @Produce json
// @Param id path string true "電影 ID"
// @Param movie body models.Movie true "電影資料"
// @Success 200 {object} models.Movie
// @Failure 400 {object} utils.ErrorResponse "無效的請求體"
// @Failure 404 {object} utils.ErrorResponse "電影未找到"
// @Router /movies/{id} [put]
func (c *MovieController) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var movie models.Movie

	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "無效的請求體")
		return
	}

	updatedMovie, err := c.movieService.UpdateMovie(id, movie)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "電影未找到")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedMovie)
}

// DeleteMovie 刪除電影
// @Summary 刪除電影
// @Description 根據 ID 刪除特定電影
// @Tags movies
// @Accept json
// @Produce json
// @Param id path string true "電影 ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} utils.ErrorResponse "電影未找到"
// @Router /movies/{id} [delete]
func (c *MovieController) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	err := c.movieService.DeleteMovie(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "電影未找到")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "電影刪除成功"})
}
