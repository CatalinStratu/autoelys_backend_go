package handlers

import (
	"net/http"
	"strconv"

	"autoelys_backend/internal/models"
	"autoelys_backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type BrandHandler struct {
	brandRepo      *repository.BrandRepository
	automobileRepo *repository.AutomobileRepository
}

func NewBrandHandler(brandRepo *repository.BrandRepository, automobileRepo *repository.AutomobileRepository) *BrandHandler {
	return &BrandHandler{
		brandRepo:      brandRepo,
		automobileRepo: automobileRepo,
	}
}

// BrandResponse represents the API response for a brand
type BrandResponse struct {
	ID   uint64 `json:"id" example:"1"`
	Name string `json:"name" example:"AUDI"`
}

// AutomobileResponse represents the API response for an automobile
type AutomobileResponse struct {
	ID   uint64 `json:"id" example:"1"`
	Name string `json:"name" example:"A4"`
}

// GetAllBrands godoc
// @Summary Get all brands
// @Description Retrieve a list of all automobile brands
// @Tags brands
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "success with brands array"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/brands [get]
func (h *BrandHandler) GetAllBrands(c *gin.Context) {
	brands, err := h.brandRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve brands",
		})
		return
	}

	// Convert to response format
	brandResponses := make([]BrandResponse, len(brands))
	for i, brand := range brands {
		brandResponses[i] = mapBrandToResponse(&brand)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    brandResponses,
	})
}

// GetAutomobilesByBrand godoc
// @Summary Get automobiles by brand
// @Description Retrieve all automobiles for a specific brand
// @Tags brands
// @Accept json
// @Produce json
// @Param id path int true "Brand ID"
// @Success 200 {object} map[string]interface{} "success with automobiles array"
// @Failure 400 {object} map[string]interface{} "invalid brand ID"
// @Failure 404 {object} map[string]interface{} "brand not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/brands/{id}/automobiles [get]
func (h *BrandHandler) GetAutomobilesByBrand(c *gin.Context) {
	// Parse brand ID from URL
	brandIDStr := c.Param("id")
	brandID, err := strconv.ParseUint(brandIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid brand ID",
		})
		return
	}

	// Check if brand exists
	brand, err := h.brandRepo.FindByID(brandID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve brand",
		})
		return
	}

	if brand == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Brand not found",
		})
		return
	}

	// Get automobiles for the brand
	automobiles, err := h.automobileRepo.GetByBrandID(brandID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve automobiles",
		})
		return
	}

	// Convert to response format
	automobileResponses := make([]AutomobileResponse, len(automobiles))
	for i, automobile := range automobiles {
		automobileResponses[i] = mapAutomobileToResponse(&automobile)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    automobileResponses,
	})
}

// Helper functions to map models to responses (DRY principle)
func mapBrandToResponse(brand *models.Brand) BrandResponse {
	response := BrandResponse{
		ID:   brand.ID,
		Name: brand.Name,
	}

	return response
}

func mapAutomobileToResponse(automobile *models.Automobile) AutomobileResponse {
	response := AutomobileResponse{
		ID:   automobile.ID,
		Name: automobile.Name,
	}

	return response
}
