package service

import (
	"context"

	"github.com/Redarcher9/Books-Management-System/internal/domain"
)

type BookInteractor struct {
	Repo          BookRepo
	KafkaProducer KafkaProducer
}

// NewBookInteractor returns a valid book interactor
func NewBookInteractor(repo BookRepo, KafkaProducer KafkaProducer) *BookInteractor {
	if repo == nil {
		return nil
	}
	return &BookInteractor{
		Repo:          repo,
		KafkaProducer: KafkaProducer,
	}
}

func (c BookInteractor) GetBooks(ctx context.Context, offset, limit int) ([]*domain.Book, error) {
	return c.Repo.GetBooks(ctx, offset, limit)
}

func (c BookInteractor) GetBookByID(ctx context.Context, ID int) (*domain.Book, error) {
	return c.Repo.GetBookByID(ctx, ID)
}

func (c BookInteractor) DeleteBookByID(ctx context.Context, ID int) error {
	err := c.Repo.DeleteBookByID(ctx, ID)
	if err != nil {
		return err
	}
	message := map[string]interface{}{
		"event": "DELETE",
		"ID":    ID,
	}
	//Publish kafka message
	c.KafkaProducer.Publish(ctx, "book_events", message)
	return nil
}

func (c BookInteractor) UpdateBookByID(ctx context.Context, ID int, book domain.Book) error {
	err := c.Repo.UpdateBookByID(ctx, ID, book)
	if err != nil {
		return err
	}
	message := map[string]interface{}{
		"event":  "UPDATE",
		"ID":     ID,
		"TITLE":  book.Title,
		"AUTHOR": book.Author,
		"YEAR":   book.Year,
	}
	//Publish kafka message
	c.KafkaProducer.Publish(ctx, "book_events", message)
	return nil
}

func (c BookInteractor) CreateBook(ctx context.Context, book *domain.Book) error {
	err := c.Repo.CreateBook(ctx, book)
	if err != nil {
		return err
	}
	message := map[string]interface{}{
		"event":  "CREATE",
		"TITLE":  book.Title,
		"AUTHOR": book.Author,
		"YEAR":   book.Year,
	}
	//Publish kafka message
	c.KafkaProducer.Publish(ctx, "book_events", message)
	return nil
}
