package handlers

import (
	"autoelys_backend/internal/middleware"
	"autoelys_backend/internal/repository"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	userRepo *repository.UserRepository
}

func NewAdminHandler(userRepo *repository.UserRepository) *AdminHandler {
	return &AdminHandler{
		userRepo: userRepo,
	}
}

// AdminUserData represents user information in admin responses
// @Description Admin user data
type AdminUserData struct {
	ID        uint64  `json:"id" example:"1"`
	UUID      string  `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	RoleID    uint64  `json:"role_id" example:"2"`
	FirstName string  `json:"first_name" example:"John"`
	LastName  string  `json:"last_name" example:"Doe"`
	Email     string  `json:"email" example:"john@example.com"`
	Phone     *string `json:"phone,omitempty" example:"+40712345678"`
	Active    bool    `json:"active" example:"true"`
	CreatedAt string  `json:"created_at" example:"2024-01-15T10:30:00Z"`
	UpdatedAt string  `json:"updated_at" example:"2024-01-15T10:30:00Z"`
}

// GetAllUsersResponse represents the paginated users response
// @Description Paginated users response
type GetAllUsersResponse struct {
	Data       []AdminUserData `json:"data"`
	Pagination PaginationMeta  `json:"pagination"`
}

// PaginationMeta represents pagination metadata
// @Description Pagination metadata
type PaginationMeta struct {
	Page       int `json:"page" example:"1"`
	Limit      int `json:"limit" example:"20"`
	Total      int `json:"total" example:"100"`
	TotalPages int `json:"total_pages" example:"5"`
}

// GetAllUsers godoc
// @Summary Get all users (Admin only)
// @Description Get a paginated list of all users on the platform
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20) maximum(100)
// @Param search query string false "Search by name or email"
// @Success 200 {object} GetAllUsersResponse "List of users"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - Admin access required"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/admin/users [get]
func (h *AdminHandler) GetAllUsers(c *gin.Context) {
	page := 1
	if pageStr := c.DefaultQuery("page", "1"); pageStr != "" {
		if val, err := strconv.Atoi(pageStr); err == nil && val > 0 {
			page = val
		}
	}

	limit := 20
	if limitStr := c.DefaultQuery("limit", "20"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			limit = val
			if limit > 100 {
				limit = 100
			}
		}
	}

	search := c.Query("search")
	offset := (page - 1) * limit

	users, total, err := h.userRepo.GetAllUsers(limit, offset, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	var userData []AdminUserData
	for _, user := range users {
		userData = append(userData, AdminUserData{
			ID:        user.ID,
			UUID:      user.UUID,
			RoleID:    user.RoleID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Phone:     user.Phone,
			Active:    user.Active,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, GetAllUsersResponse{
		Data: userData,
		Pagination: PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// AdminUpdateUserRequest represents the update user payload for admin
// @Description Admin update user request payload
type AdminUpdateUserRequest struct {
	FirstName string  `json:"first_name" example:"John"`
	LastName  string  `json:"last_name" example:"Doe"`
	Email     string  `json:"email" example:"john@example.com"`
	Phone     *string `json:"phone" example:"+40712345678"`
	RoleID    *uint64 `json:"role_id" example:"2"`
	Active    *bool   `json:"active" example:"true"`
}

// UpdateUser godoc
// @Summary Update a user (Admin only)
// @Description Update user information by UUID
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uuid path string true "User UUID"
// @Param request body AdminUpdateUserRequest true "User update details"
// @Success 200 {object} map[string]interface{} "User updated successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 422 {object} map[string]string "Email already taken"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/admin/users/{uuid} [put]
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User UUID is required"})
		return
	}

	user, err := h.userRepo.FindByUUID(uuid)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// Prevent admin from modifying themselves (optional safety)
	adminID, _ := c.Get("user_id")
	if user.ID == adminID.(uint64) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot modify your own account via admin panel"})
		return
	}

	var req AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.RoleID != nil {
		user.RoleID = *req.RoleID
	}
	if req.Active != nil {
		user.Active = *req.Active
	}

	if err := h.userRepo.AdminUpdateUser(user); err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Email already taken"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user": AdminUserData{
			ID:        user.ID,
			UUID:      user.UUID,
			RoleID:    user.RoleID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Phone:     user.Phone,
			Active:    user.Active,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		},
	})
}

// DeleteUser godoc
// @Summary Delete a user (Admin only)
// @Description Delete a user by UUID
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param uuid path string true "User UUID"
// @Success 200 {object} map[string]string "User deleted successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/admin/users/{uuid} [delete]
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User UUID is required"})
		return
	}

	user, err := h.userRepo.FindByUUID(uuid)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// Prevent admin from deleting themselves
	adminID, _ := c.Get("user_id")
	if user.ID == adminID.(uint64) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete your own account"})
		return
	}

	// Prevent deleting other admins
	if user.RoleID == middleware.AdminRoleID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete admin users"})
		return
	}

	if err := h.userRepo.DeleteUser(user.ID); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
