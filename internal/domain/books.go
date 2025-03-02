package domain

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type Book struct {
	ID     int    `json:"id" example:"1" validate:"omitempty"`
	Title  string `json:"title" validate:"required,max=255"`
	Author string `json:"author" validate:"required,max=255"`
	Year   int    `json:"year" example:"1957" validate:"required,validYear"`
}

func validYear(f1 validator.FieldLevel) bool {
	year := f1.Field().Int()
	return year >= 1450 && year <= 2025
}

// RegisterValidators registers custom validators
func RegisterValidators(v *validator.Validate) {
	v.RegisterValidation("validYear", validYear)
}

// validate is a self-contained function to validate the Book struct
func (b *Book) Validate() error {
	validate := validator.New()

	// Register custom validation for 'Year'
	validate.RegisterValidation("validYear", func(fl validator.FieldLevel) bool {
		year := fl.Field().Int()
		return year >= 1450 && year <= 2025
	})

	// Validate the struct
	if err := validate.Struct(b); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return errors.New("validation failed for field: " + e.Field())
		}
	}
	return nil
}

type BookRequest struct {
	Title  string `json:"title" validate:"required,max=255"`
	Author string `json:"author" validate:"required,max=255"`
	Year   int    `json:"year" example:"1957" validate:"required,validYear"`
}
