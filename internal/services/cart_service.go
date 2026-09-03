package services

import (
	"errors"

	"github.com/bkjonathan/go-shop/internal/dto"
	"github.com/bkjonathan/go-shop/internal/models"
	"gorm.io/gorm"
)

type CartService struct {
	db *gorm.DB
}

func NewCartService(db *gorm.DB) *CartService {
	return &CartService{db: db}
}

func (s *CartService) GetCart(userID uint) (*dto.CartResponse, error) {
	var cart models.Cart
	err := s.db.Preload("CartItems.Product.Category").
		Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		return nil, err
	}

	return s.convertToCartResponse(&cart), nil
}

func (s *CartService) AddToCart(userID uint, req *dto.AddToCartRequest) (*dto.CartResponse, error) {
	var product models.Product
	err := s.db.First(&product, req.ProductID).Error
	if err != nil {
		return nil, err
	}

	if product.Stock < req.Quantity {
		return nil, errors.New("not enough stock available")
	}

	var cart models.Cart
	if err := s.db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			cart = models.Cart{UserID: userID}
			if err := s.db.Create(&cart).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	var cartItem models.CartItem
	err = s.db.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&cartItem).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			cartItem = models.CartItem{
				CartID:    cart.ID,
				ProductID: req.ProductID,
				Quantity:  req.Quantity,
			}
			if err := s.db.Create(&cartItem).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		cartItem.Quantity += req.Quantity
		if err := s.db.Save(&cartItem).Error; err != nil {
			return nil, err
		}
	}

	return s.GetCart(userID)
}
func (s *CartService) UpdateCartItem(userID, cartItemID uint, req *dto.UpdateCartItemRequest) (*dto.CartResponse, error) {
	var cartItem models.CartItem
	err := s.db.Preload("Product").Where("id = ?", cartItemID).First(&cartItem).Error
	if err != nil {
		return nil, err
	}

	if cartItem.Cart.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	if cartItem.Product.Stock < req.Quantity {
		return nil, errors.New("not enough stock available")
	}

	cartItem.Quantity = req.Quantity
	if err := s.db.Save(&cartItem).Error; err != nil {
		return nil, err
	}

	return s.GetCart(userID)
}

func (s *CartService) RemoveCartItem(userID, cartItemID uint) (*dto.CartResponse, error) {
	var cartItem models.CartItem
	err := s.db.Preload("Cart").Where("id = ?", cartItemID).First(&cartItem).Error
	if err != nil {
		return nil, err
	}

	if cartItem.Cart.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	if err := s.db.Delete(&cartItem).Error; err != nil {
		return nil, err
	}

	return s.GetCart(userID)
}
func (s *CartService) convertToCartResponse(cart *models.Cart) *dto.CartResponse {

	cartItems := make([]dto.CartItemResponse, len(cart.CartItems)) // memory allocation
	var total float64

	for i := range cart.CartItems {
		subtotal := float64(cart.CartItems[i].Quantity) * cart.CartItems[i].Product.Price
		total += subtotal

		cartItems[i] = dto.CartItemResponse{
			ID: cart.CartItems[i].ID,
			Product: dto.ProductResponse{
				ID:          cart.CartItems[i].Product.ID,
				CategoryId:  cart.CartItems[i].Product.CategoryID,
				Name:        cart.CartItems[i].Product.Name,
				Description: cart.CartItems[i].Product.Description,
				Price:       cart.CartItems[i].Product.Price,
				Stock:       cart.CartItems[i].Product.Stock,
				SKU:         cart.CartItems[i].Product.SKU,
				IsActive:    &cart.CartItems[i].Product.IsActive,
				Category: dto.CategoryResponse{
					ID:          cart.CartItems[i].Product.Category.ID,
					Name:        cart.CartItems[i].Product.Category.Name,
					Description: cart.CartItems[i].Product.Category.Description,
					IsActive:    cart.CartItems[i].Product.Category.IsActive,
				},
			},
			Quantity:  cart.CartItems[i].Quantity,
			SubTotal:  subtotal,
			CreatedAt: cart.CartItems[i].CreatedAt,
			UpdatedAt: cart.CartItems[i].UpdatedAt,
		}
	}

	return &dto.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CartItems: cartItems,
		Total:     total,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}
}
