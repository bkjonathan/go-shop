package server

import (
	"errors"
	"net/http"

	"github.com/bkjonathan/go-shop/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Handle turns a typed handler into a gin.HandlerFunc. It plays the role of
// NestJS's ValidationPipe plus the response interceptor: bind the request into
// Req, validate it, run the handler, then render either the data envelope with
// `status` and `message` or whatever the error maps to.
//
//	auth.POST("/login", Handle(http.StatusOK, "Login successful", s.login))
//
// A handler therefore only ever contains business calls:
//
//	func (s *Server) login(ctx *gin.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
//		return s.authService.Login(req)
//	}
func Handle[Req any, Res any](status int, message string, handler func(*gin.Context, *Req) (Res, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := bindRequest(ctx, &req); err != nil {
			utils.ValidationErrorResponse(ctx, err)
			return
		}

		res, err := handler(ctx, &req)
		if err != nil {
			utils.HandleError(ctx, err)
			return
		}

		utils.DataResponse(ctx, status, message, res)
	}
}

// HandleEmpty is Handle for routes that take no request payload.
func HandleEmpty[Res any](status int, message string, handler func(*gin.Context) (Res, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := handler(ctx)
		if err != nil {
			utils.HandleError(ctx, err)
			return
		}

		utils.DataResponse(ctx, status, message, res)
	}
}

// bindRequest fills req from the path params, the query string and the body, so
// one DTO can carry `uri`, `form` and `json` tagged fields at once. Validation
// errors raised by the intermediate binds are ignored on purpose: the struct is
// only complete after the last source, so it is validated once at the end.
func bindRequest(ctx *gin.Context, req any) error {
	if len(ctx.Params) > 0 {
		if err := skipValidationErrors(ctx.ShouldBindUri(req)); err != nil {
			return err
		}
	}

	if ctx.Request.URL.RawQuery != "" {
		if err := skipValidationErrors(ctx.ShouldBindQuery(req)); err != nil {
			return err
		}
	}

	if hasBody(ctx.Request) {
		if err := skipValidationErrors(ctx.ShouldBind(req)); err != nil {
			return err
		}
	}

	return binding.Validator.ValidateStruct(req)
}

func skipValidationErrors(err error) error {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		return nil
	}
	return err
}

func hasBody(req *http.Request) bool {
	if req.Body == nil || req.ContentLength == 0 {
		return false
	}

	switch req.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}
