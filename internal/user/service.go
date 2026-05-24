package user

import (
	"fmt"

	"github.com/danielm/app_sara_backend/internal/domain"
)

type Service interface {
	GetProfile(userID string) (*domain.User, error)
	UpdateProfile(userID string, input domain.UpdateUserInput) (*domain.User, error)
	GetUser(id string) (*domain.User, error)
	ListUsers() ([]domain.User, error)
	DeleteUser(id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetProfile(userID string) (*domain.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err)
	}
	return user, nil
}

func (s *service) UpdateProfile(userID string, input domain.UpdateUserInput) (*domain.User, error) {
	user, err := s.repo.Update(userID, input)
	if err != nil {
		return nil, fmt.Errorf("update profile: %w", err)
	}
	return user, nil
}

func (s *service) GetUser(id string) (*domain.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}

func (s *service) ListUsers() ([]domain.User, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return users, nil
}

func (s *service) DeleteUser(id string) error {
	if err := s.repo.SoftDelete(id); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}
