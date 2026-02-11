package routes

import (
	"errors"
	"event-booking/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func registerEvent(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse eventId",
		})
		return
	}

	err = services.RegisterEvent(&userId, &eventId)
	if err != nil {
		if errors.Is(err, services.ErrEventNotFound) {
			context.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{
				"message": "Could not register for event",
			})
		}
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Event has been registered successfully",
	})
}

func cancelEventRegister(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse eventId",
		})
		return
	}

	err = services.CancelEvent(&userId, &eventId)
	if err != nil {
		if errors.Is(err, services.ErrEventNotFound) || errors.Is(err, services.ErrRegisterEventNotFound) {
			context.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{
				"message": "Could not register for event",
			})
		}
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Event registration has been canceled successfully",
	})
}
