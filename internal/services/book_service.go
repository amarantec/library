package services

import (
	"context"
	"errors"

	"github.com/amarantec/appserver/internal/models"
	"github.com/amarantec/appserver/internal/repositories"
)

var ErrBookNotFound = errors.New("book not found")
var ErrBookTitleEmpty = errors.New("book title is empty")
var ErrBookDescriptionEmpty = errors.New("book description is empty")
var ErrBookAuthorEmpty = errors.New("book author is empty")
var ErrBookCategoryEmpty = errors.New("book category is empty")

type Service struct {
	Repository repositories.Repository
}

func (s Service) Create(ctx context.Context, book models.Book) (models.Book, error) {
	if book.Title == "" {
		return models.Book{}, ErrBookTitleEmpty
	}
	if book.Description == "" {
		return models.Book{}, ErrBookDescriptionEmpty
	}
	if book.Author == nil {
		return models.Book{}, ErrBookAuthorEmpty
	}
	if book.Category == nil {
		return models.Book{}, ErrBookCategoryEmpty
	}

	return s.Repository.Insert(ctx, book)
}

func (s Service) Delete(ctx context.Context, id int64) error {
	return s.Repository.Delete(ctx, id)
}

func (s Service) FindOneById(ctx context.Context, id int64) (models.Book, error) {
	return s.Repository.FindOneById(ctx, id)
}

func (s Service) FindAll(ctx context.Context) ([]models.Book, error) {
	books, err := s.Repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return books, nil
}
