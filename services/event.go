package services

import (
	"errors"
	"event-booking/models"
)

type EventRepository interface {
	GetEvents() ([]models.Event, error)
	GetEventById(int64) (models.Event, error)
	CreateEvent(*models.Event) (int64, error)
	UpdateEvent(*models.Event) error
	DeleteEvent(int64) error
}

type EventService struct {
	repo EventRepository
}

var ErrForbidden = errors.New("You're not allowed to perform this action")
var ErrEventNotFound = errors.New("Event could not be retrieved")

func NewEventService(repo EventRepository) *EventService {
	return &EventService{
		repo: repo,
	}
}

func (s *EventService) GetAllEvents() ([]models.Event, error) {
	events, err := s.repo.GetEvents()
	if err != nil {
		return []models.Event{}, err
	}
	return events, nil
}

func (s *EventService) CreateEvent(e *models.Event) error {
	id, err := s.repo.CreateEvent(e)
	if err != nil {
		return err
	}

	e.Id = id
	return nil
}

func (s *EventService) GetEventById(id int64) (models.Event, error) {
	e, err := s.repo.GetEventById(id)
	if err != nil {
		return models.Event{}, err
	}
	return e, nil
}

func (s *EventService) UpdateEvent(eventId, userId int64, updatedEvent *models.Event) error {
	event, err := s.repo.GetEventById(eventId)
	if err != nil {
		return ErrEventNotFound
	}

	if event.UserId != userId {
		return ErrForbidden
	}

	updatedEvent.Id = eventId
	return s.repo.UpdateEvent(updatedEvent)
}

func (s *EventService) DeleteEvent(userId, eventId int64) error {
	event, err := s.repo.GetEventById(eventId)
	if err != nil {
		return ErrEventNotFound
	}

	if event.UserId != userId {
		return ErrForbidden
	}

	return s.repo.DeleteEvent(eventId)
}
