package user

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/danielm/app_sara_backend/internal/domain"
)

type Repository interface {
	Create(user *domain.User) error
	FindByID(id string) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	FindAll() ([]domain.User, error)
	Update(id string, input domain.UpdateUserInput) (*domain.User, error)
	SoftDelete(id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *repository) FindByID(id string) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("find by id: %w", err)
	}
	return &user, nil
}

func (r *repository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("find by email: %w", err)
	}
	return &user, nil
}

func (r *repository) FindAll() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Order("created_at DESC").Find(&users).Error
	return users, err
}

func (r *repository) Update(id string, input domain.UpdateUserInput) (*domain.User, error) {
	updates := map[string]any{}
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Phone != nil {
		updates["phone"] = *input.Phone
	}
	if input.Email != nil {
		updates["email"] = *input.Email
	}

	if len(updates) > 0 {
		if err := r.db.Model(&domain.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("update: %w", err)
		}
	}

	return r.FindByID(id)
}

func (r *repository) SoftDelete(id string) error {
	return r.db.Model(&domain.User{}).Where("id = ?", id).Update("is_active", false).Error
}
