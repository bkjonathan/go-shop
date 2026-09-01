package server

import (
	"github.com/bkjonathan/go-shop/internal/dto"
	"github.com/gin-gonic/gin"
)

func (s *Server) getProfile(ctx *gin.Context) (*dto.UserResponse, error) {
	return s.userService.GetProfile(ctx.GetUint("user_id"))
}

func (s *Server) updateProfile(ctx *gin.Context, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	return s.userService.UpdateProfile(ctx.GetUint("user_id"), req)
}
