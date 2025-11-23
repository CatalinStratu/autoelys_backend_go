package handlers

import (
	"net/http"
	"strconv"
	"time"

	"autoelys_backend/internal/models"
	"autoelys_backend/internal/repository"
	"autoelys_backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type VehicleHandler struct {
	vehicleRepo *repository.VehicleRepository
	validator   *validator.Validate
}

func NewVehicleHandler(vehicleRepo *repository.VehicleRepository, validator *validator.Validate) *VehicleHandler {
	return &VehicleHandler{
		vehicleRepo: vehicleRepo,
		validator:   validator,
	}
}

// CreateVehicleRequest represents the vehicle creation request
type CreateVehicleRequest struct {
	Title          string  `form:"title" validate:"required,min=5,max=255"`
	Category       string  `form:"category" validate:"required"`
	Description    string  `form:"description"`
	Price          float64 `form:"price" validate:"required,gt=0"`
	Currency       string  `form:"currency" validate:"required"`
	Negotiable     bool    `form:"negotiable"`
	PersonType     string  `form:"person_type" validate:"required,oneof=persoana_fizica firma"`
	Brand          string  `form:"brand" validate:"required"`
	Model          string  `form:"model" validate:"required"`
	EngineCapacity int     `form:"engine_capacity"`
	PowerHP        int     `form:"power_hp"`
	FuelType       string  `form:"fuel_type" validate:"required,oneof=benzina motorina electric hibrid gpl hybrid_benzina hybrid_motorina"`
	BodyType       string  `form:"body_type" validate:"required,oneof=sedan suv break coupe cabrio hatchback pickup van monovolum"`
	Kilometers     int     `form:"kilometers"`
	Color          string  `form:"color"`
	Year           int     `form:"year" validate:"required,min=1970,max=2030"`
	NumberOfKeys   int     `form:"number_of_keys"`
	Condition      string  `form:"condition" validate:"required,oneof=utilizat nou"`
	Transmission   string  `form:"transmission" validate:"required,oneof=manuala automata"`
	Steering       string  `form:"steering" validate:"required,oneof=stanga dreapta"`
	Registered     bool    `form:"registered"`
	City           string  `form:"city" validate:"required"`
	ContactName    string  `form:"contact_name" validate:"required"`
	Email          string  `form:"email" validate:"required,email"`
	Phone          string  `form:"phone"`
}

