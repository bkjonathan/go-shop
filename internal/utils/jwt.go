package utils

import (
	"errors"
	"time"

	"github.com/bkjonathan/go-shop/internal/config"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateTokenPair(cfg *config.JWTConfig, userID uint, email, role string) (accessToken, refreshToken string, err error) {
	now := time.Now()

	accessToken, err = signToken(cfg, userID, email, role, now, cfg.ExpiresIn)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = signToken(cfg, userID, email, role, now, cfg.RefreshTokenExpires)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func signToken(cfg *config.JWTConfig, userID uint, email, role string, issuedAt time.Time, expiresIn time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(cfg.Secret))
}

// ValidateToken check Tokens
func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
