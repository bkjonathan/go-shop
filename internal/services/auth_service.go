package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/bkjonathan/go-shop/internal/apperror"
	"github.com/bkjonathan/go-shop/internal/config"
	"github.com/bkjonathan/go-shop/internal/dto"
	"github.com/bkjonathan/go-shop/internal/models"
	"github.com/bkjonathan/go-shop/internal/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	db     *gorm.DB
	config *config.Config
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{
		db:     db,
		config: cfg,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// check if the user exists
	var existingUser models.User
	err := s.db.Where("email = ?", req.Email).First(&existingUser).Error
	if err == nil {
		return nil, apperror.Conflict("Email is already registered")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Hash Password
	hashPass, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	// Create User
	user := models.User{
		Email:     req.Email,
		Password:  hashPass,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      models.UserRoleCustomer,
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}
	// Create a cart
	cart := models.Cart{UserID: user.ID}
	if err := s.db.Create(&cart).Error; err != nil {
		fmt.Println("Unable to create cart")
	}

	// generate token
	return s.generateAuthResponse(&user)
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	var user models.User
	if err := s.db.Where("email = ? AND is_active = ?", req.Email, true).First(&user).Error; err != nil {
		return nil, apperror.Unauthorized("Invalid email or password")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, apperror.Unauthorized("Invalid email or password")
	}

	return s.generateAuthResponse(&user)
}

func (s *AuthService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	claims, err := utils.ValidateToken(req.RefreshToken, s.config.JWT.Secret)
	if err != nil {
		return nil, apperror.Unauthorized("Invalid refresh token")
	}

	var refreshToken models.RefreshToken

	if err := s.db.Where("token = ? AND expires_at > ?", req.RefreshToken, time.Now()).First(&refreshToken).Error; err != nil {
		return nil, apperror.Unauthorized("Invalid refresh token")
	}

	var user models.User
	if err := s.db.First(&user, claims.UserID).Error; err != nil {
		return nil, apperror.NotFound("User not found")
	}

	s.db.Delete(&refreshToken)
	return s.generateAuthResponse(&user)
}

func (s *AuthService) Logout(refreshToken string) error {
	return s.db.Where("token = ?", refreshToken).Delete(&models.RefreshToken{}).Error
}

func (s *AuthService) generateAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	accessToken, refreshToken, err := utils.GenerateTokenPair(&s.config.JWT, user.ID, user.Email, string(user.Role))
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
