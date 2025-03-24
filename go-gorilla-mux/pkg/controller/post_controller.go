package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	models "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/model"
	services "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/service"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/utils"
)

type PostController struct {
	postService *services.PostService
}

func NewPostController(postService *services.PostService) *PostController {
	return &PostController{
		postService: postService,
	}
}

func (c *PostController) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := c.postService.GetAllPosts()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "無法獲取帖子")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, posts)
}

func (c *PostController) GetPostByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID := params["id"]

	ID, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "無效的帖子 ID")
		return
	}

	post, err := c.postService.GetPostByID(uint(ID))
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "帖子未找到")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, post)
}

func (c *PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "無效的請求體")
		return
	}

	createdPost, err := c.postService.CreatePost(&post)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "無法創建帖子")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, createdPost)
}

func (c *PostController) UpdatePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID := params["id"]

	ID, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "無效的帖子 ID")
		return
	}

	var updateData models.Post
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "無效的請求體")
		return
	}

	post, err := c.postService.UpdatePost(uint(ID), &updateData)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "帖子未找到或更新失敗")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, post)
}

func (c *PostController) DeletePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID := params["id"]

	ID, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "無效的帖子 ID")
		return
	}

	err = c.postService.DeletePost(uint(ID))
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "無法刪除帖子")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "帖子刪除成功"})
}
