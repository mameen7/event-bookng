package db

import "event-booking/models"

type SqlEventRegisterRepository struct {
	eventRepo SqlEventRepository
}

func NewSqlEventRegisterRepository() *SqlEventRegisterRepository {
	return &SqlEventRegisterRepository{}
}

func (r *SqlEventRegisterRepository) GetEventById(id int64) (models.Event, error) {
	return r.eventRepo.GetEventById(id)
}

func (r *SqlEventRegisterRepository) RegisterEvent(userId, eventId int64) error {
	query := `
		INSERT INTO registrations (user_id, event_id)
		VALUES (?, ?);
	`
	_, err := DB.Exec(query, userId, eventId)
	return err
}

func (r *SqlEventRegisterRepository) GetRegisteredEventById(userId, eventId int64) (models.RegisterEvent, error) {
	query := `SELECT * FROM registrations WHERE user_id = ? AND event_id = ?;`
	row := DB.QueryRow(query, userId, eventId)
	var registeredEvent models.RegisterEvent
	err := row.Scan(&registeredEvent.Id, &registeredEvent.UserId, &registeredEvent.EventId)
	return registeredEvent, err
}

func (r *SqlEventRegisterRepository) DeleteRegisteredEvent(id int64) error {
	query := `DELETE FROM registrations WHERE id = ?;`
	_, err := DB.Exec(query, id)
	return err
}
