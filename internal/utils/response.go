package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    interface{}  `json:"data"`
	Error   string       `json:"error"`
	Errors  []FieldError `json:"errors,omitempty"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Response
	Meta PaginationMeta `json:"meta"`
}

// DataResponse writes a successful envelope with an explicit status code.
// SuccessResponse and CreatedResponse are the two common cases.
func DataResponse(ctx *gin.Context, statusCode int, message string, data interface{}) {
	ctx.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessResponse(ctx *gin.Context, message string, data interface{}) {
	DataResponse(ctx, http.StatusOK, message, data)
}

func CreatedResponse(ctx *gin.Context, message string, data interface{}) {
	DataResponse(ctx, http.StatusCreated, message, data)
}

func ErrorResponse(ctx *gin.Context, statusCode int, message string, err error) {
	response := Response{
		Success: false,
		Message: message,
	}
	if err != nil {
		response.Error = err.Error()
	}
	ctx.JSON(statusCode, response)
}

func BadRequestResponse(ctx *gin.Context, message string, err error) {
	ErrorResponse(ctx, http.StatusBadRequest, message, err)
}

func UnauthorizedResponse(ctx *gin.Context, message string) {
	ErrorResponse(ctx, http.StatusUnauthorized, message, nil)
}

func ForbiddenResponse(ctx *gin.Context, message string) {
	ErrorResponse(ctx, http.StatusForbidden, message, nil)
}

func NotFoundResponse(ctx *gin.Context, message string) {
	ErrorResponse(ctx, http.StatusNotFound, message, nil)
}

func InternalSServerErrorResponse(ctx *gin.Context, message string, err error) {
	ErrorResponse(ctx, http.StatusInternalServerError, message, err)
}

func PaginatedSuccessResponse(ctx *gin.Context, message string, data interface{}, meta PaginationMeta) {
	ctx.JSON(http.StatusOK, PaginatedResponse{
		Response: Response{
			Success: true,
			Message: message,
			Data:    data,
		},
		Meta: meta,
	})
}

// ValidationErrorResponse reports a request that failed binding or validation,
// listing every offending field at once instead of just the first.
func ValidationErrorResponse(ctx *gin.Context, err error) {
	response := Response{
		Success: false,
		Message: "Validation failed",
		Errors:  ValidationFieldErrors(err),
	}

	if len(response.Errors) == 0 {
		// Not a field-level problem: malformed JSON, wrong content type, ...
		response.Message = "Invalid request data"
		response.Error = err.Error()
	}

	ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
}
