package models

import "time"

// Lookup table models
type PersonType struct {
	ID          uint8  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type FuelType struct {
	ID          uint8  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type BodyType struct {
	ID          uint8  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type Condition struct {
	ID          uint8  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type Transmission struct {
	ID          uint8  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type Steering struct {
	ID          uint8  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

// Vehicle status constants
const (
	VehicleStatusActive   uint8 = 1
	VehicleStatusInactive uint8 = 2
	VehicleStatusBanned   uint8 = 3
)

// GetStatusName returns the string representation of a vehicle status
func GetStatusName(status uint8) string {
	switch status {
	case VehicleStatusActive:
		return "active"
	case VehicleStatusInactive:
		return "inactive"
	case VehicleStatusBanned:
		return "banned"
	default:
		return "unknown"
	}
}

// Vehicle model
type Vehicle struct {
	ID             uint64    `json:"id,omitempty"`
	UserID         uint64    `json:"user_id,omitempty"`
	Status         uint8     `json:"status"` // 1=active, 2=inactive, 3=banned
	StatusName     string    `json:"status_name,omitempty"`
	Recommended    bool      `json:"recommended"`
	FeaturedImage  *string   `json:"featured_image,omitempty"`
	UUID           string    `json:"uuid"`
	Slug           string    `json:"slug"`
	Title          string    `json:"title"`
	Category       string    `json:"category"`
	Description    *string   `json:"description,omitempty"`
	Price          float64   `json:"price"`
	Currency       string    `json:"currency"`
	Negotiable     bool      `json:"negotiable"`
	PersonTypeID   uint8     `json:"person_type_id"`
	PersonType     string    `json:"person_type"`
	Brand          string    `json:"brand"`
	Model          string    `json:"model"`
	EngineCapacity *int      `json:"engine_capacity,omitempty"`
	PowerHP        *int      `json:"power_hp,omitempty"`
	FuelTypeID     uint8     `json:"fuel_type_id"`
	FuelType       string    `json:"fuel_type"`
	BodyTypeID     uint8     `json:"body_type_id"`
	BodyType       string    `json:"body_type"`
	Kilometers     *int      `json:"kilometers,omitempty"`
	Color          *string   `json:"color,omitempty"`
	Year           int       `json:"year"`
	NumberOfKeys   *int      `json:"number_of_keys,omitempty"`
	ConditionID    uint8     `json:"condition_id"`
	Condition      string    `json:"condition"`
	TransmissionID uint8     `json:"transmission_id"`
	Transmission   string    `json:"transmission"`
	SteeringID     uint8     `json:"steering_id"`
	Steering       string    `json:"steering"`
	Registered     bool      `json:"registered"`
	City           string    `json:"city"`
	ContactName    string    `json:"contact_name"`
	Email          string    `json:"email"`
	Phone          *string   `json:"phone,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relationship
	Images []VehicleImage `json:"images,omitempty"`
}

type VehicleImage struct {
	ID        uint64    `json:"id"`
	VehicleID uint64    `json:"vehicle_id"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}
