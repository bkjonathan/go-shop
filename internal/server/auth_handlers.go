package server

import (
	"github.com/bkjonathan/go-shop/internal/dto"
	"github.com/gin-gonic/gin"
)

func (s *Server) register(ctx *gin.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	return s.authService.Register(req)
}

func (s *Server) login(ctx *gin.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	return s.authService.Login(req)
}

func (s *Server) refreshToken(ctx *gin.Context, req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	return s.authService.RefreshToken(req)
}

func (s *Server) logout(ctx *gin.Context, req *dto.RefreshTokenRequest) error {
	return s.authService.Logout(req.RefreshToken)
}
