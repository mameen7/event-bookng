package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"event-booking/db"
	"event-booking/routes"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.InitDB()
	server := gin.Default()

	routes.RegisterRoutes(server)

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Default fallback
	}

	server.Run(":" + port)
}
