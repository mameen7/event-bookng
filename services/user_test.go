package services

import (
	"errors"
	"testing"

	"event-booking/models"
	"event-booking/services/mocks"
	"event-booking/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func createTestUser(id int64, email, password string) models.User {
	return models.User{
		Id:       id,
		Email:    email,
		Password: password,
	}
}

func TestCreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	user := createTestUser(0, "test@example.com", "password123")
	expectedId := int64(100)

	mockRepo.EXPECT().CreateUser(&user).Return(expectedId, nil)

	err := service.CreateUser(&user)

	require.NoError(t, err)
	assert.Equal(t, expectedId, user.Id)
}

func TestCreateUser_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	user := createTestUser(0, "test@example.com", "password123")
	expectedError := errors.New("email already exists")

	mockRepo.EXPECT().CreateUser(&user).Return(int64(0), expectedError)

	err := service.CreateUser(&user)

	require.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, int64(0), user.Id)
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	user := createTestUser(0, "existing@example.com", "password123")

	mockRepo.EXPECT().CreateUser(&user).Return(int64(0), errors.New("UNIQUE constraint failed"))

	err := service.CreateUser(&user)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "UNIQUE")
}

func TestLogin_Success_ValidCredentials(t *testing.T) {
	testutil.SetupTestEnv(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	user := createTestUser(0, "test@example.com", "password123")

	mockRepo.EXPECT().ValidateCredentials(&user).Return(true, nil)

	token, err := service.Login(&user)

	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Contains(t, token, ".", "JWT should contain dots")
}

func TestLogin_InvalidCredentials(t *testing.T) {
	testutil.SetupTestEnv(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	user := createTestUser(0, "test@example.com", "wrongpassword")

	mockRepo.EXPECT().ValidateCredentials(&user).Return(false, nil)

	token, err := service.Login(&user)

	require.Error(t, err)
	assert.Equal(t, "Invalid Credentials", err.Error())
	assert.Empty(t, token)
}

func TestLogin_RepositoryError(t *testing.T) {
	testutil.SetupTestEnv(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	user := createTestUser(0, "test@example.com", "password123")
	expectedError := errors.New("database connection failed")

	mockRepo.EXPECT().ValidateCredentials(&user).Return(false, expectedError)

	token, err := service.Login(&user)

	require.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, token)
}

func TestLogin_UserNotFound(t *testing.T) {
	testutil.SetupTestEnv(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	user := createTestUser(0, "nonexistent@example.com", "password123")

	mockRepo.EXPECT().ValidateCredentials(&user).Return(false, errors.New("sql: no rows in result set"))

	token, err := service.Login(&user)

	require.Error(t, err)
	assert.Empty(t, token)
}

func TestGetUsers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	expectedUsers := []models.User{
		createTestUser(1, "user1@example.com", "hashed1"),
		createTestUser(2, "user2@example.com", "hashed2"),
	}

	mockRepo.EXPECT().GetUsers().Return(expectedUsers, nil)

	result, err := service.GetUsers()

	require.NoError(t, err)
	assert.Equal(t, expectedUsers, result)
	assert.Len(t, result, 2)
}

func TestGetUsers_EmptyList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	mockRepo.EXPECT().GetUsers().Return([]models.User{}, nil)

	result, err := service.GetUsers()

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetUsers_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	expectedError := errors.New("database connection failed")
	mockRepo.EXPECT().GetUsers().Return([]models.User{}, expectedError)

	result, err := service.GetUsers()

	require.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, result)
}