// CreateVehicle godoc
// @Summary Create a new vehicle listing (Authenticated users only)
// @Description Add a new vehicle with images and all details. Only authenticated users can create vehicles. The vehicle will be automatically assigned to the authenticated user.
// @Tags vehicles
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Vehicle title (min 5, max 255 characters)"
// @Param category formData string true "Vehicle category"
// @Param description formData string false "Vehicle description"
// @Param price formData number true "Price (must be greater than 0)"
// @Param currency formData string true "Currency (e.g., lei)"
// @Param negotiable formData boolean false "Price negotiable (default: false)"
// @Param person_type formData string true "Person type (persoana_fizica, firma)"
// @Param brand formData string true "Brand"
// @Param model formData string true "Model"
// @Param engine_capacity formData integer false "Engine capacity in cm3"
// @Param power_hp formData integer false "Power in HP"
// @Param fuel_type formData string true "Fuel type (benzina, motorina, electric, hibrid, gpl, hybrid_benzina, hybrid_motorina)"
// @Param body_type formData string true "Body type (sedan, suv, break, coupe, cabrio, hatchback, pickup, van, monovolum)"
// @Param kilometers formData integer false "Kilometers"
// @Param color formData string false "Color"
// @Param year formData integer true "Year (1970-2030)"
// @Param number_of_keys formData integer false "Number of keys"
// @Param condition formData string true "Condition (utilizat, nou)"
// @Param transmission formData string true "Transmission (manuala, automata)"
// @Param steering formData string true "Steering position (stanga, dreapta)"
// @Param registered formData boolean false "Vehicle registered (default: false)"
// @Param city formData string true "City"
// @Param contact_name formData string true "Contact name"
// @Param email formData string true "Email address (valid email format)"
// @Param phone formData string false "Phone number"
// @Param images formData file false "Vehicle images (max 8, jpeg/png/jpg)"
// @Success 201 {object} map[string]interface{} "Vehicle created successfully"
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized - Authentication required"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/user/vehicles [post]
// @Security BearerAuth
func (h *VehicleHandler) CreateVehicle(c *gin.Context) {
	var req CreateVehicleRequest

	// Bind form data
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Validation failed",
			"errors":  utils.FormatValidationErrorsSimple(err.(validator.ValidationErrors)),
		})
		return
	}

	// Additional validation: year cannot be in the future
	currentYear := time.Now().Year()
	if req.Year > currentYear+1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Year cannot be more than one year in the future",
		})
		return
	}

	// Handle image uploads
	form, err := c.MultipartForm()
	var imagePaths []string
	if err == nil && form != nil && form.File["images"] != nil {
		files := form.File["images"]
		imagePaths, err = utils.UploadVehicleImages(files, "./uploads/vehicles")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Failed to upload images",
				"error":   err.Error(),
			})
			return
		}
	}

	// Look up IDs from reference tables
	personTypeID, err := h.vehicleRepo.GetPersonTypeID(req.PersonType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid person_type value",
		})
		return
	}

	fuelTypeID, err := h.vehicleRepo.GetFuelTypeID(req.FuelType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid fuel_type value",
		})
		return
	}

	bodyTypeID, err := h.vehicleRepo.GetBodyTypeID(req.BodyType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid body_type value",
		})
		return
	}

	conditionID, err := h.vehicleRepo.GetConditionID(req.Condition)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid condition value",
		})
		return
	}

	transmissionID, err := h.vehicleRepo.GetTransmissionID(req.Transmission)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid transmission value",
		})
		return
	}

	steeringID, err := h.vehicleRepo.GetSteeringID(req.Steering)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid steering value",
		})
		return
	}

	// Get authenticated user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User not authenticated",
		})
		return
	}

	// Generate UUID and slug
	vehicleUUID := uuid.New().String()
	slug := utils.GenerateSlug(req.Title)

	// Create vehicle model
	vehicle := &models.Vehicle{
		UserID:         userID.(uint64),
		UUID:           vehicleUUID,
		Slug:           slug,
		Title:          req.Title,
		Category:       req.Category,
		Price:          req.Price,
		Currency:       req.Currency,
		Negotiable:     req.Negotiable,
		PersonTypeID:   personTypeID,
		Brand:          req.Brand,
		Model:          req.Model,
		FuelTypeID:     fuelTypeID,
		BodyTypeID:     bodyTypeID,
		Year:           req.Year,
		ConditionID:    conditionID,
		TransmissionID: transmissionID,
		SteeringID:     steeringID,
		Registered:     req.Registered,
		City:           req.City,
		ContactName:    req.ContactName,
		Email:          req.Email,
	}

	// Set optional string fields
	if req.Description != "" {
		vehicle.Description = &req.Description
	}
	if req.Color != "" {
		vehicle.Color = &req.Color
	}
	if req.Phone != "" {
		vehicle.Phone = &req.Phone
	}

	// Set optional int fields
	if req.EngineCapacity > 0 {
		vehicle.EngineCapacity = &req.EngineCapacity
	}
	if req.PowerHP > 0 {
		vehicle.PowerHP = &req.PowerHP
	}
	if req.Kilometers > 0 {
		vehicle.Kilometers = &req.Kilometers
	}
	if req.NumberOfKeys > 0 {
		vehicle.NumberOfKeys = &req.NumberOfKeys
	}

	// Save vehicle to database
	createdVehicle, err := h.vehicleRepo.Create(vehicle)
	if err != nil {
		// Cleanup uploaded images on failure
		for _, path := range imagePaths {
			_ = utils.DeleteFile(path)
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create vehicle",
			"error":   err.Error(),
		})
		return
	}

	// Save images to database
	for _, imagePath := range imagePaths {
		if err := h.vehicleRepo.CreateImage(createdVehicle.ID, imagePath); err != nil {
			// Log error but continue - vehicle is already created
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Vehicle created but failed to save some images",
				"error":   err.Error(),
			})
			return
		}
	}

	// Fetch complete vehicle with images
	completeVehicle, err := h.vehicleRepo.GetByID(createdVehicle.ID)
	if err != nil {
		completeVehicle = createdVehicle // Fallback to basic vehicle data
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":     "success",
		"message":    "Vehicle added successfully",
		"vehicle_id": createdVehicle.ID,
		"data":       completeVehicle,
	})
}

// GetUserVehicles godoc
// @Summary Get all vehicles for authenticated user
// @Description Retrieve all vehicles created by the authenticated user
// @Tags vehicles
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of user's vehicles"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/user/vehicles [get]
// @Security BearerAuth
func (h *VehicleHandler) GetUserVehicles(c *gin.Context) {
	// Get authenticated user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User not authenticated",
		})
		return
	}

	vehicles, err := h.vehicleRepo.GetByUserID(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve vehicles",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"count":  len(vehicles),
		"data":   vehicles,
	})
}

