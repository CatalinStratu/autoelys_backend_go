package handlers

import (
	"autoelys_backend/internal/models"
	"autoelys_backend/internal/repository"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	serviceRepo *repository.ServiceRepository
}

func NewServiceHandler(serviceRepo *repository.ServiceRepository) *ServiceHandler {
	return &ServiceHandler{
		serviceRepo: serviceRepo,
	}
}

// ServiceData represents service information in responses
// @Description Service data
type ServiceData struct {
	ID              uint64  `json:"id" example:"1"`
	UUID            string  `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title           string  `json:"title" example:"Oil Change"`
	Description     *string `json:"description,omitempty" example:"Full synthetic oil change service"`
	Price           float64 `json:"price" example:"150.00"`
	Currency        string  `json:"currency" example:"lei"`
	DurationMinutes *uint   `json:"duration_minutes,omitempty" example:"30"`
	Active          bool    `json:"active" example:"true"`
	CreatedAt       string  `json:"created_at" example:"2024-01-15T10:30:00Z"`
	UpdatedAt       string  `json:"updated_at" example:"2024-01-15T10:30:00Z"`
}

// GetAllServicesResponse represents the paginated services response
// @Description Paginated services response
type GetAllServicesResponse struct {
	Data       []ServiceData  `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// CreateServiceRequest represents the create service payload
// @Description Create service request payload
type CreateServiceRequest struct {
	Title           string  `json:"title" binding:"required" example:"Oil Change"`
	Description     *string `json:"description" example:"Full synthetic oil change service"`
	Price           float64 `json:"price" binding:"required" example:"150.00"`
	Currency        string  `json:"currency" example:"lei"`
	DurationMinutes *uint   `json:"duration_minutes" example:"30"`
	Active          *bool   `json:"active" example:"true"`
}

// UpdateServiceRequest represents the update service payload
// @Description Update service request payload
type UpdateServiceRequest struct {
	Title           string  `json:"title" example:"Oil Change"`
	Description     *string `json:"description" example:"Full synthetic oil change service"`
	Price           float64 `json:"price" example:"150.00"`
	Currency        string  `json:"currency" example:"lei"`
	DurationMinutes *uint   `json:"duration_minutes" example:"30"`
	Active          *bool   `json:"active" example:"true"`
}

func serviceToData(service *models.Service) ServiceData {
	return ServiceData{
		ID:              service.ID,
		UUID:            service.UUID,
		Title:           service.Title,
		Description:     service.Description,
		Price:           service.Price,
		Currency:        service.Currency,
		DurationMinutes: service.DurationMinutes,
		Active:          service.Active,
		CreatedAt:       service.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:       service.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// GetAllServices godoc
// @Summary Get all services (Admin only)
// @Description Get a paginated list of all services
// @Tags Admin - Services
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20) maximum(100)
// @Success 200 {object} GetAllServicesResponse "List of services"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - Admin access required"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/admin/services [get]
func (h *ServiceHandler) GetAllServices(c *gin.Context) {
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

	offset := (page - 1) * limit

	services, total, err := h.serviceRepo.GetAll(limit, offset, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch services"})
		return
	}

	var serviceData []ServiceData
	for _, service := range services {
		serviceData = append(serviceData, serviceToData(&service))
	}

	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, GetAllServicesResponse{
		Data: serviceData,
		Pagination: PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// CreateService godoc
// @Summary Create a new service (Admin only)
// @Description Create a new service for car owners
// @Tags Admin - Services
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateServiceRequest true "Service details"
// @Success 201 {object} map[string]interface{} "Service created successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - Admin access required"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/admin/services [post]
func (h *ServiceHandler) CreateService(c *gin.Context) {
	var req CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	currency := req.Currency
	if currency == "" {
		currency = "lei"
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	service := &models.Service{
		Title:           req.Title,
		Description:     req.Description,
		Price:           req.Price,
		Currency:        currency,
		DurationMinutes: req.DurationMinutes,
		Active:          active,
	}

	if err := h.serviceRepo.Create(service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Service created successfully",
		"service": serviceToData(service),
	})
}

// GetService godoc
// @Summary Get a service by UUID (Admin only)
// @Description Get service details by UUID
// @Tags Admin - Services
// @Produce json
// @Security BearerAuth
// @Param uuid path string true "Service UUID"
// @Success 200 {object} ServiceData "Service details"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - Admin access required"
// @Failure 404 {object} map[string]string "Service not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/admin/services/{uuid} [get]
func (h *ServiceHandler) GetService(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service UUID is required"})
		return
	}

	service, err := h.serviceRepo.FindByUUID(uuid)
	if err != nil {
		if errors.Is(err, repository.ErrServiceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch service"})
		return
	}

	c.JSON(http.StatusOK, serviceToData(service))
}

// UpdateService godoc
// @Summary Update a service (Admin only)
// @Description Update service information by UUID
// @Tags Admin - Services
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uuid path string true "Service UUID"
// @Param request body UpdateServiceRequest true "Service update details"
// @Success 200 {object} map[string]interface{} "Service updated successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - Admin access required"
// @Failure 404 {object} map[string]string "Service not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/admin/services/{uuid} [put]
func (h *ServiceHandler) UpdateService(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service UUID is required"})
		return
	}

	service, err := h.serviceRepo.FindByUUID(uuid)
	if err != nil {
		if errors.Is(err, repository.ErrServiceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch service"})
		return
	}

	var req UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if req.Title != "" {
		service.Title = req.Title
	}
	if req.Description != nil {
		service.Description = req.Description
	}
	if req.Price != 0 {
		service.Price = req.Price
	}
	if req.Currency != "" {
		service.Currency = req.Currency
	}
	if req.DurationMinutes != nil {
		service.DurationMinutes = req.DurationMinutes
	}
	if req.Active != nil {
		service.Active = *req.Active
	}

	if err := h.serviceRepo.Update(service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Service updated successfully",
		"service": serviceToData(service),
	})
}

// DeleteService godoc
// @Summary Delete a service (Admin only)
// @Description Delete a service by UUID
// @Tags Admin - Services
// @Produce json
// @Security BearerAuth
// @Param uuid path string true "Service UUID"
// @Success 200 {object} map[string]string "Service deleted successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - Admin access required"
// @Failure 404 {object} map[string]string "Service not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/admin/services/{uuid} [delete]
func (h *ServiceHandler) DeleteService(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service UUID is required"})
		return
	}

	service, err := h.serviceRepo.FindByUUID(uuid)
	if err != nil {
		if errors.Is(err, repository.ErrServiceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch service"})
		return
	}

	if err := h.serviceRepo.Delete(service.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}

// GetPublicServices godoc
// @Summary Get all active services (Public)
// @Description Get a paginated list of all active services for car owners
// @Tags Services
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20) maximum(100)
// @Success 200 {object} GetAllServicesResponse "List of services"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/services [get]
func (h *ServiceHandler) GetPublicServices(c *gin.Context) {
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

	offset := (page - 1) * limit

	services, total, err := h.serviceRepo.GetAll(limit, offset, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch services"})
		return
	}

	var serviceData []ServiceData
	for _, service := range services {
		serviceData = append(serviceData, serviceToData(&service))
	}

	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, GetAllServicesResponse{
		Data: serviceData,
		Pagination: PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}
