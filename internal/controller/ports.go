package controller

import (
	"context"

	"github.com/Redarcher9/Books-Management-System/internal/domain"
)

type (
	BookService interface {
		GetBooks(ctx context.Context, offset, limit int) ([]*domain.Book, error)
		GetBookByID(ctx context.Context, ID int) (*domain.Book, error)
		DeleteBookByID(ctx context.Context, ID int) error
		UpdateBookByID(ctx context.Context, ID int, book domain.UpdateBookRequest) error
		CreateBook(ctx context.Context, book *domain.Book) error
	}
)