// GetAllVehicles godoc
// @Summary Get all vehicles with search and filters (Public)
// @Description Public endpoint to retrieve all active vehicles with optional search and filtering. No authentication required. Perfect for browsing and searching the vehicle marketplace.
// @Tags vehicles
// @Accept json
// @Produce json
// @Param search query string false "Search by title, brand, model, or description"
// @Param brand query string false "Filter by brand name"
// @Param model query string false "Filter by model name"
// @Param fuel_type query string false "Filter by fuel type (benzina, motorina, electric, hibrid, gpl, hybrid_benzina, hybrid_motorina)"
// @Param body_type query string false "Filter by body type (sedan, suv, break, coupe, cabrio, hatchback, pickup, van, monovolum)"
// @Param transmission query string false "Filter by transmission (manuala, automata)"
// @Param condition query string false "Filter by condition (utilizat, nou)"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param min_year query int false "Minimum year"
// @Param max_year query int false "Maximum year"
// @Param city query string false "Filter by city"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20, max: 100)"
// @Success 200 {object} map[string]interface{} "List of vehicles with pagination"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/vehicles [get]
func (h *VehicleHandler) GetAllVehicles(c *gin.Context) {
	// Parse query parameters
	search := c.DefaultQuery("search", "")
	brand := c.DefaultQuery("brand", "")
	model := c.DefaultQuery("model", "")
	fuelType := c.DefaultQuery("fuel_type", "")
	bodyType := c.DefaultQuery("body_type", "")
	transmission := c.DefaultQuery("transmission", "")
	condition := c.DefaultQuery("condition", "")
	city := c.DefaultQuery("city", "")

	var minPrice, maxPrice float64
	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		if val, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			minPrice = val
		}
	}
	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		if val, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			maxPrice = val
		}
	}

	var minYear, maxYear int
	if minYearStr := c.Query("min_year"); minYearStr != "" {
		if val, err := strconv.Atoi(minYearStr); err == nil {
			minYear = val
		}
	}
	if maxYearStr := c.Query("max_year"); maxYearStr != "" {
		if val, err := strconv.Atoi(maxYearStr); err == nil {
			maxYear = val
		}
	}

	// Parse pagination parameters
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
				limit = 100 // Max limit
			}
		}
	}

	offset := (page - 1) * limit

	// Build search parameters
	params := repository.VehicleSearchParams{
		Search:       search,
		Brand:        brand,
		Model:        model,
		FuelType:     fuelType,
		BodyType:     bodyType,
		Transmission: transmission,
		Condition:    condition,
		MinPrice:     minPrice,
		MaxPrice:     maxPrice,
		MinYear:      minYear,
		MaxYear:      maxYear,
		City:         city,
		Limit:        limit,
		Offset:       offset,
	}

	// Get vehicles from repository
	vehicles, total, err := h.vehicleRepo.GetAll(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve vehicles",
			"error":   err.Error(),
		})
		return
	}

	// Calculate pagination info
	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   vehicles,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetVehicle godoc
// @Summary Get vehicle by slug (Public)
// @Description Public endpoint to retrieve detailed information about a specific vehicle using its SEO-friendly slug. Returns all vehicle details, images, and specifications. No authentication required.
// @Tags vehicles
// @Accept json
// @Produce json
// @Param slug path string true "Vehicle slug (SEO-friendly URL identifier)"
// @Success 200 {object} map[string]interface{} "Vehicle details with complete information"
// @Failure 404 {object} map[string]interface{} "Vehicle not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/vehicles/{slug} [get]
func (h *VehicleHandler) GetVehicle(c *gin.Context) {
	slug := c.Param("slug")

	vehicle, err := h.vehicleRepo.GetBySlug(slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve vehicle",
			"error":   err.Error(),
		})
		return
	}

	if vehicle == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Vehicle not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   vehicle,
	})
}

