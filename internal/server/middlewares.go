package server

import (
	"strings"

	"github.com/bkjonathan/go-shop/internal/models"
	"github.com/bkjonathan/go-shop/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(ctx, "Authorization header required")
			ctx.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.UnauthorizedResponse(ctx, "Invalid authorization header format")
			ctx.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenParts[1], s.config.JWT.Secret)
		if err != nil {
			utils.UnauthorizedResponse(ctx, "Invalid Token")
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("user_email", claims.Email)
		ctx.Set("user_role", claims.Role)

		ctx.Next()
	}
}

func (s *Server) adminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("user_role")
		if !exists {
			utils.ForbiddenResponse(ctx, "Forbidden")
			ctx.Abort()
			return
		}

		if role != string(models.UserRoleAdmin) {
			utils.ForbiddenResponse(ctx, "Forbidden")
			ctx.Abort()
			return
		}
	}
}
