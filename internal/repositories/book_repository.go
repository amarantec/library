package repositories

import (
	"context"
	"errors"

	"github.com/amarantec/appserver/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Insert(ctx context.Context, book models.Book) (models.Book, error)
	FindAll(ctx context.Context) ([]models.Book, error)
	FindOneById(ctx context.Context, id int64) (models.Book, error)
	//Update(ctx context.Context, bool models.Book) error
	Delete(ctx context.Context, id int64) error
}

type RepositoryPostgres struct {
	Conn *pgxpool.Pool
}

func (r *RepositoryPostgres) Insert(ctx context.Context, book models.Book) (models.Book, error) {
	err := r.Conn.QueryRow(
		ctx,
		`INSERT INTO books (title, description, author, category) VALUES ($1, $2, $3, $4) RETURNIN id, title, description, author, category`,
		book.Title,
		book.Description,
		book.Author,
		book.Category).Scan(&book.ID, &book.Title, &book.Description, &book.Author, &book.Category)

	if err != nil {
		return models.Book{}, err
	}

	return book, nil
}

func (r *RepositoryPostgres) Delete(ctx context.Context, id int64) error {
	book, err := r.Conn.Exec(
		ctx,
		`DELETE FROM books WHERE id = $1`,
		id)

	if book.RowsAffected() == 0 {
		return errors.New("book not found")
	}

	return err
}

func (r *RepositoryPostgres) FindAll(ctx context.Context) ([]models.Book, error) {
	rows, err := r.Conn.Query(
		ctx,
		`SELECT id, title, description, author, category FROM books`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var books []models.Book

	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Description, &book.Author, &book.Category); err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (r *RepositoryPostgres) FindOneById(ctx context.Context, id int64) (models.Book, error) {

	var book = models.Book{ID: id}
	err := r.Conn.QueryRow(
		ctx,
		`SELECT title, description, author, category FROM books WHERE id=$1`,
		id).Scan(&book.Title, &book.Description, &book.Author, &book.Category)

	if err == pgx.ErrNoRows {
		return models.Book{}, errors.New("book not found")
	}

	if err != nil {
		return models.Book{}, err
	}

	return book, nil
}
