package services

import (
	"errors"
	"event-booking/models"
	"event-booking/utils"
)

type UserRepository interface {
	CreateUser(*models.User) (int64, error)
	ValidateCredentials(*models.User) (bool, error)
	GetUsers() ([]models.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(u *models.User) error {
	id, err := s.repo.CreateUser(u)
	if err != nil {
		return err
	}

	u.Id = id
	return nil
}

func (s *UserService) Login(u *models.User) (string, error) {
	isValid, err := s.repo.ValidateCredentials(u)
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

func (s *UserService) GetUsers() ([]models.User, error) {
	users, err := s.repo.GetUsers()
	return users, err
}
