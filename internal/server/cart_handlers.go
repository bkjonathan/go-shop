package server

import (
	"strconv"

	"github.com/bkjonathan/go-shop/internal/apperror"
	"github.com/bkjonathan/go-shop/internal/dto"
	"github.com/gin-gonic/gin"
)

func (s *Server) getCart(ctx *gin.Context) (*dto.CartResponse, error) {
	return s.cartService.GetCart(ctx.GetUint("user_id"))
}

func (s *Server) addToCart(ctx *gin.Context, req *dto.AddToCartRequest) (*dto.CartResponse, error) {
	return s.cartService.AddToCart(ctx.GetUint("user_id"), req)
}

func (s *Server) updateCartItem(ctx *gin.Context, req *dto.UpdateCartItemRequest) (*dto.CartResponse, error) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return nil, apperror.BadRequest("Invalid cart item ID")
	}

	return s.cartService.UpdateCartItem(ctx.GetUint("user_id"), uint(id), req)
}

func (s *Server) removeCartItem(ctx *gin.Context) (*dto.CartResponse, error) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return nil, apperror.BadRequest("Invalid cart item ID")
	}

	return s.cartService.RemoveCartItem(ctx.GetUint("user_id"), uint(id))
}

func (s *Server) clearCart(ctx *gin.Context) error {
	return s.cartService.ClearCart(ctx.GetUint("user_id"))
}
