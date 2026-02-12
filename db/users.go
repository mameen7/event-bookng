package db

import (
	"event-booking/models"
	"event-booking/utils"
)

type SqlUserRepository struct{}

func NewSqlUserRepository() *SqlUserRepository {
	return &SqlUserRepository{}
}

func (r *SqlUserRepository) CreateUser(u *models.User) (int64, error) {
	query := `
	INSERT INTO users (email, password)
	VALUES (?, ?);
	`
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return 0, err
	}
	result, err := DB.Exec(query, u.Email, hashedPassword)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int64(id), err
}

func (r *SqlUserRepository) ValidateCredentials(u *models.User) (bool, error) {
	query := `
	SELECT id, password FROM users WHERE email = ?
	`
	row := DB.QueryRow(query, u.Email)

	var retrievedPassword string
	err := row.Scan(&u.Id, &retrievedPassword)
	if err != nil {
		return false, err
	}

	return utils.CheckPasswordHash(u.Password, retrievedPassword), nil

}

func (r *SqlUserRepository) GetUsers() ([]models.User, error) {
	query := `SELECT * FROM users;`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	users := []models.User{}
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.Id, &u.Email, &u.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
