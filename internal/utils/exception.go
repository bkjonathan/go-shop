package utils

import (
	"errors"
	"net/http"

	"github.com/bkjonathan/go-shop/internal/apperror"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HandleError is the single place an error coming out of a handler becomes an
// HTTP response, like a NestJS global exception filter. Anything a service does
// not describe with apperror is treated as a bug and reported as a 500.
func HandleError(ctx *gin.Context, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		ErrorResponse(ctx, appErr.Status, appErr.Message, nil)
		ctx.Abort()
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		NotFoundResponse(ctx, "Resource not found")
		ctx.Abort()
		return
	}

	// Never leak internals to clients in production; keep them while developing.
	var detail error
	if gin.Mode() != gin.ReleaseMode {
		detail = err
	}
	ErrorResponse(ctx, http.StatusInternalServerError, "Internal server error", detail)
	ctx.Abort()
}
