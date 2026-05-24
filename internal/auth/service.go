package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/danielm/app_sara_backend/internal/domain"
	userRepoPkg "github.com/danielm/app_sara_backend/internal/user"
)

type Service interface {
	Register(input domain.RegisterInput) (*domain.AuthResponse, error)
	Login(input domain.LoginInput) (*domain.AuthResponse, error)
	Refresh(refreshToken string) (*domain.AuthResponse, error)
}

type service struct {
	userRepo         userRepoPkg.Repository
	jwtSecret        string
	jwtRefreshSecret string
	adminSecret      string
}

func NewService(userRepo userRepoPkg.Repository, jwtSecret, jwtRefreshSecret, adminSecret string) Service {
	return &service{
		userRepo:         userRepo,
		jwtSecret:        jwtSecret,
		jwtRefreshSecret: jwtRefreshSecret,
		adminSecret:      adminSecret,
	}
}

func (s *service) Register(input domain.RegisterInput) (*domain.AuthResponse, error) {
	existing, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("check email: %w", err)
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	if input.Role == domain.RoleAdmin {
		if input.AdminSecret == "" || input.AdminSecret != s.adminSecret {
			return nil, errors.New("invalid admin secret")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	now := time.Now()
	user := &domain.User{
		ID:           uuid.New().String(),
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		Phone:        input.Phone,
		Role:         input.Role,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return &domain.AuthResponse{
		User:         *user,
		AccessToken:  tokens.accessToken,
		RefreshToken: tokens.refreshToken,
	}, nil
}

func (s *service) Login(input domain.LoginInput) (*domain.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return &domain.AuthResponse{
		User:         *user,
		AccessToken:  tokens.accessToken,
		RefreshToken: tokens.refreshToken,
	}, nil
}

func (s *service) Refresh(refreshToken string) (*domain.AuthResponse, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &domain.Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtRefreshSecret), nil
	})
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(*domain.Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	if user == nil || !user.IsActive {
		return nil, errors.New("user not found or inactive")
	}

	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return &domain.AuthResponse{
		User:         *user,
		AccessToken:  tokens.accessToken,
		RefreshToken: tokens.refreshToken,
	}, nil
}

type tokenPair struct {
	accessToken  string
	refreshToken string
}

func (s *service) generateTokens(user *domain.User) (*tokenPair, error) {
	accessClaims := domain.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	refreshClaims := domain.Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(s.jwtRefreshSecret))
	if err != nil {
		return nil, err
	}

	return &tokenPair{
		accessToken:  accessToken,
		refreshToken: refreshToken,
	}, nil
}
