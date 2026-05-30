package domain

import "time"

type Service struct {
	ID          string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	Duration    int       `json:"duration,omitempty"`
	Category    string    `json:"category,omitempty" gorm:"type:varchar(50)"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Service) TableName() string {
	return "services"
}

type CreateServiceInput struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Duration    int     `json:"duration"`
	Category    string  `json:"category"`
}

type UpdateServiceInput struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price" validate:"omitempty,gt=0"`
	Duration    *int     `json:"duration"`
	Category    *string  `json:"category"`
}
