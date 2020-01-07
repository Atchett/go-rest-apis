package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Book - the model of the data
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   string `json:"year"`
}

var books []Book

func main() {
	r := mux.NewRouter()

	books = append(books,
		Book{ID: 1, Title: "Golang pointers", Author: "Mr. Golang", Year: "2010"},
		Book{ID: 2, Title: "Goroutines", Author: "Mr. Goroutine", Year: "2011"},
		Book{ID: 3, Title: "Golang routers", Author: "Mrs. Gorouter", Year: "2012"},
		Book{ID: 4, Title: "Golang concurrency", Author: "Mrs. Goconcurrency", Year: "2013"},
		Book{ID: 5, Title: "Goland good parts", Author: "Mrs. Good", Year: "2014"})

	r.HandleFunc("/books", booksAll).Methods("GET")
	r.HandleFunc("/books/{id}", bookByID).Methods("GET")
	r.HandleFunc("/books", bookAdd).Methods("POST")
	r.HandleFunc("/books", bookUpdate).Methods("PUT")
	r.HandleFunc("/books/{id}", bookRemove).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}

func booksAll(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(books)
}

func bookByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	// if err != nil {
	// 	http.Error(w, "Error converting param to int", http.StatusInternalServerError)
	// }

	for _, book := range books {
		if book.ID == id {
			json.NewEncoder(w).Encode(&book)
		}
	}
}

func bookAdd(w http.ResponseWriter, r *http.Request) {
	var book Book

	_ = json.NewDecoder(r.Body).Decode(&book)
	books = append(books, book)
	json.NewEncoder(w).Encode(books)
}

func bookUpdate(w http.ResponseWriter, r *http.Request) {
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	for i, item := range books {
		if item.ID == book.ID {
			books[i] = book
		}
	}

	json.NewEncoder(w).Encode(books)
}

func bookRemove(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for i, item := range books {
		if item.ID == id {
			books = append(books[:i], books[i+1:]...)
		}
	}
	json.NewEncoder(w).Encode(books)
}
