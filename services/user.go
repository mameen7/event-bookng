package services

import (
	"errors"
	"event-booking/db"
	"event-booking/models"
	"event-booking/utils"
)

func CreateUser(u *models.User) error {
	id, err := db.CreateUser(u)
	if err != nil {
		return err
	}

	u.Id = id
	return nil
}

func Login(u *models.User) (string, error) {
	isValid, err := db.ValidateCredentials(u)
	if err != nil {
		return "", err
	}
	if !isValid {
		return "", errors.New("Invalid Credentials")
	}

	token, err := utils.GenerateToken(u.Email, u.Id)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GetUsers() ([]models.User, error) {
	users, err := db.GetUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}
