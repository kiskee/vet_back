package vets

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/danielm/app_sara_backend/internal/domain"
)

type Repository interface {
	Create(vet *domain.Vet) error
	FindByID(id string) (*domain.Vet, error)
	FindByUserID(userID string) (*domain.Vet, error)
	FindAllActive(lat, lng *float64, radius *int) ([]domain.Vet, error)
	Update(id string, input domain.UpdateVetInput) (*domain.Vet, error)
	UpdateStatus(id string, status string) error
	UpdateLocation(id string, lat, lng float64) error
}

type repository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewRepository(db *gorm.DB, rdb *redis.Client) Repository {
	return &repository{db: db, rdb: rdb}
}

func (r *repository) Create(vet *domain.Vet) error {
	return r.db.Create(vet).Error
}

func (r *repository) FindByID(id string) (*domain.Vet, error) {
	var vet domain.Vet
	err := r.db.First(&vet, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("find by id: %w", err)
	}
	return &vet, nil
}

func (r *repository) FindByUserID(userID string) (*domain.Vet, error) {
	var vet domain.Vet
	err := r.db.Where("user_id = ?", userID).First(&vet).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("find by user id: %w", err)
	}
	return &vet, nil
}

func (r *repository) FindAllActive(lat, lng *float64, radius *int) ([]domain.Vet, error) {
	var vets []domain.Vet

	query := r.db.Where("is_active = ?", true)

	if lat != nil && lng != nil && radius != nil {
		raw := fmt.Sprintf(`
			SELECT v.*, 
				ST_Distance(
					location,
					ST_SetSRID(ST_MakePoint(%f, %f), 4326)
				) AS distance
			FROM vets v
			WHERE v.is_active = true 
				AND ST_DWithin(
					location,
					ST_SetSRID(ST_MakePoint(%f, %f), 4326),
					%d
				)
			ORDER BY distance
		`, *lng, *lat, *lng, *lat, *radius)

		if err := r.db.Raw(raw).Scan(&vets).Error; err != nil {
			return nil, fmt.Errorf("find near: %w", err)
		}
		return vets, nil
	}

	if err := query.Order("created_at DESC").Find(&vets).Error; err != nil {
		return nil, fmt.Errorf("find all: %w", err)
	}
	return vets, nil
}

func (r *repository) Update(id string, input domain.UpdateVetInput) (*domain.Vet, error) {
	updates := map[string]any{}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.ClinicName != nil {
		updates["clinic_name"] = *input.ClinicName
	}
	if input.ConsultationFee != nil {
		updates["consultation_fee"] = *input.ConsultationFee
	}
	if input.MaxConcurrent != nil {
		updates["max_concurrent"] = *input.MaxConcurrent
	}

	if len(updates) > 0 {
		if err := r.db.Model(&domain.Vet{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("update: %w", err)
		}
	}

	return r.FindByID(id)
}

func (r *repository) UpdateStatus(id string, status string) error {
	return r.db.Model(&domain.Vet{}).Where("id = ?", id).Update("status", status).Error
}

func (r *repository) UpdateLocation(id string, lat, lng float64) error {
	raw := fmt.Sprintf(`
		UPDATE vets 
		SET location = ST_SetSRID(ST_MakePoint(%f, %f), 4326),
		    updated_at = NOW()
		WHERE id = '%s'
	`, lng, lat, id)

	return r.db.Exec(raw).Error
}
