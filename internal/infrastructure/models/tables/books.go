package tables

import "github.com/Redarcher9/Books-Management-System/internal/domain"

type Books struct {
	ID     int    `gorm:"column:id;primaryKey;autoIncrement"`
	Title  string `gorm:"column:title"`
	Author string `gorm:"column:author"`
	Year   int    `gorm:"column:year"`
}

func (b Books) TableName() string {
	return "books"
}

func (b Books) ToDomain() *domain.Book {
	res := &domain.Book{
		ID:     b.ID,
		Author: b.Author,
		Title:  b.Title,
		Year:   b.Year,
	}
	return res
}
