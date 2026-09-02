package server

import (
	"strconv"

	"github.com/bkjonathan/go-shop/internal/dto"
	"github.com/bkjonathan/go-shop/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) createCategory(ctx *gin.Context, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	return s.productService.CreateCategory(req)
}

func (s *Server) listCategories(ctx *gin.Context) ([]*dto.CategoryResponse, error) {
	return s.productService.ListCategories()
}

func (s *Server) updateCategory(ctx *gin.Context, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid category ID", err)
		return nil, err
	}

	return s.productService.UpdateCategory(uint(id), req)

}

func (s *Server) deleteCategory(ctx *gin.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid category ID", err)
		return err
	}
	return s.productService.DeleteCategory(uint(id))
}

func (s *Server) createProduct(ctx *gin.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	return s.productService.CreateProduct(req)
}

func (s *Server) listProducts(ctx *gin.Context) (*dto.ProductListResponse, error) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	items, meta, err := s.productService.ListProducts(page, limit)
	if err != nil {
		return nil, err
	}

	return &dto.ProductListResponse{Items: items, Meta: *meta}, nil
}

func (s *Server) getProduct(ctx *gin.Context) (*dto.ProductResponse, error) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid product ID", err)
		return nil, err
	}

	return s.productService.GetProduct(uint(id))
}

func (s *Server) updateProduct(ctx *gin.Context, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid product ID", err)
		return nil, err
	}

	return s.productService.UpdateProduct(uint(id), req)
}

func (s *Server) deleteProduct(ctx *gin.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid product ID", err)
		return err
	}
	return s.productService.DeleteProduct(uint(id))
}
