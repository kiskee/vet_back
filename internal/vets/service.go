package vets

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/danielm/app_sara_backend/internal/domain"
)

type Service interface {
	CreateFromUser(input domain.CreateVetFromUserInput) (*domain.Vet, error)
	GetByID(id string) (*domain.Vet, error)
	List(lat, lng *float64, radius *int) ([]domain.Vet, error)
	Update(id, userID string, input domain.UpdateVetInput, isAdmin bool) (*domain.Vet, error)
	UpdateStatus(id, userID, status string, isAdmin bool) error
	UpdateLocation(id, userID string, input domain.UpdateVetLocationInput) error
}

type service struct {
	repo Repository
	rdb  *redis.Client
}

func NewService(repo Repository, rdb *redis.Client) Service {
	return &service{repo: repo, rdb: rdb}
}

func (s *service) CreateFromUser(input domain.CreateVetFromUserInput) (*domain.Vet, error) {
	vet := &domain.Vet{
		UserID:        input.UserID,
		Status:        "offline",
		MaxConcurrent: 1,
	}
	if err := s.repo.Create(vet); err != nil {
		return nil, fmt.Errorf("create vet: %w", err)
	}
	return vet, nil
}

func (s *service) GetByID(id string) (*domain.Vet, error) {
	vet, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}
	return vet, nil
}

func (s *service) List(lat, lng *float64, radius *int) ([]domain.Vet, error) {
	vets, err := s.repo.FindAllActive(lat, lng, radius)
	if err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}
	return vets, nil
}

func (s *service) Update(id, userID string, input domain.UpdateVetInput, isAdmin bool) (*domain.Vet, error) {
	vet, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("find vet: %w", err)
	}
	if vet == nil {
		return nil, nil
	}
	if !isAdmin && vet.UserID != userID {
		return nil, fmt.Errorf("forbidden")
	}

	return s.repo.Update(id, input)
}

func (s *service) UpdateStatus(id, userID, status string, isAdmin bool) error {
	vet, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("find vet: %w", err)
	}
	if vet == nil {
		return fmt.Errorf("vet not found")
	}
	if !isAdmin && vet.UserID != userID {
		return fmt.Errorf("forbidden")
	}

	if err := s.repo.UpdateStatus(id, status); err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	key := fmt.Sprintf("vet:%s:status", id)
	if status == "offline" {
		s.rdb.Del(context.Background(), key)
	} else {
		s.rdb.Set(context.Background(), key, status, 24*time.Hour)
	}

	return nil
}

func (s *service) UpdateLocation(id, userID string, input domain.UpdateVetLocationInput) error {
	vet, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("find vet: %w", err)
	}
	if vet == nil {
		return fmt.Errorf("vet not found")
	}
	if vet.UserID != userID {
		return fmt.Errorf("forbidden")
	}

	if err := s.repo.UpdateLocation(id, input.Lat, input.Lng); err != nil {
		return fmt.Errorf("update location: %w", err)
	}

	key := fmt.Sprintf("vet:%s:location", id)
	s.rdb.GeoAdd(context.Background(), key, &redis.GeoLocation{
		Name:      id,
		Longitude: input.Lng,
		Latitude:  input.Lat,
	})
	s.rdb.Expire(context.Background(), key, 30*time.Second)

	return nil
}
