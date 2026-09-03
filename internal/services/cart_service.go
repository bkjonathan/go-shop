package services

import (
	"errors"

	"github.com/bkjonathan/go-shop/internal/apperror"
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

// GetCart returns the user's cart, or an empty one when they have never added
// anything. A customer who has not shopped yet is not a 404.
func (s *CartService) GetCart(userID uint) (*dto.CartResponse, error) {
	var cart models.Cart
	err := s.db.Preload("CartItems.Product.Category").
		Preload("CartItems.Product.Images").
		Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dto.CartResponse{UserID: userID, CartItems: []dto.CartItemResponse{}}, nil
		}
		return nil, err
	}

	return convertToCartResponse(&cart), nil
}

func (s *CartService) AddToCart(userID uint, req *dto.AddToCartRequest) (*dto.CartResponse, error) {
	var product models.Product
	if err := s.db.First(&product, req.ProductID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("Product not found")
		}
		return nil, err
	}

	cart, err := s.cartFor(userID)
	if err != nil {
		return nil, err
	}

	var cartItem models.CartItem
	err = s.db.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&cartItem).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Adding to a line that is already in the cart tops it up, so the stock
	// check has to cover what is in there already.
	quantity := req.Quantity
	if err == nil {
		quantity += cartItem.Quantity
	}
	if product.Stock < quantity {
		return nil, apperror.BadRequest("Not enough stock available")
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		cartItem = models.CartItem{
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		if err := s.db.Create(&cartItem).Error; err != nil {
			return nil, err
		}
	} else {
		cartItem.Quantity = quantity
		if err := s.db.Save(&cartItem).Error; err != nil {
			return nil, err
		}
	}

	return s.GetCart(userID)
}

func (s *CartService) UpdateCartItem(userID, cartItemID uint, req *dto.UpdateCartItemRequest) (*dto.CartResponse, error) {
	cartItem, err := s.ownedCartItem(userID, cartItemID)
	if err != nil {
		return nil, err
	}

	if cartItem.Product.Stock < req.Quantity {
		return nil, apperror.BadRequest("Not enough stock available")
	}

	cartItem.Quantity = req.Quantity
	if err := s.db.Save(cartItem).Error; err != nil {
		return nil, err
	}

	return s.GetCart(userID)
}

func (s *CartService) RemoveCartItem(userID, cartItemID uint) (*dto.CartResponse, error) {
	cartItem, err := s.ownedCartItem(userID, cartItemID)
	if err != nil {
		return nil, err
	}

	if err := s.db.Delete(cartItem).Error; err != nil {
		return nil, err
	}

	return s.GetCart(userID)
}

func (s *CartService) ClearCart(userID uint) error {
	cart, err := s.cartFor(userID)
	if err != nil {
		return err
	}

	return s.db.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error
}

// cartFor returns the user's cart, creating it the first time they add
// something.
func (s *CartService) cartFor(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := s.db.Where("user_id = ?", userID).First(&cart).Error
	if err == nil {
		return &cart, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	cart = models.Cart{UserID: userID}
	if err := s.db.Create(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

// ownedCartItem loads a cart item together with the cart it belongs to and the
// product it points at, and refuses it when it is not this user's. Another
// user's item is reported as missing rather than forbidden, so the endpoint
// does not confirm that the id exists.
func (s *CartService) ownedCartItem(userID, cartItemID uint) (*models.CartItem, error) {
	var cartItem models.CartItem
	err := s.db.Preload("Cart").Preload("Product").First(&cartItem, cartItemID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("Cart item not found")
		}
		return nil, err
	}

	if cartItem.Cart.UserID != userID {
		return nil, apperror.NotFound("Cart item not found")
	}

	return &cartItem, nil
}

func convertToCartResponse(cart *models.Cart) *dto.CartResponse {
	cartItems := make([]dto.CartItemResponse, len(cart.CartItems))
	var total float64

	for i := range cart.CartItems {
		item := &cart.CartItems[i]
		subtotal := float64(item.Quantity) * item.Product.Price
		total += subtotal

		cartItems[i] = dto.CartItemResponse{
			ID:        item.ID,
			Product:   toProductResponse(&item.Product),
			Quantity:  item.Quantity,
			SubTotal:  subtotal,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
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
