package server

import (
	"strconv"

	"github.com/bkjonathan/go-shop/internal/apperror"
	"github.com/bkjonathan/go-shop/internal/dto"
	"github.com/gin-gonic/gin"
)

func (s *Server) createOrder(ctx *gin.Context) (*dto.OrderResponse, error) {
	return s.orderService.CreateOrder(ctx.GetUint("user_id"))
}

func (s *Server) getOrders(ctx *gin.Context) (*dto.OrderListResponse, error) {
	userID := ctx.GetUint("user_id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	items, meta, err := s.orderService.GetOrders(userID, page, limit)
	if err != nil {
		return nil, apperror.Internal("Failed to get orders", err)
	}

	return &dto.OrderListResponse{
		Items: items,
		Meta:  *meta,
	}, nil
}

func (s *Server) getOrder(ctx *gin.Context) (*dto.OrderResponse, error) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return nil, apperror.BadRequest("Invalid order ID")
	}

	return s.orderService.GetOrder(ctx.GetUint("user_id"), uint(id))
}
