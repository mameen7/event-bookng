package routes

import (
	"event-booking/middleware"
	"event-booking/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	server *gin.Engine,
	userService *services.UserService,
	eventService *services.EventService,
	eventRegisterService *services.EventRegisterService,
) {
	authenticated := server.Group("/")
	authenticated.Use(middleware.Authenticate)
	authenticated.GET("/events", func(c *gin.Context) {
		getEvents(c, eventService)
	})
	authenticated.GET("/events/:id", func(c *gin.Context) {
		getEventById(c, eventService)
	})
	authenticated.POST("/events", func(c *gin.Context) {
		createEvent(c, eventService)
	})
	authenticated.PUT("/events/:id", func(c *gin.Context) {
		updateEvent(c, eventService)
	})
	authenticated.DELETE("/events/:id", func(c *gin.Context) {
		deleteEvent(c, eventService)
	})

	authenticated.POST("/events/:id/register", func(c *gin.Context) {
		registerEvent(c, eventRegisterService)
	})
	authenticated.DELETE("/events/:id/register", func(c *gin.Context) {
		cancelEventRegister(c, eventRegisterService)
	})

	authenticated.GET("/users", func(c *gin.Context) {
		getAllUsers(c, userService)
	})

	server.POST("/signup", func(c *gin.Context) {
		signup(c, userService)
	})
	server.POST("/login", func(c *gin.Context) {
		login(c, userService)
	})
}
