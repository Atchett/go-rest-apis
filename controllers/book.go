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
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
		}

		// get row by ID from db
		row := db.QueryRow("select * from books where id=$1", id)
		err = row.Scan(&book.ID, &book.Title, &book.Author, &book.Year)

		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
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

		json.NewDecoder(r.Body).Decode(&book)

		err := db.QueryRow("insert into books (title, author, year) values($1, $2, $3) RETURNING id;", book.Title, book.Author, book.Year).Scan(&bookID)

		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
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

		result, err := db.Exec("update books set title=$1, author=$2, year=$3 where id=$4 RETURNING id", &book.Title, &book.Author, &book.Year, &book.ID)

		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
		}

		rowsUpdated, err := result.RowsAffected()

		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
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
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
		}

		result, err := db.Exec("delete from books where id = $1", id)
		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
		}

		rowsDeleted, err := result.RowsAffected()
		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
		}

		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, rowsDeleted)

	}
}
