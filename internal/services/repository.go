package services

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/danielm/app_sara_backend/internal/domain"
)

type Repository interface {
	Create(svc *domain.Service) error
	FindByID(id string) (*domain.Service, error)
	FindAllActive() ([]domain.Service, error)
	Update(id string, input domain.UpdateServiceInput) (*domain.Service, error)
	SoftDelete(id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(svc *domain.Service) error {
	return r.db.Create(svc).Error
}

func (r *repository) FindByID(id string) (*domain.Service, error) {
	var svc domain.Service
	err := r.db.First(&svc, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("find by id: %w", err)
	}
	return &svc, nil
}

func (r *repository) FindAllActive() ([]domain.Service, error) {
	var services []domain.Service
	err := r.db.Where("is_active = ?", true).Order("name ASC").Find(&services).Error
	return services, err
}

func (r *repository) Update(id string, input domain.UpdateServiceInput) (*domain.Service, error) {
	updates := map[string]any{}
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Price != nil {
		updates["price"] = *input.Price
	}
	if input.Duration != nil {
		updates["duration"] = *input.Duration
	}
	if input.Category != nil {
		updates["category"] = *input.Category
	}

	if len(updates) > 0 {
		if err := r.db.Model(&domain.Service{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("update: %w", err)
		}
	}

	return r.FindByID(id)
}

func (r *repository) SoftDelete(id string) error {
	return r.db.Model(&domain.Service{}).Where("id = ?", id).Update("is_active", false).Error
}
