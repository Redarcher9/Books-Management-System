package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Redarcher9/Books-Management-System/internal/domain"
	"github.com/Redarcher9/Books-Management-System/internal/infrastructure/models/tables"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type Books struct {
	gormDB  *gorm.DB
	redisDB *redis.Client
}

const (
	bookListCacheKey    = "books:all"
	bookByIDCacheFormat = "books:%d"
	cacheTTL            = 10 * time.Minute // Cache expiration time
)

func NewBooksRepo(gormDB *gorm.DB, redisDB *redis.Client) *Books {
	return &Books{
		gormDB:  gormDB,
		redisDB: redisDB,
	}
}

func (b *Books) GetBooks(ctx context.Context, offset, limit int) ([]*domain.Book, error) {
	cacheKey := fmt.Sprintf("%s:%d:%d", bookListCacheKey, offset, limit)

	// Check if data is available in Redis
	if cachedData, err := b.redisDB.Get(cacheKey).Result(); err == nil {
		var domainBooks []*domain.Book
		if err := json.Unmarshal([]byte(cachedData), &domainBooks); err == nil {
			fmt.Println("Fetched from cache")
			return domainBooks, nil
		}
	}

	var books []*tables.Books
	result := b.gormDB.Model(&tables.Books{}).
		Limit(limit).
		Offset(offset).
		Find(&books)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// if book not found return nil
			return nil, nil
		}
		return nil, result.Error
	}
	if len(books) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	domainBooks := make([]*domain.Book, 0, len(books))
	for _, b := range books {
		domainBooks = append(domainBooks, b.ToDomain())
	}

	// Cache the result in Redis
	data, _ := json.Marshal(domainBooks)
	b.redisDB.Set(cacheKey, data, cacheTTL)

	return domainBooks, nil
}

func (b *Books) GetBookByID(ctx context.Context, ID int) (*domain.Book, error) {
	cacheKey := fmt.Sprintf(bookByIDCacheFormat, ID)

	// Check if data is available in Redis
	if cachedData, err := b.redisDB.Get(cacheKey).Result(); err == nil {
		var book domain.Book
		if err := json.Unmarshal([]byte(cachedData), &book); err == nil {
			fmt.Println("Fetched from cache")
			return &book, nil
		}
	}

	var book tables.Books // Note: Not a pointer here
	result := b.gormDB.
		Where("id = ?", ID).
		First(&book)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound // Return nil if not found
		}
		return nil, fmt.Errorf("failed to get book by ID: %w", result.Error)
	}

	// Cache the result in Redis
	data, _ := json.Marshal(book.ToDomain())
	b.redisDB.Set(cacheKey, data, cacheTTL)

	return book.ToDomain(), nil
}

func (b *Books) CreateBook(ctx context.Context, book *domain.Book) error {
	// Check if the book already exists (by Title and Author)
	var existingBook tables.Books
	if err := b.gormDB.Where("title = ? AND author = ?", book.Title, book.Author).First(&existingBook).Error; err == nil {
		return gorm.ErrDuplicatedKey
	}

	newBook := &tables.Books{
		Title:  book.Title,
		Author: book.Author,
		Year:   book.Year,
	}
	if err := b.gormDB.Create(newBook).Error; err != nil {
		return err
	}

	b.expireCache()
	return nil
}

func (b *Books) UpdateBookByID(ctx context.Context, ID int, book domain.Book) error {
	response := b.gormDB.Model(&tables.Books{}).Where("id = ?", ID).Updates(book)
	if response.Error != nil {
		return response.Error
	}

	if response.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	b.expireCache()
	b.redisDB.Del(fmt.Sprintf(bookByIDCacheFormat, ID))
	return nil
}

func (b *Books) DeleteBookByID(ctx context.Context, ID int) error {
	response := b.gormDB.Where("id = ?", ID).Delete(&tables.Books{})
	if response.Error != nil {
		return response.Error
	}
	b.expireCache()
	b.redisDB.Del(fmt.Sprintf(bookByIDCacheFormat, ID))
	return nil
}

// expireCache clears relevant cache entries in Redis
func (b *Books) expireCache() {
	// Delete all book list caches
	keys, err := b.redisDB.Keys("books:all:*").Result()
	if err == nil && len(keys) > 0 {
		b.redisDB.Del(keys...)
	}
}
