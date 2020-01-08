package bookrepository

import (
	"database/sql"
	"go-rest-apis/models"
)

// BookRepository - empty struct to add methods to
type BookRepository struct{}

// BooksAll - return the books to the user
func (b BookRepository) BooksAll(db *sql.DB, book models.Book, books []models.Book) ([]models.Book, error) {

	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		return []models.Book{}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		books = append(books, book)
	}

	if err != nil {
		return []models.Book{}, err
	}

	return books, nil
}

// BookByID - get the book by ID
func (b BookRepository) BookByID(db *sql.DB, book models.Book, id int) (models.Book, error) {

	// get row by ID from db
	row := db.QueryRow("select * from books where id=$1", id)
	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Year)

	return book, err
}

// BookAdd - add a book to the table
func (b BookRepository) BookAdd(db *sql.DB, book models.Book) (int, error) {

	err := db.QueryRow("insert into books (title, author, year) values($1, $2, $3) RETURNING id;", book.Title, book.Author, book.Year).Scan(&book.ID)

	if err != nil {
		return 0, err
	}
	return book.ID, nil
}

// BookUpdate - updates a book based on params passed in via body
func (b BookRepository) BookUpdate(db *sql.DB, book models.Book) (int64, error) {

	result, err := db.Exec("update books set title=$1, author=$2, year=$3 where id=$4 RETURNING id", &book.Title, &book.Author, &book.Year, &book.ID)
	if err != nil {
		return 0, err
	}

	rowsUpdated, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsUpdated, nil
}

// BookDelete - deletes a book based on the ID
func (b BookRepository) BookDelete(db *sql.DB, id int) (int64, error) {

	result, err := db.Exec("delete from books where id = $1", id)
	if err != nil {
		return 0, err
	}
	rowsDeleted, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsDeleted, nil

}
