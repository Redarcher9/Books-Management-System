package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Redarcher9/Books-Management-System/internal/domain"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookController struct {
	BookInteractor BookService
}

func NewBookController(bookService BookService) *BookController {
	if bookService == nil {
		return nil
	}
	return &BookController{
		BookInteractor: bookService,
	}
}

// GetAllBooks godoc
// @Summary Get all books with pagination
// @Description Retrieve all books with pagination. If the provided offset or limit is less than 0, default values of limit = 10 and offset = 0 will be applied automatically.
// @Tags books
// @Accept json
// @Produce json
// @Param offset query int false "Offset for pagination" default(0) min(0)
// @Param limit query int false "Limit for pagination" default(10) min(1) max(100)
// @Success 200 {array} []domain.Book
// @Failure 500  "Internal Server Error"
// @Router /books [get]
func (bc *BookController) GetBooks(g *gin.Context) {
	// Get Query parameters with default values
	offset, err := strconv.Atoi(g.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	limit, err := strconv.Atoi(g.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	books, err := bc.BookInteractor.GetBooks(g, offset, limit)
	if err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: "Internal Server Error",
		})
		return
	}
	g.JSON(http.StatusOK, books)
}

// GetBookByID godoc
// @Summary Get a book by ID
// @Description Fetch detailed information about a book using its unique ID.
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} domain.Book
// @Failure 400 "Invalid ID format"
// @Failure 404 "Book not found"
// @Failure 500  "Internal Server Error"
// @Router /books/{id} [get]
func (bc *BookController) GetBookByID(g *gin.Context) {
	// Get the 'ID' parameter
	IDParam := g.Param("id")

	// Convert 'ID' to Int
	id, err := strconv.Atoi(IDParam)
	if err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Invalid ID format",
		})
		return
	}

	book, err := bc.BookInteractor.GetBookByID(g, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			g.JSON(http.StatusNotFound, domain.ErrorResponse{
				Message: fmt.Sprintf("Book for ID %d not found", id),
			})
			return
		}
		g.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: "Internal Server Error",
		})
		return
	}
	g.JSON(http.StatusOK, book)
}

// DeleteBookByID handles DELETE /books/:id
// @Summary Delete a book by ID
// @Description Delete a book by its ID
// @Tags books
// @Param id path int true "Book ID"
// @Success 200 "Book Deleted Successfully"
// @Failure 500 "Internal Server Error"
// @Router /books/{id} [delete]
func (bc *BookController) DeleteBookByID(g *gin.Context) {
	// Get the 'ID' parameter
	IDParam := g.Param("id")

	// Convert 'ID' to Int
	id, err := strconv.Atoi(IDParam)
	if err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Invalid ID format",
		})
		return
	}

	err = bc.BookInteractor.DeleteBookByID(g, id)
	if err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: "Internal Server Error",
		})
		return
	}
	g.Status(http.StatusOK)
}

// UpdateBookByID godoc
// @Summary Update a book by ID
// @Description Update an existing book by its ID with the provided data
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body domain.BookRequest true "Book data to update"
// @Success 200 "Book updated successfully"
// @Failure 400 "Validation Error"
// @Failure 404  "Book to update not found"
// @Failure 500  "Internal Server Error"
// @Router /books/{id} [put]
func (bc *BookController) UpdateBookByID(g *gin.Context) {
	// Get the 'ID' parameter
	IDParam := g.Param("id")

	// Convert 'ID' to Int
	id, err := strconv.Atoi(IDParam)
	if err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Invalid ID format",
		})
		return
	}

	var req domain.Book
	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: fmt.Sprintf("Invalid data: %s", err.Error()),
		})
		return
	}

	// Validate the input data using the Book's Validate method
	if err := req.Validate(); err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: fmt.Sprintf("Invalid data: %s", err.Error()),
		})
		return
	}

	// call Service for updating
	err = bc.BookInteractor.UpdateBookByID(g, id, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			g.JSON(http.StatusNotFound, domain.ErrorResponse{
				Message: fmt.Sprintf("Book for ID %d not found", id),
			})
			return
		}
		g.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: "Internal Server Error",
		})
		return
	}
	g.JSON(http.StatusOK, gin.H{"message": "book updated successfully"})
}

// CreateBook godoc
// @Summary Create a new book
// @Description Create a new book with title, author, and year
// @Tags books
// @Accept json
// @Produce json
// @Param book body domain.BookRequest true "Book data to create"
// @Success 201 "Book Created Successfully"
// @Failure 400 "Validation Error"
// @Failure 409  "Book with provided Title and Author already exists"
// @Failure 500  "Internal Server Error"
// @Router /books [post]
func (bc *BookController) CreateBook(g *gin.Context) {
	var req domain.Book
	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: fmt.Sprintf("Invalid data: %s", err.Error()),
		})
		return
	}

	// Validate the input data using the Book's Validate method
	if err := req.Validate(); err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: fmt.Sprintf("Invalid data: %s", err.Error()),
		})
		return
	}

	err := bc.BookInteractor.CreateBook(g, &req)
	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			g.JSON(http.StatusConflict, domain.ErrorResponse{
				Message: fmt.Sprintf("Book with Title %s and Author %s already exists", req.Title, req.Author),
			})
			return
		}
		g.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	g.JSON(http.StatusCreated, req)
}
