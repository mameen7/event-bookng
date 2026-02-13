package db

import "database/sql"

// createTablesForDB creates tables in the provided database
func createTablesForDB(database *sql.DB) error {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);
	`
	_, err := database.Exec(createUsersTable)
	if err != nil {
		return err
	}

	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		datetime DATETIME NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
	`

	_, err = database.Exec(createEventsTable)
	if err != nil {
		return err
	}

	createRegistrationsTable := `
	CREATE TABLE IF NOT EXISTS registrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_id INTEGER,
		user_id INTEGER,
		FOREIGN KEY(event_id) REFERENCES events(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
	`

	_, err = database.Exec(createRegistrationsTable)
	return err
}

// createTables creates tables in the global DB
func createTables() error {
	return createTablesForDB(DB)
}
