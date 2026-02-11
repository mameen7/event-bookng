package db

import (
	"errors"
	"event-booking/models"
)

var ErrEventNotFound = errors.New("event could not be found")

func CreateEvent(e *models.Event) (int64, error) {
	query := `
	INSERT INTO events (name, description, location, datetime, user_id)
	VALUES (?, ?, ?, ?, ?);
	`
	result, err := DB.Exec(query, e.Name, e.Description, e.Location, e.DateTime, e.UserId)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func GetEvents() ([]models.Event, error) {
	query := `SELECT * FROM events;`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []models.Event{}
	for rows.Next() {
		var e models.Event
		err = rows.Scan(&e.Id, &e.Name, &e.Description, &e.Location, &e.DateTime, &e.UserId)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func GetEventById(id int64) (models.Event, error) {
	query := `SELECT * FROM events WHERE id = ?`
	row := DB.QueryRow(query, id)

	var e models.Event
	err := row.Scan(&e.Id, &e.Name, &e.Description, &e.Location, &e.DateTime, &e.UserId)
	return e, err
}

func UpdateEvent(e *models.Event) error {
	query := `
	UPDATE events
	SET name = ?, description = ?, location = ?, datetime = ?
	WHERE id = ?
	`
	_, err := DB.Exec(query, e.Name, e.Description, e.Location, e.DateTime, e.Id)
	return err
}

func DeleteEvent(id int64) error {
	query := `
	DELETE FROM events WHERE id = ?;
	`
	_, err := DB.Exec(query, id)
	return err
}
