package services

import (
	"dynamodb-basic-crud/models"
	"dynamodb-basic-crud/repositories"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) CreateUser(user models.User) error {
	return s.repo.CreateUser(user)
}

func (s *UserService) UpdateUser(id string, name string, email string) error {
	err := s.repo.UpdateUser(id, name, email)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) DeleteUser(id string) error {
	err := s.repo.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}
