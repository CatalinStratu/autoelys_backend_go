package handlers

import (
	"autoelys_backend/internal/auth"
	"autoelys_backend/internal/models"
	"autoelys_backend/internal/repository"
	"autoelys_backend/internal/services"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	userRepo     *repository.UserRepository
	passwordRepo *repository.PasswordResetRepository
	emailService *services.EmailService
	validator    *validator.Validate
}

func NewAuthHandler(
	userRepo *repository.UserRepository,
	passwordRepo *repository.PasswordResetRepository,
	emailService *services.EmailService,
	validate *validator.Validate,
) *AuthHandler {
	return &AuthHandler{
		userRepo:     userRepo,
		passwordRepo: passwordRepo,
		emailService: emailService,
		validator:    validate,
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
	UUID      string `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
	Email     string `json:"email" example:"john@example.com"`
}

// ErrorResponse represents an error response
// @Description Error response
type ErrorResponse struct {
	Errors map[string][]string `json:"errors"`
}

// LoginRequest represents the login payload
// @Description User login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required" example:"Password123"`
}

// LoginResponse represents the successful login response
// @Description User login response
type LoginResponse struct {
	Message string   `json:"message" example:"Login successful."`
	User    UserData `json:"user"`
	Token   string   `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// ForgotPasswordRequest represents the forgot password payload
// @Description Forgot password request payload
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email" example:"john@example.com"`
}

// ResetPasswordRequest represents the reset password payload
// @Description Reset password request payload
type ResetPasswordRequest struct {
	Token                string `json:"token" validate:"required" example:"a1b2c3d4e5f6..."`
	Password             string `json:"password" validate:"required,strong_password" example:"NewPassword123"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password" example:"NewPassword123"`
}

// UserProfileResponse represents the user profile response
// @Description User profile response
type UserProfileResponse struct {
	User UserData `json:"user"`
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
			UUID:      user.UUID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
		Token: token,
	}

	c.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse "Login successful"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 422 {object} ErrorResponse "Validation error"
// @Failure 429 {object} map[string]string "Too many requests"
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		errors := formatValidationErrors(err.(validator.ValidationErrors))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errors})
		return
	}

	user, err := h.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if !user.Active {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is inactive"})
		return
	}

	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email, user.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	response := LoginResponse{
		Message: "Login successful.",
		User: UserData{
			UUID:      user.UUID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
		Token: token,
	}

	c.JSON(http.StatusOK, response)
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send password reset email to user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "Email address"
// @Success 200 {object} map[string]string "Reset email sent"
// @Failure 422 {object} ErrorResponse "Validation error"
// @Failure 429 {object} map[string]string "Too many requests"
// @Router /api/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		validationErrors := formatValidationErrors(err.(validator.ValidationErrors))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": validationErrors})
		return
	}

	user, err := h.userRepo.FindByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "If the email exists, a password reset link has been sent.",
		})
		return
	}

	if !user.Active {
		c.JSON(http.StatusOK, gin.H{
			"message": "If the email exists, a password reset link has been sent.",
		})
		return
	}

	_ = h.passwordRepo.DeleteUserTokens(user.ID)

	resetToken, err := h.passwordRepo.Create(user.ID, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if err := h.emailService.SendPasswordResetEmail(user.Email, resetToken.Token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "If the email exists, a password reset link has been sent.",
	})
}

// ResetPassword godoc
// @Summary Reset password with token
// @Description Reset user password using reset token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Reset token and new password"
// @Success 200 {object} map[string]string "Password reset successful"
// @Failure 400 {object} map[string]string "Invalid or expired token"
// @Failure 422 {object} ErrorResponse "Validation error"
// @Router /api/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		validationErrors := formatValidationErrors(err.(validator.ValidationErrors))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": validationErrors})
		return
	}

	resetToken, err := h.passwordRepo.ValidateToken(req.Token)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidToken) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset token"})
			return
		}
		if errors.Is(err, repository.ErrTokenAlreadyUsed) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Reset token has already been used"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	user, err := h.userRepo.FindByID(resetToken.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if err := h.userRepo.UpdatePassword(user.ID, hashedPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	if err := h.passwordRepo.MarkAsUsed(resetToken.ID); err != nil {
		// Log but don't fail
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password has been reset successfully. You can now login with your new password.",
	})
}

// GetMe godoc
// @Summary Get current user profile
// @Description Get the authenticated user's profile information
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserProfileResponse "User profile"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "User not found"
// @Router /api/auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.userRepo.FindByID(userID.(uint64))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	response := UserProfileResponse{
		User: UserData{
			UUID:      user.UUID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
	}

	c.JSON(http.StatusOK, response)
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
