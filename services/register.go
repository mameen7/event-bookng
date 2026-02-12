package services

import (
	"errors"
	"event-booking/models"
)

type RegisterRepository interface {
	GetEventById(int64) (models.Event, error)
	GetRegisteredEventById(int64, int64) (models.RegisterEvent, error)
	RegisterEvent(int64, int64) error
	DeleteRegisteredEvent(int64) error
}

type EventRegisterService struct {
	repo RegisterRepository
}

var ErrRegisterEventNotFound = errors.New("Event registration could not be retrieved")

func NewEventRegisterService(repo RegisterRepository) *EventRegisterService {
	return &EventRegisterService{
		repo: repo,
	}
}

func (s *EventRegisterService) RegisterEvent(userId, eventId int64) error {
	_, err := s.repo.GetEventById(eventId)
	if err != nil {
		return ErrEventNotFound
	}

	return s.repo.RegisterEvent(userId, eventId)
}

func (s *EventRegisterService) CancelEvent(userId, eventId int64) error {
	_, err := s.repo.GetEventById(eventId)
	if err != nil {
		return ErrEventNotFound
	}

	registeredEvent, err := s.repo.GetRegisteredEventById(userId, eventId)
	if err != nil {
		return ErrRegisterEventNotFound
	}

	return s.repo.DeleteRegisteredEvent(registeredEvent.Id)
}
