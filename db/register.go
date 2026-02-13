package db

import (
	"database/sql"
	"event-booking/models"
)

type SqlEventRegisterRepository struct {
	db        *sql.DB
	eventRepo *SqlEventRepository
}

func NewSqlEventRegisterRepository(database *sql.DB) *SqlEventRegisterRepository {
	return &SqlEventRegisterRepository{
		db:        database,
		eventRepo: NewSqlEventRepository(database),
	}
}

func (r *SqlEventRegisterRepository) GetEventById(id int64) (models.Event, error) {
	return r.eventRepo.GetEventById(id)
}

func (r *SqlEventRegisterRepository) RegisterEvent(userId, eventId int64) error {
	query := `
		INSERT INTO registrations (user_id, event_id)
		VALUES (?, ?);
	`
	_, err := r.db.Exec(query, userId, eventId)
	return err
}

func (r *SqlEventRegisterRepository) GetRegisteredEventById(userId, eventId int64) (models.RegisterEvent, error) {
	query := `SELECT * FROM registrations WHERE user_id = ? AND event_id = ?;`
	row := r.db.QueryRow(query, userId, eventId)
	var registeredEvent models.RegisterEvent
	err := row.Scan(&registeredEvent.Id, &registeredEvent.EventId, &registeredEvent.UserId)
	return registeredEvent, err
}

func (r *SqlEventRegisterRepository) DeleteRegisteredEvent(id int64) error {
	query := `DELETE FROM registrations WHERE id = ?;`
	_, err := r.db.Exec(query, id)
	return err
}
