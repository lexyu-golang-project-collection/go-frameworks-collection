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

type BookController struct {
	bookService *services.BookService
}

func NewBookController(bookService *services.BookService) *BookController {
	return &BookController{
		bookService: bookService,
	}
}

func (c *BookController) GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := c.bookService.GetAllBooks()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "無法獲取書籍")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, books)
}

func (c *BookController) GetBookByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookID := params["id"]

	ID, err := strconv.ParseInt(bookID, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "無效的書籍 ID")
		return
	}

	book, err := c.bookService.GetBookByID(uint(ID))
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "書籍未找到")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, book)
}

func (c *BookController) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "無效的請求體")
		return
	}

	createdBook, err := c.bookService.CreateBook(&book)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "無法創建書籍")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, createdBook)
}

func (c *BookController) UpdateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookID := params["id"]

	ID, err := strconv.ParseInt(bookID, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "無效的書籍 ID")
		return
	}

	var updateData models.Book
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "無效的請求體")
		return
	}

	book, err := c.bookService.UpdateBook(uint(ID), &updateData)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "書籍未找到或更新失敗")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, book)
}

func (c *BookController) DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookID := params["id"]

	ID, err := strconv.ParseInt(bookID, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "無效的書籍 ID")
		return
	}

	err = c.bookService.DeleteBook(uint(ID))
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "無法刪除書籍")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "書籍刪除成功"})
}
