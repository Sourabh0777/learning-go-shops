package services

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"learning-go-shop/internal/config"
	"learning-go-shop/internal/dto"
	"learning-go-shop/internal/models"
	"learning-go-shop/internal/util"
)

type AuthService struct {
	db     *gorm.DB
	config *config.Config
}

func NewAuthService(db *gorm.DB, config *config.Config) *AuthService {
	return &AuthService{
		db:     db,
		config: config,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	var existingUser models.User
	// First() → finds the first matching record and fills the data into that variable.
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	// hashedPassword, err: = util.HashPassword(req.Password)

	// Create user
	user := models.User{
		Email:     req.Email,
		Password:  "hashedPassword",
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      models.UserRoleCustomer,
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}
	// Create cart
	cart := models.Cart{
		UserID: user.ID,
	}
	if err := s.db.Create(&cart).Error; err != nil {
		fmt.Println("Unable to create cart")
		_ = err

	}
	// Generate token
	// generate token
	return s.generateAuthResponse(&user)
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	var user models.User
	if err := s.db.Where("email=? AND is_active=?", req.Email, true).First(&user).Error; err != nil {
		return nil, errors.New("Invalid credentials")
	}
	if !util.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("Invalid credentials")
	}
	return s.generateAuthResponse(&user)
}

func (s *AuthService) RefereshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	claims, err := util.ValidateToken(req.RefreshToken, s.config.JWT.Secret)
	if err != nil {
		return nil, errors.New("Invalid refresh token")
	}
	var refereshToken models.RefreshToken
	if err := s.db.Where("token= ? AND expires_at > ?", req.RefreshToken, time.Now()).First(&refereshToken).Error; err != nil {
		return nil, errors.New("Refresh token not found or expired")
	}

	var user models.User
	if err := s.db.First(&user, claims.UserID).Error; err != nil {
		return nil, errors.New("User not found")
	}
	s.db.Delete(&refereshToken)

	return s.generateAuthResponse(&user)

}
func (s *AuthService) Logout(user *models.User) (*dto.AuthResponse, error) {
	accessToken, refreshToken, err := util.GenerateTokenPair(&s.config.JWT, user.Email, user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}
	refreshTokenModel := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.config.JWT.RefreshTokenExpires),
	}
	s.db.Create(&refreshTokenModel)
	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Role:      string(user.Role),
			IsActive:  user.IsActive,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
