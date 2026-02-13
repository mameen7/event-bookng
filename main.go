package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"

	"event-booking/db"
	"event-booking/routes"
	"event-booking/services"
	"event-booking/utils"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("futuredate", utils.ValidateFutureDate)
	}

	db.InitDB()
	eventRepo := db.NewSqlEventRepository(db.DB)
	eventRegisterRepo := db.NewSqlEventRegisterRepository(db.DB)
	userRepo := db.NewSqlUserRepository(db.DB)

	eventService := services.NewEventService(eventRepo)
	eventRegisterService := services.NewEventRegisterService(eventRegisterRepo)
	userService := services.NewUserService(userRepo)

	server := gin.Default()
	routes.RegisterRoutes(server, userService, eventService, eventRegisterService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	server.Run(":" + port)
}
