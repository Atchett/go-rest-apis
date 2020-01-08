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
