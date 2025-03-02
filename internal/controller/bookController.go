package controller

import (
	"net/http"
	"strconv"

	"github.com/Redarcher9/Books-Management-System/internal/domain"
	"github.com/gin-gonic/gin"
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
// @Description Get all books for the provided limit and offset
// @Tags books
// @Accept json
// @Produce json
// @Param offset query int false "Offset for pagination" default(0) min(0)
// @Param limit query int false "Limit for pagination" default(10) min(1) max(100)
// @Success 200 {array} []domain.Book
// @Failure 400 {object} string
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
		g.JSON(http.StatusNotFound, gin.H{"error": "Books not found"})
		return
	}
	g.JSON(http.StatusOK, books)
}

// GetBookByID godoc
// @Summary Get a book by ID
// @Description Retrieve a book by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} domain.Book
// @Failure 400 {object} string "Invalid ID format"
// @Failure 404 {object} string "Book not found"
// @Router /books/{id} [get]
func (bc *BookController) GetBookByID(g *gin.Context) {
	// Get the 'ID' parameter
	IDParam := g.Param("id")

	// Convert 'ID' to Int
	id, err := strconv.Atoi(IDParam)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	book, err := bc.BookInteractor.GetBookByID(g, id)
	if err != nil {
		g.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	g.JSON(http.StatusOK, book)
}

// DeleteBookByID handles DELETE /books/:id
// @Summary Delete a book by ID
// @Description Delete a book by its ID
// @Tags books
// @Param id path int true "Book ID"
// @Success 204 "No Content"
// @Router /books/{id} [delete]
func (bc *BookController) DeleteBookByID(g *gin.Context) {
	// Get the 'ID' parameter
	IDParam := g.Param("id")

	// Convert 'ID' to Int
	id, err := strconv.Atoi(IDParam)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = bc.BookInteractor.DeleteBookByID(g, id)
	if err != nil {
		g.JSON(http.StatusNotFound, gin.H{"error": "Could not Delete Book"})
		return
	}
	g.Status(http.StatusNoContent)
}

// UpdateBookByID godoc
// @Summary Update a book by ID
// @Description Update an existing book by its ID with the provided data
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body domain.UpdateBookRequest true "Book data to update"
// @Success 200 {string} string "Book updated successfully"
// @Failure 400 {object} string "Invalid request data"
// @Failure 404 {string} string "Book not found"
// @Router /books/{id} [put]
func (bc *BookController) UpdateBookByID(g *gin.Context) {
	// Get the 'ID' parameter
	IDParam := g.Param("id")

	// Convert 'ID' to Int
	id, err := strconv.Atoi(IDParam)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req domain.UpdateBookRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// call Service for updating
	err = bc.BookInteractor.UpdateBookByID(g, id, req)
	if err != nil {
		g.Status(http.StatusNotFound)
	}
	g.Status(http.StatusOK)
}

// CreateBook godoc
// @Summary Create a new book
// @Description Create a new book with title, author, and year
// @Tags books
// @Accept json
// @Produce json
// @Param book body domain.Book true "Book data"
// @Success 201 {object} domain.Book
// @Failure 400 {object} string
// @Router /books [post]
func (bc *BookController) CreateBook(g *gin.Context) {
	var book domain.Book
	if err := g.BindJSON(&book); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Validate the input data using the Book's Validate method
	if err := book.Validate(); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := bc.BookInteractor.CreateBook(g, &book)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusCreated, book)
}
