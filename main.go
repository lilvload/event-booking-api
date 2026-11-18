package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lilvload/event-booking-api/db"
	"github.com/lilvload/event-booking-api/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is not set. Copy .env.example to .env and set the JWT_SECRET variable.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	origins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if origins == "" {
		origins = "http://localhost:5173"
	}

	db.InitDB()
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOrigins: strings.Split(origins, ","),
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		MaxAge:       12 * time.Hour,
	}))

	routes.RegisterRoutes(server)

	server.Run(":" + port)
}
