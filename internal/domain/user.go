package domain

import "time"

type UserRole string

const (
	RoleUser         UserRole = "user"
	RoleVeterinarian UserRole = "veterinarian"
	RoleAdmin        UserRole = "admin"
)

type User struct {
	ID           string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name         string    `json:"name" gorm:"not null"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"column:password_hash;not null"`
	Phone        *string   `json:"phone,omitempty" gorm:"type:varchar(20)"`
	Role         UserRole  `json:"role" gorm:"type:varchar(20);not null;default:'user'"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

type RegisterInput struct {
	Name     string   `json:"name" validate:"required"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=6"`
	Phone    *string  `json:"phone,omitempty"`
	Role        UserRole `json:"role" validate:"required"`
	AdminSecret string   `json:"admin_secret,omitempty"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserInput struct {
	Name  *string `json:"name,omitempty"`
	Phone *string `json:"phone,omitempty"`
	Email *string `json:"email,omitempty"`
}

type AuthResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
