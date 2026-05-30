package domain

import "time"

type Vet struct {
	ID              string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID          string    `json:"user_id" gorm:"type:uuid;uniqueIndex;not null"`
	Description     string    `json:"description,omitempty"`
	ClinicName      string    `json:"clinic_name,omitempty"`
	ConsultationFee float64   `json:"consultation_fee,omitempty" gorm:"type:decimal(10,2)"`
	MaxConcurrent   int       `json:"max_concurrent" gorm:"default:1"`
	Status          string    `json:"status" gorm:"type:varchar(20);default:offline"`
	RatingAvg       float64   `json:"rating_avg" gorm:"type:decimal(2,1);default:0"`
	ReviewsCount    int       `json:"reviews_count" gorm:"default:0"`
	IsActive        bool      `json:"is_active" gorm:"default:true"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (Vet) TableName() string {
	return "vets"
}

type CreateVetFromUserInput struct {
	UserID string
}

type UpdateVetInput struct {
	Description     *string  `json:"description"`
	ClinicName      *string  `json:"clinic_name"`
	ConsultationFee *float64 `json:"consultation_fee" validate:"omitempty,gt=0"`
	MaxConcurrent   *int     `json:"max_concurrent" validate:"omitempty,gt=0"`
}

type UpdateVetStatusInput struct {
	Status string `json:"status" validate:"required,oneof=available busy offline"`
}

type UpdateVetLocationInput struct {
	Lat float64 `json:"lat" validate:"required"`
	Lng float64 `json:"lng" validate:"required"`
}
