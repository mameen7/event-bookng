package db

import (
	"testing"

	"event-booking/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	repo := NewSqlUserRepository(testDB)

	user := &models.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	id, err := repo.CreateUser(user)

	require.NoError(t, err)
	assert.Greater(t, id, int64(0))

	users, _ := repo.GetUsers()
	require.Len(t, users, 1)
	assert.NotEqual(t, "password123", users[0].Password, "Password should be hashed in database")
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	repo := NewSqlUserRepository(testDB)

	user1 := &models.User{
		Email:    "duplicate@example.com",
		Password: "password123",
	}
	user2 := &models.User{
		Email:    "duplicate@example.com",
		Password: "password456",
	}

	_, err1 := repo.CreateUser(user1)
	require.NoError(t, err1)

	_, err2 := repo.CreateUser(user2)
	require.Error(t, err2, "Should fail due to UNIQUE constraint")
	assert.Contains(t, err2.Error(), "UNIQUE")
}

func TestCreateUser_PasswordTooShort(t *testing.T) {
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	repo := NewSqlUserRepository(testDB)

	user := &models.User{
		Email:    "test@example.com",
		Password: "short",
	}

	_, err := repo.CreateUser(user)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "at least 8 characters")
}

func TestValidateCredentials_ValidPassword(t *testing.T) {
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	repo := NewSqlUserRepository(testDB)

	user := &models.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	repo.CreateUser(user)

	loginUser := &models.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	valid, err := repo.ValidateCredentials(loginUser)

	require.NoError(t, err)
	assert.True(t, valid)
	assert.Greater(t, loginUser.Id, int64(0), "User ID should be set")
}

func TestValidateCredentials_InvalidPassword(t *testing.T) {
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	repo := NewSqlUserRepository(testDB)

	user := &models.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	repo.CreateUser(user)

	loginUser := &models.User{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	valid, err := repo.ValidateCredentials(loginUser)

	require.NoError(t, err)
	assert.False(t, valid)
}

func TestValidateCredentials_UserNotFound(t *testing.T) {
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	repo := NewSqlUserRepository(testDB)

	loginUser := &models.User{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	valid, err := repo.ValidateCredentials(loginUser)

	require.Error(t, err)
	assert.False(t, valid)
}

func TestGetUsers(t *testing.T) {
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	repo := NewSqlUserRepository(testDB)

	user1 := &models.User{Email: "user1@example.com", Password: "password123"}
	user2 := &models.User{Email: "user2@example.com", Password: "password456"}

	repo.CreateUser(user1)
	repo.CreateUser(user2)

	users, err := repo.GetUsers()

	require.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "user1@example.com", users[0].Email)
	assert.Equal(t, "user2@example.com", users[1].Email)
	assert.NotEqual(t, "password123", users[0].Password)
	assert.NotEqual(t, "password456", users[1].Password)
}

func TestGetUsers_EmptyTable(t *testing.T) {
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	repo := NewSqlUserRepository(testDB)

	users, err := repo.GetUsers()

	require.NoError(t, err)
	assert.Empty(t, users)
}
