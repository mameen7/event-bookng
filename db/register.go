package db

import "event-booking/models"

func RegisterEvent(userId, eventId *int64) error {
	query := `
		INSERT INTO registrations (user_id, event_id)
		VALUES (?, ?);
	`
	_, err := DB.Exec(query, *userId, *eventId)
	return err
}

func GetRegisteredEventById(userId, eventId *int64) (models.RegisterEvent, error) {
	query := `SELECT * FROM registrations WHERE user_id = ? AND event_id = ?;`
	row := DB.QueryRow(query, *userId, *eventId)
	var registeredEvent models.RegisterEvent
	err := row.Scan(&registeredEvent.Id, &registeredEvent.UserId, &registeredEvent.EventId)
	return registeredEvent, err
}

func DeleteRegisteredEvent(id *int64) error {
	query := `DELETE FROM registrations WHERE id = ?;`
	_, err := DB.Exec(query, id)
	return err
}
