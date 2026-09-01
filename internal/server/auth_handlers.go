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
