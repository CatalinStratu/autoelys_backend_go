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
// @Summary Create a new vehicle listing
// @Description Add a new vehicle with images and all details
// @Tags vehicles
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Vehicle title"
// @Param category formData string true "Vehicle category"
// @Param description formData string false "Vehicle description"
// @Param price formData number true "Price"
// @Param currency formData string true "Currency (e.g., lei)"
// @Param negotiable formData boolean false "Price negotiable"
// @Param person_type formData string true "Person type (persoana_fizica, firma)"
// @Param brand formData string true "Brand"
// @Param model formData string true "Model"
// @Param engine_capacity formData integer false "Engine capacity (cm3)"
// @Param power_hp formData integer false "Power (HP)"
// @Param fuel_type formData string true "Fuel type"
// @Param body_type formData string true "Body type"
// @Param kilometers formData integer false "Kilometers"
// @Param color formData string false "Color"
// @Param year formData integer true "Year"
// @Param number_of_keys formData integer false "Number of keys"
// @Param condition formData string true "Condition (utilizat, nou)"
// @Param transmission formData string true "Transmission (manuala, automata)"
// @Param steering formData string true "Steering (stanga, dreapta)"
// @Param registered formData boolean false "Registered"
// @Param city formData string true "City"
// @Param contact_name formData string true "Contact name"
// @Param email formData string true "Email"
// @Param phone formData string false "Phone"
// @Param images formData file false "Vehicle images (max 8)"
// @Success 201 {object} map[string]interface{} "Vehicle created successfully"
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/vehicles [post]
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

	// Generate UUID and slug
	vehicleUUID := uuid.New().String()
	slug := utils.GenerateSlug(req.Title)

	// Create vehicle model
	vehicle := &models.Vehicle{
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


// GetVehicle godoc
// @Summary Get vehicle by ID
// @Description Retrieve a vehicle with all its details and images
// @Tags vehicles
// @Accept json
// @Produce json
// @Param id path int true "Vehicle ID"
// @Success 200 {object} map[string]interface{} "Vehicle details"
// @Failure 404 {object} map[string]interface{} "Vehicle not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/vehicles/{id} [get]
func (h *VehicleHandler) GetVehicle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid vehicle ID",
		})
		return
	}

	vehicle, err := h.vehicleRepo.GetByID(id)
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
