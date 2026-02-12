package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"event-booking/models"
	"event-booking/services"
)

func signup(context *gin.Context, userService *services.UserService) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot parse request data",
		})
		return
	}

	err = userService.CreateUser(&user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create user",
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
	})
}

func login(context *gin.Context, userService *services.UserService) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot parse request data",
		})
		return
	}
	token, err := userService.Login(&user)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Could not authenticate user",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Logged in successfully",
		"token":   token,
	})

}

func getAllUsers(context *gin.Context, userService *services.UserService) {
	users, err := userService.GetUsers()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve users",
		})
		return
	}
	context.JSON(http.StatusOK, users)
}