// GetVehicleByUUID godoc
// @Summary Get vehicle by UUID (Owner/Admin only)
// @Description Retrieve a vehicle with all its details and images using UUID. Only accessible by vehicle owner or admin.
// @Tags vehicles
// @Accept json
// @Produce json
// @Param uuid path string true "Vehicle UUID"
// @Success 200 {object} map[string]interface{} "Vehicle details"
// @Failure 403 {object} map[string]interface{} "Forbidden - Not owner or admin"
// @Failure 404 {object} map[string]interface{} "Vehicle not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/user/vehicles/{uuid} [get]
// @Security BearerAuth
func (h *VehicleHandler) GetVehicleByUUID(c *gin.Context) {
	vehicleUUID := c.Param("uuid")

	// Get authenticated user info from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User not authenticated",
		})
		return
	}

	roleID, exists := c.Get("role_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User role not found",
		})
		return
	}

	vehicle, err := h.vehicleRepo.GetByUUID(vehicleUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve vehicle",
			"error":   err.Error(),
		})
		return
	}

	if vehicle == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Vehicle not found",
		})
		return
	}

	// Authorization check: Only owner or admin can view
	isOwner := vehicle.UserID == userID.(uint64)
	isAdmin := roleID.(uint64) == 1 // Admin role ID is 1

	if !isOwner && !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "You don't have permission to view this vehicle",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   vehicle,
	})
}

// UpdateVehicleRequest represents the vehicle update request
type UpdateVehicleRequest struct {
	Title          string  `form:"title" validate:"omitempty,min=5,max=255"`
	Category       string  `form:"category"`
	Description    string  `form:"description"`
	Price          float64 `form:"price" validate:"omitempty,gt=0"`
	Currency       string  `form:"currency"`
	Negotiable     bool    `form:"negotiable"`
	PersonType     string  `form:"person_type" validate:"omitempty,oneof=persoana_fizica firma"`
	Brand          string  `form:"brand"`
	Model          string  `form:"model"`
	EngineCapacity int     `form:"engine_capacity"`
	PowerHP        int     `form:"power_hp"`
	FuelType       string  `form:"fuel_type" validate:"omitempty,oneof=benzina motorina electric hibrid gpl hybrid_benzina hybrid_motorina"`
	BodyType       string  `form:"body_type" validate:"omitempty,oneof=sedan suv break coupe cabrio hatchback pickup van monovolum"`
	Kilometers     int     `form:"kilometers"`
	Color          string  `form:"color"`
	Year           int     `form:"year" validate:"omitempty,min=1970,max=2030"`
	NumberOfKeys   int     `form:"number_of_keys"`
	Condition      string  `form:"condition" validate:"omitempty,oneof=utilizat nou"`
	Transmission   string  `form:"transmission" validate:"omitempty,oneof=manuala automata"`
	Steering       string  `form:"steering" validate:"omitempty,oneof=stanga dreapta"`
	Registered     bool    `form:"registered"`
	City           string  `form:"city"`
	ContactName    string  `form:"contact_name"`
	Email          string  `form:"email" validate:"omitempty,email"`
	Phone          string  `form:"phone"`
}

