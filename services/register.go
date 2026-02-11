package services

import (
	"errors"
	"event-booking/db"
)

var ErrRegisterEventNotFound = errors.New("Event registration could not be retrieved")

func RegisterEvent(userId, eventId *int64) error {
	_, err := db.GetEventById(*eventId)
	if err != nil {
		return ErrEventNotFound
	}

	return db.RegisterEvent(userId, eventId)
}

func CancelEvent(userId, eventId *int64) error {
	_, err := db.GetEventById(*eventId)
	if err != nil {
		return ErrEventNotFound
	}

	registeredEvent, err := db.GetRegisteredEventById(userId, eventId)
	if err != nil {
		return ErrRegisterEventNotFound
	}

	return db.DeleteRegisteredEvent(&registeredEvent.Id)
}
