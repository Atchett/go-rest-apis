package controllers

import (
	"database/sql"
	"encoding/json"
	"go-rest-apis/models"
	bookrepository "go-rest-apis/repository/book"
	"go-rest-apis/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Controller - Controller struct
type Controller struct {
}

var books []models.Book

// BooksAll - gets the books from the db
func (c Controller) BooksAll(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book models.Book
		var error models.Error
		books := []models.Book{}
		bookRepo := bookrepository.BookRepository{}

		books, err := bookRepo.BooksAll(db, book, books)

		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
		}
		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, books)

	}
}

// BookByID - gets the book by its ID
func (c Controller) BookByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var book models.Book
		var error models.Error

		// get Id from params
		params := mux.Vars(r)
		id, _ := strconv.Atoi(params["id"])

		bookRepo := bookrepository.BookRepository{}
		book, err := bookRepo.BookByID(db, book, id)

		if err != nil {
			if err == sql.ErrNoRows {
				error.Message = "Not found"
				utils.SendError(w, http.StatusNotFound, error)
				return
			}
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, book)

	}
}

// BookAdd - add a book
func (c Controller) BookAdd(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book models.Book
		var bookID int
		var error models.Error

		// decode body into the book
		json.NewDecoder(r.Body).Decode(&book)

		// check book has the right values
		if book.Title == "" || book.Author == "" || book.Year == "" {
			error.Message = "Enter missing fields"
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		bookRepo := bookrepository.BookRepository{}
		bookID, err := bookRepo.BookAdd(db, book)

		if err != nil {
			error.Message = "Server Error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, bookID)

	}
}

// BookUpdate - update a book
func (c Controller) BookUpdate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book models.Book
		var error models.Error

		json.NewDecoder(r.Body).Decode(&book)

		// check book has the right values
		if book.ID == 0 || book.Title == "" || book.Author == "" || book.Year == "" {
			error.Message = "All fields are required"
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		bookRepo := bookrepository.BookRepository{}
		rowsUpdated, err := bookRepo.BookUpdate(db, book)

		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, rowsUpdated)

	}
}

// BookRemove - update a book
func (c Controller) BookRemove(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var error models.Error
		params := mux.Vars(r)
		id, _ := strconv.Atoi(params["id"])

		bookRepo := bookrepository.BookRepository{}

		rowsDeleted, err := bookRepo.BookDelete(db, id)
		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		if rowsDeleted == 0 {
			error.Message = "Not found"
			utils.SendError(w, http.StatusNotFound, error)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, rowsDeleted)

	}
}
