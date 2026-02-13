package services

import (
	"errors"
	"testing"
	"time"

	"event-booking/models"
	"event-booking/services/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func createTestEvent(id, userId int64) models.Event {
	return models.Event{
		Id:          id,
		Name:        "Test Event",
		Description: "Test Description for event",
		Location:    "Test Location",
		DateTime:    time.Now().Add(24 * time.Hour),
		UserId:      userId,
	}
}

func TestGetAllEvents_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepository(ctrl)
	service := NewEventService(mockRepo)

	expectedEvents := []models.Event{
		createTestEvent(1, 1),
		createTestEvent(2, 2),
	}

	mockRepo.EXPECT().GetEvents().Return(expectedEvents, nil)

	result, err := service.GetAllEvents()

	require.NoError(t, err)
	assert.Equal(t, expectedEvents, result)
	assert.Len(t, result, 2)
}

func TestGetAllEvents_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepository(ctrl)
	service := NewEventService(mockRepo)

	expectedError := errors.New("database connection failed")
	mockRepo.EXPECT().GetEvents().Return([]models.Event{}, expectedError)

	result, err := service.GetAllEvents()

	require.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, result)
}

func TestGetAllEvents_EmptyList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepository(ctrl)
	service := NewEventService(mockRepo)

	mockRepo.EXPECT().GetEvents().Return([]models.Event{}, nil)

	result, err := service.GetAllEvents()

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetEventById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepository(ctrl)
	service := NewEventService(mockRepo)

	expectedEvent := createTestEvent(1, 1)
	mockRepo.EXPECT().GetEventById(int64(1)).Return(expectedEvent, nil)

	result, err := service.GetEventById(1)

	require.NoError(t, err)
	assert.Equal(t, expectedEvent, result)
}

func TestGetEventById_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepository(ctrl)
	service := NewEventService(mockRepo)

	mockRepo.EXPECT().GetEventById(int64(999)).Return(models.Event{}, errors.New("not found"))

	result, err := service.GetEventById(999)

	require.Error(t, err)
	assert.Equal(t, models.Event{}, result)
}

func TestCreateEvent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepository(ctrl)
	service := NewEventService(mockRepo)

	event := createTestEvent(0, 1) // ID is 0 before creation
	expectedId := int64(100)

	mockRepo.EXPECT().CreateEvent(&event).Return(expectedId, nil)

	err := service.CreateEvent(&event)

	require.NoError(t, err)
	assert.Equal(t, expectedId, event.Id)
}

func TestCreateEvent_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepository(ctrl)
	service := NewEventService(mockRepo)

	event := createTestEvent(0, 1)
	expectedError := errors.New("database insert failed")

	mockRepo.EXPECT().CreateEvent(&event).Return(int64(0), expectedError)

	err := service.CreateEvent(&event)

	require.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, int64(0), event.Id) // ID should not be set on error
}

func TestUpdateEvent_Success_OwnerCanUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepository(ctrl)
	service := NewEventService(mockRepo)

	eventId := int64(1)
	userId := int64(10)
	existingEvent := createTestEvent(eventId, userId)
	updatedEvent := createTestEvent(0, userId)
	updatedEvent.Name = "Updated Name"

	mockRepo.EXPECT().GetEventById(eventId).Return(existingEvent, nil)
	mockRepo.EXPECT().UpdateEvent(gomock.Any()).DoAndReturn(func(e *models.Event) error {
		assert.Equal(t, eventId, e.Id)
		return nil
	})

	err := service.UpdateEvent(eventId, userId, &updatedEvent)

	require.NoError(t, err)
	assert.Equal(t, eventId, updatedEvent.Id)
}

func TestUpdateEvent_Forbidden_NonOwnerCannotUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepository(ctrl)
	service := NewEventService(mockRepo)

	eventId := int64(1)
	ownerUserId := int64(10)
	otherUserId := int64(20)
	existingEvent := createTestEvent(eventId, ownerUserId)
	updatedEvent := createTestEvent(0, otherUserId)

	mockRepo.EXPECT().GetEventById(eventId).Return(existingEvent, nil)

	err := service.UpdateEvent(eventId, otherUserId, &updatedEvent)

	require.Error(t, err)
	assert.Equal(t, ErrForbidden, err)
}

func TestUpdateEvent_EventNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepository(ctrl)
	service := NewEventService(mockRepo)

	eventId := int64(999)
	userId := int64(10)
	updatedEvent := createTestEvent(0, userId)

	mockRepo.EXPECT().GetEventById(eventId).Return(models.Event{}, errors.New("not found"))

	err := service.UpdateEvent(eventId, userId, &updatedEvent)

	require.Error(t, err)
	assert.Equal(t, ErrEventNotFound, err)
}

func TestUpdateEvent_RepositoryUpdateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepository(ctrl)
	service := NewEventService(mockRepo)

	eventId := int64(1)
	userId := int64(10)
	existingEvent := createTestEvent(eventId, userId)
	updatedEvent := createTestEvent(0, userId)
	expectedError := errors.New("update failed")

	mockRepo.EXPECT().GetEventById(eventId).Return(existingEvent, nil)
	mockRepo.EXPECT().UpdateEvent(gomock.Any()).Return(expectedError)

	err := service.UpdateEvent(eventId, userId, &updatedEvent)

	require.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestDeleteEvent_Success_OwnerCanDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepository(ctrl)
	service := NewEventService(mockRepo)

	eventId := int64(1)
	userId := int64(10)
	existingEvent := createTestEvent(eventId, userId)

	mockRepo.EXPECT().GetEventById(eventId).Return(existingEvent, nil)
	mockRepo.EXPECT().DeleteEvent(eventId).Return(nil)

	err := service.DeleteEvent(userId, eventId)

	require.NoError(t, err)
}

func TestDeleteEvent_Forbidden_NonOwnerCannotDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepository(ctrl)
	service := NewEventService(mockRepo)

	eventId := int64(1)
	ownerUserId := int64(10)
	otherUserId := int64(20)
	existingEvent := createTestEvent(eventId, ownerUserId)

	mockRepo.EXPECT().GetEventById(eventId).Return(existingEvent, nil)

	err := service.DeleteEvent(otherUserId, eventId)

	require.Error(t, err)
	assert.Equal(t, ErrForbidden, err)
}

func TestDeleteEvent_EventNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepository(ctrl)
	service := NewEventService(mockRepo)

	eventId := int64(999)
	userId := int64(10)

	mockRepo.EXPECT().GetEventById(eventId).Return(models.Event{}, errors.New("not found"))

	err := service.DeleteEvent(userId, eventId)

	require.Error(t, err)
	assert.Equal(t, ErrEventNotFound, err)
}

func TestDeleteEvent_RepositoryDeleteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepository(ctrl)
	service := NewEventService(mockRepo)

	eventId := int64(1)
	userId := int64(10)
	existingEvent := createTestEvent(eventId, userId)
	expectedError := errors.New("delete failed")

	mockRepo.EXPECT().GetEventById(eventId).Return(existingEvent, nil)
	mockRepo.EXPECT().DeleteEvent(eventId).Return(expectedError)

	err := service.DeleteEvent(userId, eventId)

	require.Error(t, err)
	assert.Equal(t, expectedError, err)
}
