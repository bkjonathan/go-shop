package server

import (
	"net/http"

	"github.com/bkjonathan/go-shop/internal/config"
	"github.com/bkjonathan/go-shop/internal/interfaces"
	"github.com/bkjonathan/go-shop/internal/providers"
	"github.com/bkjonathan/go-shop/internal/services"
	"github.com/bkjonathan/go-shop/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Server struct {
	config         *config.Config
	db             *gorm.DB
	logger         *zerolog.Logger
	authService    *services.AuthService
	userService    *services.UserService
	productService *services.ProductService
	uploadService  *services.UploadService
	cartService    *services.CartService
	orderService   *services.OrderService
}

// TODO: Add a NewServer function to initialize the server with its dependencies
func New(cfg *config.Config, db *gorm.DB, logger *zerolog.Logger) *Server {
	var uploadProvider interfaces.UploadProvider
	if cfg.Upload.UploadProvider == "s3" {
		uploadProvider = providers.NewS3Provider(cfg)
	} else {
		uploadProvider = providers.NewLocalUploadProvider(cfg.Upload.Path)
	}

	return &Server{
		config:         cfg,
		db:             db,
		logger:         logger,
		authService:    services.NewAuthService(db, cfg),
		userService:    services.NewUserService(db),
		productService: services.NewProductService(db),
		uploadService:  services.NewUploadService(uploadProvider),
		cartService:    services.NewCartService(db),
		orderService:   services.NewOrderService(db),
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
	router.Static("/uploads", "./uploads")

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

	// Category routes
	categories := api.Group("/categories")
	categories.Use(s.authMiddleware())
	categories.POST("/", s.adminMiddleware(), Handle(http.StatusCreated, "Category created successfully", s.createCategory))
	categories.GET("/", HandleEmpty(http.StatusOK, "Categories retrieved successfully", s.listCategories))
	categories.PUT("/:id", s.adminMiddleware(), Handle(http.StatusOK, "Category updated successfully", s.updateCategory))
	categories.DELETE("/:id", s.adminMiddleware(), HandleEmptyNoContent(http.StatusOK, "Category deleted successfully", s.deleteCategory))

	// Product routes
	products := api.Group("/products")
	products.Use(s.authMiddleware())
	products.POST("/", s.adminMiddleware(), Handle(http.StatusCreated, "Product created successfully", s.createProduct))
	products.GET("/", HandleEmpty(http.StatusOK, "Products retrieved successfully", s.listProducts))
	products.GET("/:id", HandleEmpty(http.StatusOK, "Product retrieved successfully", s.getProduct))
	products.PUT("/:id", s.adminMiddleware(), Handle(http.StatusOK, "Product updated successfully", s.updateProduct))
	products.DELETE("/:id", s.adminMiddleware(), HandleEmptyNoContent(http.StatusOK, "Product deleted successfully", s.deleteProduct))
	products.POST("/:id/images", s.adminMiddleware(), HandleEmpty(http.StatusCreated, "Product image uploaded successfully", s.uploadProductImage))

	// Cart routes — always scoped to the authenticated user, never to an id in the path.
	cart := api.Group("/cart")
	cart.Use(s.authMiddleware())
	cart.GET("", HandleEmpty(http.StatusOK, "Cart retrieved successfully", s.getCart))
	cart.POST("/items", Handle(http.StatusCreated, "Item added to cart successfully", s.addToCart))
	cart.PUT("/items/:id", Handle(http.StatusOK, "Cart item updated successfully", s.updateCartItem))
	cart.DELETE("/items/:id", HandleEmpty(http.StatusOK, "Cart item removed successfully", s.removeCartItem))
	cart.DELETE("", HandleEmptyNoContent(http.StatusOK, "Cart cleared successfully", s.clearCart))

	// Order routes — like the cart, always scoped to the authenticated user.
	orders := api.Group("/orders")
	orders.Use(s.authMiddleware())
	orders.POST("", HandleEmpty(http.StatusCreated, "Order created successfully", s.createOrder))
	orders.GET("", HandleEmpty(http.StatusOK, "Orders retrieved successfully", s.getOrders))
	orders.GET("/:id", HandleEmpty(http.StatusOK, "Order retrieved successfully", s.getOrder))

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
