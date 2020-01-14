package main

import (
	"database/sql"
	"fmt"
	"go-rest-apis/controllers"
	"go-rest-apis/driver"
	"go-rest-apis/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

var books []models.Book

var (
	db *sql.DB
	_  = gotenv.Load()
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db = driver.ConnectDB()
	controller := controllers.Controller{}

	r := mux.NewRouter()

	r.HandleFunc("/books", controller.BooksAll(db)).Methods("GET")
	r.HandleFunc("/books/{id}", controller.BookByID(db)).Methods("GET")
	r.HandleFunc("/books", controller.BookAdd(db)).Methods("POST")
	r.HandleFunc("/books", controller.BookUpdate(db)).Methods("PUT")
	r.HandleFunc("/books/{id}", controller.BookRemove(db)).Methods("DELETE")

	fmt.Println("Server is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
