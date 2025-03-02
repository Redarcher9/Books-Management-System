package routes

import (
	"github.com/Redarcher9/Books-Management-System/internal/controller"
	"github.com/Redarcher9/Books-Management-System/internal/infrastructure/kafka"
	"github.com/Redarcher9/Books-Management-System/internal/infrastructure/repository"
	"github.com/Redarcher9/Books-Management-System/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func NewBookRouter(group *gin.RouterGroup, db *gorm.DB, kafka *kafka.KafkaProducer, redis *redis.Client) {
	//Instantiate Repository, Service and Controller through dependency injection
	bookRepo := repository.NewBooksRepo(db, redis)
	bookService := service.NewBookInteractor(bookRepo, kafka)
	bookController := controller.NewBookController(bookService)

	//Initialise Routes
	group.GET("/books", bookController.GetBooks)
	group.GET("/books/:id", bookController.GetBookByID)
	group.DELETE("/books/:id", bookController.DeleteBookByID)
	group.PUT("/books/:id", bookController.UpdateBookByID)
	group.POST("/books", bookController.CreateBook)
}
