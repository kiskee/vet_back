package services

import (
	"fmt"

	"github.com/danielm/app_sara_backend/internal/domain"
)

type Service interface {
	Create(input domain.CreateServiceInput) (*domain.Service, error)
	GetByID(id string) (*domain.Service, error)
	List() ([]domain.Service, error)
	Update(id string, input domain.UpdateServiceInput) (*domain.Service, error)
	Delete(id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(input domain.CreateServiceInput) (*domain.Service, error) {
	svc := &domain.Service{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Duration:    input.Duration,
		Category:    input.Category,
	}
	if err := s.repo.Create(svc); err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return svc, nil
}

func (s *service) GetByID(id string) (*domain.Service, error) {
	svc, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}
	return svc, nil
}

func (s *service) List() ([]domain.Service, error) {
	services, err := s.repo.FindAllActive()
	if err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}
	return services, nil
}

func (s *service) Update(id string, input domain.UpdateServiceInput) (*domain.Service, error) {
	svc, err := s.repo.Update(id, input)
	if err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}
	return svc, nil
}

func (s *service) Delete(id string) error {
	if err := s.repo.SoftDelete(id); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}
