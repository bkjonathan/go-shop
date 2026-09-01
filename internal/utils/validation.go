package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// FieldError is one failed rule on one field, shaped for the JSON response.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// RegisterValidationTagNames makes validation errors report the wire name
// ("first_name") instead of the Go field name ("FirstName"). Call it once at
// startup, before any route is served.
func RegisterValidationTagNames() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		// Whichever source the field is bound from: body, query or path.
		for _, tag := range []string{"json", "form", "uri"} {
			name := strings.SplitN(fld.Tag.Get(tag), ",", 2)[0]
			if name != "" && name != "-" {
				return name
			}
		}
		return fld.Name
	})
}

// ValidationFieldErrors turns a binding failure into per-field messages.
// It returns nil for errors that are not about a specific field (malformed
// JSON, for instance), so the caller can fall back to a plain message.
func ValidationFieldErrors(err error) []FieldError {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		fields := make([]FieldError, 0, len(validationErrs))
		for _, fieldErr := range validationErrs {
			fields = append(fields, FieldError{
				Field:   fieldErr.Field(),
				Message: validationMessage(fieldErr),
			})
		}
		return fields
	}

	// A value of the wrong JSON type ({"stock": "ten"}) never reaches the
	// validator, but it is still a field-level problem.
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) && typeErr.Field != "" {
		return []FieldError{{
			Field:   typeErr.Field,
			Message: fmt.Sprintf("%s must be of type %s", typeErr.Field, typeErr.Type),
		}}
	}

	return nil
}

func validationMessage(fieldErr validator.FieldError) string {
	field := fieldErr.Field()
	param := fieldErr.Param()

	switch fieldErr.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "url":
		return field + " must be a valid URL"
	case "uuid", "uuid4":
		return field + " must be a valid UUID"
	case "numeric":
		return field + " must be numeric"
	case "alphanum":
		return field + " must contain only letters and numbers"
	case "min":
		return fmt.Sprintf("%s must be at least %s%s", field, param, sizeUnit(fieldErr))
	case "max":
		return fmt.Sprintf("%s must be at most %s%s", field, param, sizeUnit(fieldErr))
	case "len":
		return fmt.Sprintf("%s must be exactly %s%s", field, param, sizeUnit(fieldErr))
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case "gte":
		return fmt.Sprintf("%s must be %s or greater", field, param)
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, param)
	case "lte":
		return fmt.Sprintf("%s must be %s or less", field, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, strings.Join(strings.Fields(param), ", "))
	case "eqfield":
		return fmt.Sprintf("%s must match %s", field, param)
	default:
		return fmt.Sprintf("%s failed the %q rule", field, fieldErr.Tag())
	}
}

// sizeUnit spells out what min/max/len are counting for this field's type.
func sizeUnit(fieldErr validator.FieldError) string {
	switch fieldErr.Kind() {
	case reflect.String:
		return " characters"
	case reflect.Slice, reflect.Array, reflect.Map:
		return " items"
	default:
		return ""
	}
}
