package services

import (
	"errors"
	"event-booking/db"
	"event-booking/models"
)

var ErrForbidden = errors.New("You're not allowed to perform this action")
var ErrEventNotFound = errors.New("Event could not be retrieved")

func GetAllEvents() ([]models.Event, error) {
	events, err := db.GetEvents()
	if err != nil {
		return []models.Event{}, err
	}
	return events, nil
}

func CreateEvent(e *models.Event) error {
	id, err := db.CreateEvent(e)
	if err != nil {
		return err
	}

	e.Id = id
	return nil
}

func GetEventById(id int64) (models.Event, error) {
	e, err := db.GetEventById(id)
	if err != nil {
		return models.Event{}, err
	}
	return e, nil
}

func UpdateEvent(eventId, userId int64, updatedEvent *models.Event) error {
	event, err := db.GetEventById(eventId)
	if err != nil {
		return ErrEventNotFound
	}

	if event.UserId != userId {
		return ErrForbidden
	}

	updatedEvent.Id = eventId
	return db.UpdateEvent(updatedEvent)
}

func DeleteEvent(userId, eventId int64) error {
	event, err := db.GetEventById(eventId)
	if err != nil {
		return ErrEventNotFound
	}

	if event.UserId != userId {
		return ErrForbidden
	}

	return db.DeleteEvent(eventId)
}
