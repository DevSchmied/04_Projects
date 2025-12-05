package cache

import "bookmanager-go/internal/model"

type BookCache interface {
	GetBookList() ([]model.Book, error)
	SetBookList([]model.Book) error
	InvalidateBookList() error
}
