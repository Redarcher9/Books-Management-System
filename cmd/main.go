package main

import (
	"fmt"
	"time"

	"github.com/Redarcher9/Books-Management-System/config"
	"github.com/Redarcher9/Books-Management-System/internal/infrastructure/kafka"
	"github.com/Redarcher9/Books-Management-System/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var envConfig = config.Init()

func main() {
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Update with specific origins in production
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	//Instantiate Database,Kafka and Redis Server
	dbInstance := setUpDatabase()
	kafkaInstance := setUpKafkaProducer()
	redisInstance := setUpRedis()

	//Setup Routes and Swagger URLs
	routes.SetupRoutes(r, dbInstance, kafkaInstance, redisInstance)
	routes.SetupSwagger(r)

	//Run the Gin server on specified port
	port := fmt.Sprintf(":%s", envConfig.APIPort)
	r.Run(port)
}

// Sets up Database and Returns Database Instance
func setUpDatabase() *gorm.DB {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=%s", envConfig.DbUsername, envConfig.DbPassword, envConfig.DbName, envConfig.DbPort, envConfig.DbSSLMode)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("failed to connect to database %w", err))
	}
	return db
}

// Sets up kafka and Returns kafka Instance
func setUpKafkaProducer() *kafka.KafkaProducer {
	return kafka.NewKafkaProducer(envConfig.KafkaAddress, envConfig.KafkaTopic)
}

// Sets up Redis and Returns Redis Instance
func setUpRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     envConfig.RedisAddress,
		Password: envConfig.RedisPassword,
		DB:       envConfig.RedisDB,
	})

	return rdb
}
