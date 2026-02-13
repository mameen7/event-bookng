package services

import (
	"errors"
	"testing"

	"event-booking/models"
	"event-booking/services/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func createTestRegisteredEvent(id, userId, eventId int64) models.RegisterEvent {
	return models.RegisterEvent{
		Id:      id,
		UserId:  userId,
		EventId: eventId,
	}
}

func TestRegisterEvent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRegisterRepository(ctrl)
	service := NewEventRegisterService(mockRepo)

	userId := int64(10)
	eventId := int64(1)
	event := createTestEvent(eventId, 5)

	mockRepo.EXPECT().GetEventById(eventId).Return(event, nil)
	mockRepo.EXPECT().RegisterEvent(userId, eventId).Return(nil)

	err := service.RegisterEvent(userId, eventId)

	require.NoError(t, err)
}

func TestRegisterEvent_EventNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRegisterRepository(ctrl)
	service := NewEventRegisterService(mockRepo)

	userId := int64(10)
	eventId := int64(999)

	mockRepo.EXPECT().GetEventById(eventId).Return(models.Event{}, errors.New("not found"))

	err := service.RegisterEvent(userId, eventId)

	require.Error(t, err)
	assert.Equal(t, ErrEventNotFound, err)
}

func TestRegisterEvent_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRegisterRepository(ctrl)
	service := NewEventRegisterService(mockRepo)

	userId := int64(10)
	eventId := int64(1)
	event := createTestEvent(eventId, 5)
	expectedError := errors.New("registration insert failed")

	mockRepo.EXPECT().GetEventById(eventId).Return(event, nil)
	mockRepo.EXPECT().RegisterEvent(userId, eventId).Return(expectedError)

	err := service.RegisterEvent(userId, eventId)

	require.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestRegisterEvent_AlreadyRegistered(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRegisterRepository(ctrl)
	service := NewEventRegisterService(mockRepo)

	userId := int64(10)
	eventId := int64(1)
	event := createTestEvent(eventId, 5)

	mockRepo.EXPECT().GetEventById(eventId).Return(event, nil)
	mockRepo.EXPECT().RegisterEvent(userId, eventId).Return(errors.New("UNIQUE constraint failed"))

	err := service.RegisterEvent(userId, eventId)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "UNIQUE")
}

func TestCancelEvent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRegisterRepository(ctrl)
	service := NewEventRegisterService(mockRepo)

	userId := int64(10)
	eventId := int64(1)
	event := createTestEvent(eventId, 5)
	registeredEvent := createTestRegisteredEvent(100, userId, eventId)

	mockRepo.EXPECT().GetEventById(eventId).Return(event, nil)
	mockRepo.EXPECT().GetRegisteredEventById(userId, eventId).Return(registeredEvent, nil)
	mockRepo.EXPECT().DeleteRegisteredEvent(registeredEvent.Id).Return(nil)

	err := service.CancelEvent(userId, eventId)

	require.NoError(t, err)
}

func TestCancelEvent_EventNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRegisterRepository(ctrl)
	service := NewEventRegisterService(mockRepo)

	userId := int64(10)
	eventId := int64(999)

	mockRepo.EXPECT().GetEventById(eventId).Return(models.Event{}, errors.New("not found"))

	err := service.CancelEvent(userId, eventId)

	require.Error(t, err)
	assert.Equal(t, ErrEventNotFound, err)
}

func TestCancelEvent_RegistrationNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRegisterRepository(ctrl)
	service := NewEventRegisterService(mockRepo)

	userId := int64(10)
	eventId := int64(1)
	event := createTestEvent(eventId, 5)

	mockRepo.EXPECT().GetEventById(eventId).Return(event, nil)
	mockRepo.EXPECT().GetRegisteredEventById(userId, eventId).Return(models.RegisterEvent{}, errors.New("not found"))

	err := service.CancelEvent(userId, eventId)

	require.Error(t, err)
	assert.Equal(t, ErrRegisterEventNotFound, err)
}

func TestCancelEvent_DeleteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRegisterRepository(ctrl)
	service := NewEventRegisterService(mockRepo)

	userId := int64(10)
	eventId := int64(1)
	event := createTestEvent(eventId, 5)
	registeredEvent := createTestRegisteredEvent(100, userId, eventId)
	expectedError := errors.New("delete failed")

	mockRepo.EXPECT().GetEventById(eventId).Return(event, nil)
	mockRepo.EXPECT().GetRegisteredEventById(userId, eventId).Return(registeredEvent, nil)
	mockRepo.EXPECT().DeleteRegisteredEvent(registeredEvent.Id).Return(expectedError)

	err := service.CancelEvent(userId, eventId)

	require.Error(t, err)
	assert.Equal(t, expectedError, err)
}