// UpdateVehicle godoc
// @Summary Update vehicle by UUID
// @Description Update a vehicle listing using UUID
// @Tags vehicles
// @Accept multipart/form-data
// @Produce json
// @Param uuid path string true "Vehicle UUID"
// @Param title formData string false "Vehicle title"
// @Param category formData string false "Vehicle category"
// @Param description formData string false "Vehicle description"
// @Param price formData number false "Price"
// @Param currency formData string false "Currency (e.g., lei)"
// @Param negotiable formData boolean false "Price negotiable"
// @Param person_type formData string false "Person type (persoana_fizica, firma)"
// @Param brand formData string false "Brand"
// @Param model formData string false "Model"
// @Param engine_capacity formData integer false "Engine capacity (cm3)"
// @Param power_hp formData integer false "Power (HP)"
// @Param fuel_type formData string false "Fuel type"
// @Param body_type formData string false "Body type"
// @Param kilometers formData integer false "Kilometers"
// @Param color formData string false "Color"
// @Param year formData integer false "Year"
// @Param number_of_keys formData integer false "Number of keys"
// @Param condition formData string false "Condition (utilizat, nou)"
// @Param transmission formData string false "Transmission (manuala, automata)"
// @Param steering formData string false "Steering (stanga, dreapta)"
// @Param registered formData boolean false "Registered"
// @Param city formData string false "City"
// @Param contact_name formData string false "Contact name"
// @Param email formData string false "Email"
// @Param phone formData string false "Phone"
// @Success 200 {object} map[string]interface{} "Vehicle updated successfully"
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 403 {object} map[string]interface{} "Forbidden - Not owner or admin"
// @Failure 404 {object} map[string]interface{} "Vehicle not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/user/vehicles/{uuid} [put]
// @Security BearerAuth
func (h *VehicleHandler) UpdateVehicle(c *gin.Context) {
	vehicleUUID := c.Param("uuid")

	// Get authenticated user info from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User not authenticated",
		})
		return
	}

	roleID, exists := c.Get("role_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User role not found",
		})
		return
	}

	// Verify vehicle exists
	existingVehicle, err := h.vehicleRepo.GetByUUID(vehicleUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve vehicle",
			"error":   err.Error(),
		})
		return
	}

	if existingVehicle == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Vehicle not found",
		})
		return
	}

	// Authorization check: Only owner or admin can update
	isOwner := existingVehicle.UserID == userID.(uint64)
	isAdmin := roleID.(uint64) == 1 // Admin role ID is 1

	if !isOwner && !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "You don't have permission to update this vehicle",
		})
		return
	}

	var req UpdateVehicleRequest

	// Bind form data
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Validation failed",
			"errors":  utils.FormatValidationErrorsSimple(err.(validator.ValidationErrors)),
		})
		return
	}

	// Update only provided fields
	if req.Title != "" {
		existingVehicle.Title = req.Title
		existingVehicle.Slug = utils.GenerateSlug(req.Title)
	}
	if req.Category != "" {
		existingVehicle.Category = req.Category
	}
	if req.Description != "" {
		existingVehicle.Description = &req.Description
	}
	if req.Price > 0 {
		existingVehicle.Price = req.Price
	}
	if req.Currency != "" {
		existingVehicle.Currency = req.Currency
	}
	existingVehicle.Negotiable = req.Negotiable

	if req.PersonType != "" {
		personTypeID, err := h.vehicleRepo.GetPersonTypeID(req.PersonType)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid person_type value",
			})
			return
		}
		existingVehicle.PersonTypeID = personTypeID
	}

	if req.Brand != "" {
		existingVehicle.Brand = req.Brand
	}
	if req.Model != "" {
		existingVehicle.Model = req.Model
	}
	if req.EngineCapacity > 0 {
		existingVehicle.EngineCapacity = &req.EngineCapacity
	}
	if req.PowerHP > 0 {
		existingVehicle.PowerHP = &req.PowerHP
	}

	if req.FuelType != "" {
		fuelTypeID, err := h.vehicleRepo.GetFuelTypeID(req.FuelType)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid fuel_type value",
			})
			return
		}
		existingVehicle.FuelTypeID = fuelTypeID
	}

	if req.BodyType != "" {
		bodyTypeID, err := h.vehicleRepo.GetBodyTypeID(req.BodyType)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid body_type value",
			})
			return
		}
		existingVehicle.BodyTypeID = bodyTypeID
	}

	if req.Kilometers > 0 {
		existingVehicle.Kilometers = &req.Kilometers
	}
	if req.Color != "" {
		existingVehicle.Color = &req.Color
	}
	if req.Year > 0 {
		// Validate year
		currentYear := time.Now().Year()
		if req.Year > currentYear+1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Year cannot be more than one year in the future",
			})
			return
		}
		existingVehicle.Year = req.Year
	}
	if req.NumberOfKeys > 0 {
		existingVehicle.NumberOfKeys = &req.NumberOfKeys
	}

	if req.Condition != "" {
		conditionID, err := h.vehicleRepo.GetConditionID(req.Condition)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid condition value",
			})
			return
		}
		existingVehicle.ConditionID = conditionID
	}

	if req.Transmission != "" {
		transmissionID, err := h.vehicleRepo.GetTransmissionID(req.Transmission)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid transmission value",
			})
			return
		}
		existingVehicle.TransmissionID = transmissionID
	}

	if req.Steering != "" {
		steeringID, err := h.vehicleRepo.GetSteeringID(req.Steering)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid steering value",
			})
			return
		}
		existingVehicle.SteeringID = steeringID
	}

	existingVehicle.Registered = req.Registered

	if req.City != "" {
		existingVehicle.City = req.City
	}
	if req.ContactName != "" {
		existingVehicle.ContactName = req.ContactName
	}
	if req.Email != "" {
		existingVehicle.Email = req.Email
	}
	if req.Phone != "" {
		existingVehicle.Phone = &req.Phone
	}

	// Update vehicle in database
	if err := h.vehicleRepo.Update(vehicleUUID, existingVehicle); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update vehicle",
			"error":   err.Error(),
		})
		return
	}

	// Fetch updated vehicle
	updatedVehicle, err := h.vehicleRepo.GetByUUID(vehicleUUID)
	if err != nil {
		updatedVehicle = existingVehicle // Fallback
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Vehicle updated successfully",
		"data":    updatedVehicle,
	})
}
