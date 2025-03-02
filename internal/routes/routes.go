package routes

import (
	docs "github.com/Redarcher9/Books-Management-System/docs"
	"github.com/Redarcher9/Books-Management-System/internal/infrastructure/kafka"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRoutes(gin *gin.Engine, gormDB *gorm.DB, kafka *kafka.KafkaProducer, redis *redis.Client) {
	// @BasePath /api/v1
	Router := gin.Group("/api/v1")
	NewHelloWorldRouter(Router)
	NewBookRouter(Router, gormDB, kafka, redis)
}

func SetupSwagger(gin *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/api/v1"
	gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
