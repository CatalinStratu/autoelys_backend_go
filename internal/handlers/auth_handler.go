package handlers

import (
	"autoelys_backend/internal/auth"
	"autoelys_backend/internal/models"
	"autoelys_backend/internal/repository"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	userRepo  *repository.UserRepository
	validator *validator.Validate
}

func NewAuthHandler(userRepo *repository.UserRepository, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{
		userRepo:  userRepo,
		validator: validate,
	}
}

// RegisterRequest represents the registration payload
// @Description User registration request payload
type RegisterRequest struct {
	FirstName            string `json:"first_name" validate:"required,min=2" example:"John"`
	LastName             string `json:"last_name" validate:"required,min=2" example:"Doe"`
	Email                string `json:"email" validate:"required,email" example:"john@example.com"`
	Phone                string `json:"phone" validate:"omitempty,phone_e164" example:"+40712345678"`
	Password             string `json:"password" validate:"required,strong_password" example:"Password123"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password" example:"Password123"`
	AcceptedTerms        bool   `json:"accepted_terms" validate:"required,eq=true" example:"true"`
}

// RegisterResponse represents the successful registration response
// @Description User registration response
type RegisterResponse struct {
	Message string   `json:"message" example:"Account created successfully."`
	User    UserData `json:"user"`
	Token   string   `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// UserData represents user information in responses
// @Description User data
type UserData struct {
	ID        uint64 `json:"id" example:"1"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
	Email     string `json:"email" example:"john@example.com"`
}

// ErrorResponse represents an error response
// @Description Error response
type ErrorResponse struct {
	Errors map[string][]string `json:"errors"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} RegisterResponse "Account created successfully"
// @Failure 422 {object} ErrorResponse "Validation error"
// @Failure 429 {object} map[string]string "Too many requests"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		errors := formatValidationErrors(err.(validator.ValidationErrors))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errors})
		return
	}

	emailExists, err := h.userRepo.EmailExists(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if emailExists {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": map[string][]string{
				"email": {"The email is already taken."},
			},
		})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	userRoleID, err := h.userRepo.GetRoleIDByName(models.RoleUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	now := time.Now()
	user := &models.User{
		RoleID:          userRoleID,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Email:           req.Email,
		Phone:           stringPtr(req.Phone),
		PasswordHash:    hashedPassword,
		Active:          true,
		AcceptedTermsAt: &now,
	}

	if err := h.userRepo.Create(user); err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"errors": map[string][]string{
					"email": {"The email is already taken."},
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email, user.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	response := RegisterResponse{
		Message: "Account created successfully.",
		User: UserData{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
		Token: token,
	}

	c.JSON(http.StatusCreated, response)
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func formatValidationErrors(errs validator.ValidationErrors) map[string][]string {
	errors := make(map[string][]string)

	for _, err := range errs {
		fieldName := toSnakeCase(err.Field())

		var message string
		switch err.Tag() {
		case "required":
			message = "This field is required."
		case "min":
			message = "Must be at least " + err.Param() + " characters."
		case "email":
			message = "Must be a valid email address."
		case "phone_e164":
			message = "Must be a valid phone number in E.164 format (e.g., +407xxxxxxxx)."
		case "strong_password":
			message = "Must be at least 8 characters long and contain both letters and digits."
		case "eqfield":
			message = "Passwords do not match."
		case "eq":
			if err.Field() == "AcceptedTerms" {
				message = "You must accept the terms and conditions."
			} else {
				message = "Invalid value."
			}
		default:
			message = "Invalid value for " + fieldName + "."
		}

		errors[fieldName] = append(errors[fieldName], message)
	}

	return errors
}

func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
