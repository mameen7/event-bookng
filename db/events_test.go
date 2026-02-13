package db

import (
	"testing"
	"time"

	"event-booking/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateEvent(t *testing.T) {
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	repo := NewSqlEventRepository(testDB)

	event := &models.Event{
		Name:        "Tech Conference",
		Description: "Annual tech conference for developers",
		Location:    "San Francisco",
		DateTime:    time.Now().Add(24 * time.Hour),
		UserId:      1,
	}

	id, err := repo.CreateEvent(event)

	require.NoError(t, err)
	assert.Greater(t, id, int64(0))
}

func TestGetEvents(t *testing.T) {
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	repo := NewSqlEventRepository(testDB)

	event1 := &models.Event{
		Name: "Event 1", Description: "Description 1", Location: "Location 1",
		DateTime: time.Now().Add(24 * time.Hour), UserId: 1,
	}
	event2 := &models.Event{
		Name: "Event 2", Description: "Description 2", Location: "Location 2",
		DateTime: time.Now().Add(48 * time.Hour), UserId: 1,
	}

	repo.CreateEvent(event1)
	repo.CreateEvent(event2)

	events, err := repo.GetEvents()

	require.NoError(t, err)
	assert.Len(t, events, 2)
	assert.Equal(t, "Event 1", events[0].Name)
	assert.Equal(t, "Event 2", events[1].Name)
}

func TestGetEvents_EmptyTable(t *testing.T) {
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	repo := NewSqlEventRepository(testDB)

	events, err := repo.GetEvents()

	require.NoError(t, err)
	assert.Empty(t, events)
}

func TestGetEventById(t *testing.T) {
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	repo := NewSqlEventRepository(testDB)

	event := &models.Event{
		Name:        "Test Event",
		Description: "Test Description",
		Location:    "Test Location",
		DateTime:    time.Now().Add(24 * time.Hour),
		UserId:      1,
	}

	id, _ := repo.CreateEvent(event)

	retrievedEvent, err := repo.GetEventById(id)

	require.NoError(t, err)
	assert.Equal(t, event.Name, retrievedEvent.Name)
	assert.Equal(t, event.Description, retrievedEvent.Description)
	assert.Equal(t, event.Location, retrievedEvent.Location)
	assert.Equal(t, event.UserId, retrievedEvent.UserId)
}

func TestGetEventById_NotFound(t *testing.T) {
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	repo := NewSqlEventRepository(testDB)

	_, err := repo.GetEventById(999)

	require.Error(t, err)
}

func TestUpdateEvent(t *testing.T) {
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	repo := NewSqlEventRepository(testDB)

	event := &models.Event{
		Name:        "Original Name",
		Description: "Original Description",
		Location:    "Original Location",
		DateTime:    time.Now().Add(24 * time.Hour),
		UserId:      1,
	}

	id, _ := repo.CreateEvent(event)

	event.Id = id
	event.Name = "Updated Name"
	event.Description = "Updated Description"

	err := repo.UpdateEvent(event)

	require.NoError(t, err)

	retrievedEvent, _ := repo.GetEventById(id)
	assert.Equal(t, "Updated Name", retrievedEvent.Name)
	assert.Equal(t, "Updated Description", retrievedEvent.Description)
}

func TestDeleteEvent(t *testing.T) {
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	repo := NewSqlEventRepository(testDB)

	event := &models.Event{
		Name:        "Event to Delete",
		Description: "This will be deleted",
		Location:    "Delete Location",
		DateTime:    time.Now().Add(24 * time.Hour),
		UserId:      1,
	}

	id, _ := repo.CreateEvent(event)
	err := repo.DeleteEvent(id)

	require.NoError(t, err)

	_, err = repo.GetEventById(id)
	assert.Error(t, err, "Event should not exist after deletion")
}
