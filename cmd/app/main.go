package main

import (
	"log"
	"os"

	"tournament-app/internal/db"
	"tournament-app/internal/router"

	_ "tournament-app/docs" // this creates error on build and its necessary for Swagger to work

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Tournament App API
// @version 1.0
// @description This is a sample server for a tournament application.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 10.0.2.10:8080
// @BasePath /

// @securityDefinitions.basic BasicAuth
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("POSTGRES_DSN environment variable is not set")
	}

	if err := db.InitPostgres(dsn); err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}

	db.InitRedis(0)
}

func main() {
	r := gin.Default()

	// CORS
	r.Use(cors.Default())

	router.UserRoutes(r)
	router.TournamentRoutes(r)

	// Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// run on all interfaces
	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
