package middleware

import (
	"backend/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	
	// Use JSON tag names in error messages
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validators
	registerCustomValidators()
}

// ValidationMiddleware provides centralized request validation
func ValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// ValidateStruct validates a struct and returns formatted error messages
func ValidateStruct(s interface{}) map[string]string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	
	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		tag := err.Tag()
		
		// Create user-friendly error messages
		message := getErrorMessage(field, tag, err.Param())
		errors[field] = message
	}
	
	return errors
}

// ValidateJSON validates JSON request body against a struct
func ValidateJSON(r *http.Request, dest interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}
	
	if errors := ValidateStruct(dest); errors != nil {
		return &ValidationError{Errors: errors}
	}
	
	return nil
}

// ValidationError represents validation errors
type ValidationError struct {
	Errors map[string]string `json:"errors"`
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

// HandleValidationError writes validation errors to response
func HandleValidationError(w http.ResponseWriter, err error) {
	if validationErr, ok := err.(*ValidationError); ok {
		response := map[string]interface{}{
			"success": false,
			"message": "Validation failed",
			"errors":  validationErr.Errors,
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// For other errors, use the standard error response
	utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
}

// registerCustomValidators adds custom validation rules
func registerCustomValidators() {
	// Password strength validator
	validate.RegisterValidation("strong_password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		return isStrongPassword(password)
	})
	
	// Role validator
	validate.RegisterValidation("valid_role", func(fl validator.FieldLevel) bool {
		role := fl.Field().String()
		validRoles := []string{"user", "instructor", "admin"}
		for _, validRole := range validRoles {
			if role == validRole {
				return true
			}
		}
		return false
	})
	
	// Course level validator
	validate.RegisterValidation("course_level", func(fl validator.FieldLevel) bool {
		level := fl.Field().String()
		validLevels := []string{"beginner", "intermediate", "advanced"}
		for _, validLevel := range validLevels {
			if level == validLevel {
				return true
			}
		}
		return false
	})
}

// isStrongPassword checks password strength
func isStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	
	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}
	
	return hasUpper && hasLower && hasNumber && hasSpecial
}

// getErrorMessage returns user-friendly error messages
func getErrorMessage(field, tag, param string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, param)
	case "max":
		return fmt.Sprintf("%s must be no more than %s characters long", field, param)
	case "strong_password":
		return fmt.Sprintf("%s must contain at least 8 characters including uppercase, lowercase, number and special character", field)
	case "valid_role":
		return fmt.Sprintf("%s must be one of: user, instructor, admin", field)
	case "course_level":
		return fmt.Sprintf("%s must be one of: beginner, intermediate, advanced", field)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, param)
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, param)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, param)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}