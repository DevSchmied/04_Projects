package cache

import (
	"bookmanager-go/internal/model"
	"context"
)

type BookCacher interface {
	GetBookList(ctx context.Context) ([]model.Book, error)
	SetBookList(ctx context.Context, books []model.Book) error
	InvalidateBookList(ctx context.Context) error
}
