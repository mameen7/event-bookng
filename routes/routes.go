package routes

import (
	"event-booking/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	authenticated := server.Group("/")
	authenticated.Use(middleware.Authenticate)
	authenticated.GET("/events", getEvents)
	authenticated.GET("/events/:id", getEventById)
	authenticated.POST("/events", createEvent)
	authenticated.PUT("/events/:id", updateEvent)
	authenticated.DELETE("/events/:id", deleteEvent)

	authenticated.POST("/events/:id/register", registerEvent)
	authenticated.DELETE("/events/:id/register", cancelEventRegister)

	authenticated.GET("/users", getAllUsers)

	server.POST("/signup", signup)
	server.POST("/login", login)
}
