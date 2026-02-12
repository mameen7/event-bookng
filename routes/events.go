package routes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"event-booking/models"
	"event-booking/services"
)

func getEvents(context *gin.Context, eventService *services.EventService) {
	events, err := eventService.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve events",
		})
		return
	}
	context.JSON(http.StatusOK, events)
}

func getEventById(context *gin.Context, eventService *services.EventService) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse event id",
		})
		return
	}

	event, err := eventService.GetEventById(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve event",
		})
	}

	context.JSON(http.StatusOK, event)
}

func createEvent(context *gin.Context, eventService *services.EventService) {
	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot parse request data",
		})
		return
	}

	event.UserId = context.GetInt64("userId")
	err = eventService.CreateEvent(&event)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create event",
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Event created successfully",
		"event":   event,
	})
}

func updateEvent(context *gin.Context, eventService *services.EventService) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse event id",
		})
		return
	}

	var updatedEvent models.Event
	err = context.ShouldBindJSON(&updatedEvent)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot parse request data",
		})
		return
	}

	userId := context.GetInt64("userId")
	err = eventService.UpdateEvent(eventId, userId, &updatedEvent)
	if err != nil {
		if errors.Is(err, services.ErrEventNotFound) {
			context.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
		} else if errors.Is(err, services.ErrForbidden) {
			context.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
			})
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{
				"message": "Event could not be updated",
			})
		}
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Event has been updated successfully",
	})
}

func deleteEvent(context *gin.Context, eventService *services.EventService) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse event id",
		})
		return
	}

	UserId := context.GetInt64("userId")
	err = eventService.DeleteEvent(UserId, eventId)
	if err != nil {
		if errors.Is(err, services.ErrEventNotFound) {
			context.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
		} else if errors.Is(err, services.ErrForbidden) {
			context.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
			})
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{
				"message": "Event could not be deleted",
			})
		}
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Event has been deleted successfully",
	})
}
