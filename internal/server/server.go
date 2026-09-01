package server

import (
	"net/http"

	"github.com/bkjonathan/go-shop/internal/config"
	"github.com/bkjonathan/go-shop/internal/services"
	"github.com/bkjonathan/go-shop/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Server struct {
	config      *config.Config
	db          *gorm.DB
	logger      *zerolog.Logger
	authService *services.AuthService
	userService *services.UserService
}

func New(cfg *config.Config, db *gorm.DB, logger *zerolog.Logger) *Server {
	return &Server{
		config:      cfg,
		db:          db,
		logger:      logger,
		authService: services.NewAuthService(db, cfg),
		userService: services.NewUserService(db),
	}
}

func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.New()

	// Report validation errors with json field names ("first_name"), not Go ones.
	utils.RegisterValidationTagNames()

	// Add middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())

	// Add Routes
	router.GET("/health", s.healthCheck)

	api := router.Group("/api/v1")

	auth := api.Group("/auth")
	auth.POST("/register", Handle(http.StatusCreated, "User registered successfully", s.register))
	auth.POST("/login", Handle(http.StatusOK, "Login successful", s.login))
	auth.POST("/refresh", Handle(http.StatusOK, "Token refreshed successfully", s.refreshToken))
	auth.POST("/logout", HandleNoContent(http.StatusOK, "Logout successful", s.logout))

	users := api.Group("/users")
	users.Use(s.authMiddleware())
	users.GET("/me", HandleEmpty(http.StatusOK, "Profile retrieved successfully", s.getProfile))
	users.PUT("/me", Handle(http.StatusOK, "Profile updated successfully", s.updateProfile))

	return router
}

func (s *Server) healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}
	}
}
