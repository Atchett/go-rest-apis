package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/atchett/go-rest-apis/models"
	"github.com/atchett/go-rest-apis/driver"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

var books []models.Book
var db *sql.DB

func init() {
	gotenv.Load()
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db = driver.ConnectDB()

	r := mux.NewRouter()

	r.HandleFunc("/books", booksAll).Methods("GET")
	r.HandleFunc("/books/{id}", bookByID).Methods("GET")
	r.HandleFunc("/books", bookAdd).Methods("POST")
	r.HandleFunc("/books", bookUpdate).Methods("PUT")
	r.HandleFunc("/books/{id}", bookRemove).Methods("DELETE")

	fmt.Println("Server is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func booksAll(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	books = []models.Book{}

	rows, err := db.Query("SELECT * FROM books")
	logFatal(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		logFatal(err)
		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)

}

func bookByID(w http.ResponseWriter, r *http.Request) {

	var book models.Book

	// get Id from params
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	logFatal(err)

	// get row by ID from db
	row := db.QueryRow("select * from books where id=$1", id)
	err = row.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	logFatal(err)

	json.NewEncoder(w).Encode(book)

}

func bookAdd(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	var bookID int

	json.NewDecoder(r.Body).Decode(&book)

	err := db.QueryRow("insert into books (title, author, year) values($1, $2, $3) RETURNING id;", book.Title, book.Author, book.Year).Scan(&bookID)

	logFatal(err)

	json.NewEncoder(w).Encode(bookID)

}

func bookUpdate(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	json.NewDecoder(r.Body).Decode(&book)

	result, err := db.Exec("update books set title=$1, author=$2, year=$3 where id=$4 RETURNING id", &book.Title, &book.Author, &book.Year, &book.ID)
	logFatal(err)

	rowsUpdated, err := result.RowsAffected()
	logFatal(err)

	json.NewEncoder(w).Encode(rowsUpdated)

}

func bookRemove(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	logFatal(err)

	result, err := db.Exec("delete from books where id = $1", id)
	logFatal(err)

	rowsDeleted, err := result.RowsAffected()
	logFatal(err)

	json.NewEncoder(w).Encode(rowsDeleted)

}
