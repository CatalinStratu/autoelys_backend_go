package repository

import (
	"autoelys_backend/internal/models"
	"database/sql"
)

type VehicleRepository struct {
	db *sql.DB
}

func NewVehicleRepository(db *sql.DB) *VehicleRepository {
	return &VehicleRepository{db: db}
}

// GetIDByName helper functions for lookup tables
func (r *VehicleRepository) GetPersonTypeID(name string) (uint8, error) {
	var id uint8
	err := r.db.QueryRow("SELECT id FROM person_types WHERE name = ?", name).Scan(&id)
	return id, err
}

func (r *VehicleRepository) GetFuelTypeID(name string) (uint8, error) {
	var id uint8
	err := r.db.QueryRow("SELECT id FROM fuel_types WHERE name = ?", name).Scan(&id)
	return id, err
}

func (r *VehicleRepository) GetBodyTypeID(name string) (uint8, error) {
	var id uint8
	err := r.db.QueryRow("SELECT id FROM body_types WHERE name = ?", name).Scan(&id)
	return id, err
}

func (r *VehicleRepository) GetConditionID(name string) (uint8, error) {
	var id uint8
	err := r.db.QueryRow("SELECT id FROM conditions WHERE name = ?", name).Scan(&id)
	return id, err
}

func (r *VehicleRepository) GetTransmissionID(name string) (uint8, error) {
	var id uint8
	err := r.db.QueryRow("SELECT id FROM transmissions WHERE name = ?", name).Scan(&id)
	return id, err
}

func (r *VehicleRepository) GetSteeringID(name string) (uint8, error) {
	var id uint8
	err := r.db.QueryRow("SELECT id FROM steerings WHERE name = ?", name).Scan(&id)
	return id, err
}

// Create inserts a new vehicle and returns the created vehicle with ID
func (r *VehicleRepository) Create(vehicle *models.Vehicle) (*models.Vehicle, error) {
	// Default status to active if not set
	if vehicle.Status == 0 {
		vehicle.Status = models.VehicleStatusActive
	}

	query := `INSERT INTO vehicles (
		user_id, status, uuid, slug, title, category, description, price, currency, negotiable,
		person_type_id, brand, model, engine_capacity, power_hp,
		fuel_type_id, body_type_id, kilometers, color, year, number_of_keys,
		condition_id, transmission_id, steering_id, registered,
		city, contact_name, email, phone
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query,
		vehicle.UserID,
		vehicle.Status,
		vehicle.UUID,
		vehicle.Slug,
		vehicle.Title,
		vehicle.Category,
		vehicle.Description,
		vehicle.Price,
		vehicle.Currency,
		vehicle.Negotiable,
		vehicle.PersonTypeID,
		vehicle.Brand,
		vehicle.Model,
		vehicle.EngineCapacity,
		vehicle.PowerHP,
		vehicle.FuelTypeID,
		vehicle.BodyTypeID,
		vehicle.Kilometers,
		vehicle.Color,
		vehicle.Year,
		vehicle.NumberOfKeys,
		vehicle.ConditionID,
		vehicle.TransmissionID,
		vehicle.SteeringID,
		vehicle.Registered,
		vehicle.City,
		vehicle.ContactName,
		vehicle.Email,
		vehicle.Phone,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	vehicle.ID = uint64(id)
	return vehicle, nil
}

// CreateImage inserts a vehicle image
func (r *VehicleRepository) CreateImage(vehicleID uint64, imageURL string) error {
	query := `INSERT INTO vehicle_images (vehicle_id, image_url) VALUES (?, ?)`
	_, err := r.db.Exec(query, vehicleID, imageURL)
	return err
}

// GetByID retrieves a vehicle by UUID with its images and lookup table data
func (r *VehicleRepository) GetByID(uuid uint64) (*models.Vehicle, error) {
	query := `SELECT
		v.uuid, v.status, v.slug, v.title, v.category, v.description, v.price, v.currency, v.negotiable,
		v.person_type_id, pt.name as person_type_name,
		v.brand, v.model, v.engine_capacity, v.power_hp,
		v.fuel_type_id, ft.name as fuel_type_name,
		v.body_type_id, bt.name as body_type_name,
		v.kilometers, v.color, v.year, v.number_of_keys,
		v.condition_id, c.name as condition_name,
		v.transmission_id, t.name as transmission_name,
		v.steering_id, s.name as steering_name,
		v.registered,
		v.city, v.contact_name, v.email, v.phone, v.created_at, v.updated_at
	FROM vehicles v
	LEFT JOIN person_types pt ON v.person_type_id = pt.id
	LEFT JOIN fuel_types ft ON v.fuel_type_id = ft.id
	LEFT JOIN body_types bt ON v.body_type_id = bt.id
	LEFT JOIN conditions c ON v.condition_id = c.id
	LEFT JOIN transmissions t ON v.transmission_id = t.id
	LEFT JOIN steerings s ON v.steering_id = s.id
	WHERE v.uuid = ?`

	vehicle := &models.Vehicle{}
	var personTypeName, fuelTypeName, bodyTypeName, conditionName, transmissionName, steeringName string

	err := r.db.QueryRow(query, uuid).Scan(
		&vehicle.UUID,
		&vehicle.Status,
		&vehicle.Slug,
		&vehicle.Title,
		&vehicle.Category,
		&vehicle.Description,
		&vehicle.Price,
		&vehicle.Currency,
		&vehicle.Negotiable,
		&vehicle.PersonTypeID,
		&personTypeName,
		&vehicle.Brand,
		&vehicle.Model,
		&vehicle.EngineCapacity,
		&vehicle.PowerHP,
		&vehicle.FuelTypeID,
		&fuelTypeName,
		&vehicle.BodyTypeID,
		&bodyTypeName,
		&vehicle.Kilometers,
		&vehicle.Color,
		&vehicle.Year,
		&vehicle.NumberOfKeys,
		&vehicle.ConditionID,
		&conditionName,
		&vehicle.TransmissionID,
		&transmissionName,
		&vehicle.SteeringID,
		&steeringName,
		&vehicle.Registered,
		&vehicle.City,
		&vehicle.ContactName,
		&vehicle.Email,
		&vehicle.Phone,
		&vehicle.CreatedAt,
		&vehicle.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Set the name fields
	vehicle.PersonType = personTypeName
	vehicle.FuelType = fuelTypeName
	vehicle.BodyType = bodyTypeName
	vehicle.Condition = conditionName
	vehicle.Transmission = transmissionName
	vehicle.Steering = steeringName
	vehicle.StatusName = models.GetStatusName(vehicle.Status)

	// Get images
	images, err := r.GetImagesByVehicleID(uuid)
	if err != nil {
		return nil, err
	}
	vehicle.Images = images

	return vehicle, nil
}

// GetBySlug retrieves a vehicle by slug with its images and lookup table data
func (r *VehicleRepository) GetBySlug(slug string) (*models.Vehicle, error) {
	query := `SELECT
		v.id, v.user_id, v.status, v.uuid, v.slug, v.title, v.category, v.description, v.price, v.currency, v.negotiable,
		v.person_type_id, pt.name as person_type_name,
		v.brand, v.model, v.engine_capacity, v.power_hp,
		v.fuel_type_id, ft.name as fuel_type_name,
		v.body_type_id, bt.name as body_type_name,
		v.kilometers, v.color, v.year, v.number_of_keys,
		v.condition_id, c.name as condition_name,
		v.transmission_id, t.name as transmission_name,
		v.steering_id, s.name as steering_name,
		v.registered,
		v.city, v.contact_name, v.email, v.phone, v.created_at, v.updated_at
	FROM vehicles v
	LEFT JOIN person_types pt ON v.person_type_id = pt.id
	LEFT JOIN fuel_types ft ON v.fuel_type_id = ft.id
	LEFT JOIN body_types bt ON v.body_type_id = bt.id
	LEFT JOIN conditions c ON v.condition_id = c.id
	LEFT JOIN transmissions t ON v.transmission_id = t.id
	LEFT JOIN steerings s ON v.steering_id = s.id
	WHERE v.slug = ?`

	vehicle := &models.Vehicle{}
	var personTypeName, fuelTypeName, bodyTypeName, conditionName, transmissionName, steeringName string

	err := r.db.QueryRow(query, slug).Scan(
		&vehicle.ID,
		&vehicle.UserID,
		&vehicle.Status,
		&vehicle.UUID,
		&vehicle.Slug,
		&vehicle.Title,
		&vehicle.Category,
		&vehicle.Description,
		&vehicle.Price,
		&vehicle.Currency,
		&vehicle.Negotiable,
		&vehicle.PersonTypeID,
		&personTypeName,
		&vehicle.Brand,
		&vehicle.Model,
		&vehicle.EngineCapacity,
		&vehicle.PowerHP,
		&vehicle.FuelTypeID,
		&fuelTypeName,
		&vehicle.BodyTypeID,
		&bodyTypeName,
		&vehicle.Kilometers,
		&vehicle.Color,
		&vehicle.Year,
		&vehicle.NumberOfKeys,
		&vehicle.ConditionID,
		&conditionName,
		&vehicle.TransmissionID,
		&transmissionName,
		&vehicle.SteeringID,
		&steeringName,
		&vehicle.Registered,
		&vehicle.City,
		&vehicle.ContactName,
		&vehicle.Email,
		&vehicle.Phone,
		&vehicle.CreatedAt,
		&vehicle.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Set the name fields
	vehicle.PersonType = personTypeName
	vehicle.FuelType = fuelTypeName
	vehicle.BodyType = bodyTypeName
	vehicle.Condition = conditionName
	vehicle.Transmission = transmissionName
	vehicle.Steering = steeringName
	vehicle.StatusName = models.GetStatusName(vehicle.Status)

	// Get images
	images, err := r.GetImagesByVehicleID(vehicle.ID)
	if err != nil {
		return nil, err
	}
	vehicle.Images = images

	return vehicle, nil
}

// GetByUUID retrieves a vehicle by UUID with its images and lookup table data
func (r *VehicleRepository) GetByUUID(uuid string) (*models.Vehicle, error) {
	query := `SELECT
		v.id, v.user_id, v.status, v.uuid, v.slug, v.title, v.category, v.description, v.price, v.currency, v.negotiable,
		v.person_type_id, pt.name as person_type_name,
		v.brand, v.model, v.engine_capacity, v.power_hp,
		v.fuel_type_id, ft.name as fuel_type_name,
		v.body_type_id, bt.name as body_type_name,
		v.kilometers, v.color, v.year, v.number_of_keys,
		v.condition_id, c.name as condition_name,
		v.transmission_id, t.name as transmission_name,
		v.steering_id, s.name as steering_name,
		v.registered,
		v.city, v.contact_name, v.email, v.phone, v.created_at, v.updated_at
	FROM vehicles v
	LEFT JOIN person_types pt ON v.person_type_id = pt.id
	LEFT JOIN fuel_types ft ON v.fuel_type_id = ft.id
	LEFT JOIN body_types bt ON v.body_type_id = bt.id
	LEFT JOIN conditions c ON v.condition_id = c.id
	LEFT JOIN transmissions t ON v.transmission_id = t.id
	LEFT JOIN steerings s ON v.steering_id = s.id
	WHERE v.uuid = ?`

	vehicle := &models.Vehicle{}
	var personTypeName, fuelTypeName, bodyTypeName, conditionName, transmissionName, steeringName string

	err := r.db.QueryRow(query, uuid).Scan(
		&vehicle.ID,
		&vehicle.UserID,
		&vehicle.Status,
		&vehicle.UUID,
		&vehicle.Slug,
		&vehicle.Title,
		&vehicle.Category,
		&vehicle.Description,
		&vehicle.Price,
		&vehicle.Currency,
		&vehicle.Negotiable,
		&vehicle.PersonTypeID,
		&personTypeName,
		&vehicle.Brand,
		&vehicle.Model,
		&vehicle.EngineCapacity,
		&vehicle.PowerHP,
		&vehicle.FuelTypeID,
		&fuelTypeName,
		&vehicle.BodyTypeID,
		&bodyTypeName,
		&vehicle.Kilometers,
		&vehicle.Color,
		&vehicle.Year,
		&vehicle.NumberOfKeys,
		&vehicle.ConditionID,
		&conditionName,
		&vehicle.TransmissionID,
		&transmissionName,
		&vehicle.SteeringID,
		&steeringName,
		&vehicle.Registered,
		&vehicle.City,
		&vehicle.ContactName,
		&vehicle.Email,
		&vehicle.Phone,
		&vehicle.CreatedAt,
		&vehicle.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Set the name fields
	vehicle.PersonType = personTypeName
	vehicle.FuelType = fuelTypeName
	vehicle.BodyType = bodyTypeName
	vehicle.Condition = conditionName
	vehicle.Transmission = transmissionName
	vehicle.Steering = steeringName
	vehicle.StatusName = models.GetStatusName(vehicle.Status)

	// Get images
	images, err := r.GetImagesByVehicleID(vehicle.ID)
	if err != nil {
		return nil, err
	}
	vehicle.Images = images

	return vehicle, nil
}

// Update updates a vehicle by UUID
func (r *VehicleRepository) Update(uuid string, vehicle *models.Vehicle) error {
	query := `UPDATE vehicles SET
		slug = ?, title = ?, category = ?, description = ?, price = ?, currency = ?, negotiable = ?,
		person_type_id = ?, brand = ?, model = ?, engine_capacity = ?, power_hp = ?,
		fuel_type_id = ?, body_type_id = ?, kilometers = ?, color = ?, year = ?, number_of_keys = ?,
		condition_id = ?, transmission_id = ?, steering_id = ?, registered = ?,
		city = ?, contact_name = ?, email = ?, phone = ?
	WHERE uuid = ?`

	_, err := r.db.Exec(query,
		vehicle.Slug,
		vehicle.Title,
		vehicle.Category,
		vehicle.Description,
		vehicle.Price,
		vehicle.Currency,
		vehicle.Negotiable,
		vehicle.PersonTypeID,
		vehicle.Brand,
		vehicle.Model,
		vehicle.EngineCapacity,
		vehicle.PowerHP,
		vehicle.FuelTypeID,
		vehicle.BodyTypeID,
		vehicle.Kilometers,
		vehicle.Color,
		vehicle.Year,
		vehicle.NumberOfKeys,
		vehicle.ConditionID,
		vehicle.TransmissionID,
		vehicle.SteeringID,
		vehicle.Registered,
		vehicle.City,
		vehicle.ContactName,
		vehicle.Email,
		vehicle.Phone,
		uuid,
	)

	return err
}

// GetImagesByVehicleID retrieves all images for a vehicle
func (r *VehicleRepository) GetImagesByVehicleID(vehicleID uint64) ([]models.VehicleImage, error) {
	query := `SELECT id, vehicle_id, image_url, created_at FROM vehicle_images WHERE vehicle_id = ?`

	rows, err := r.db.Query(query, vehicleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []models.VehicleImage
	for rows.Next() {
		var image models.VehicleImage
		err := rows.Scan(&image.ID, &image.VehicleID, &image.ImageURL, &image.CreatedAt)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	return images, rows.Err()
}

// GetAllPersonTypes retrieves all person types
func (r *VehicleRepository) GetAllPersonTypes() ([]models.PersonType, error) {
	query := "SELECT id, name, display_name FROM person_types"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []models.PersonType
	for rows.Next() {
		var t models.PersonType
		if err := rows.Scan(&t.ID, &t.Name, &t.DisplayName); err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	return types, rows.Err()
}

// GetAllFuelTypes retrieves all fuel types
func (r *VehicleRepository) GetAllFuelTypes() ([]models.FuelType, error) {
	query := "SELECT id, name, display_name FROM fuel_types"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []models.FuelType
	for rows.Next() {
		var t models.FuelType
		if err := rows.Scan(&t.ID, &t.Name, &t.DisplayName); err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	return types, rows.Err()
}

// GetAllBodyTypes retrieves all body types
func (r *VehicleRepository) GetAllBodyTypes() ([]models.BodyType, error) {
	query := "SELECT id, name, display_name FROM body_types"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []models.BodyType
	for rows.Next() {
		var t models.BodyType
		if err := rows.Scan(&t.ID, &t.Name, &t.DisplayName); err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	return types, rows.Err()
}

// GetAllConditions retrieves all conditions
func (r *VehicleRepository) GetAllConditions() ([]models.Condition, error) {
	query := "SELECT id, name, display_name FROM conditions"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []models.Condition
	for rows.Next() {
		var t models.Condition
		if err := rows.Scan(&t.ID, &t.Name, &t.DisplayName); err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	return types, rows.Err()
}

// GetAllTransmissions retrieves all transmissions
func (r *VehicleRepository) GetAllTransmissions() ([]models.Transmission, error) {
	query := "SELECT id, name, display_name FROM transmissions"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []models.Transmission
	for rows.Next() {
		var t models.Transmission
		if err := rows.Scan(&t.ID, &t.Name, &t.DisplayName); err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	return types, rows.Err()
}

// GetAllSteerings retrieves all steerings
func (r *VehicleRepository) GetAllSteerings() ([]models.Steering, error) {
	query := "SELECT id, name, display_name FROM steerings"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []models.Steering
	for rows.Next() {
		var t models.Steering
		if err := rows.Scan(&t.ID, &t.Name, &t.DisplayName); err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	return types, rows.Err()
}

// VehicleSearchParams holds all search and filter parameters
type VehicleSearchParams struct {
	Search       string
	Brand        string
	Model        string
	FuelType     string
	BodyType     string
	Transmission string
	Condition    string
	MinPrice     float64
	MaxPrice     float64
	MinYear      int
	MaxYear      int
	City         string
	Limit        int
	Offset       int
}

// GetAll retrieves all vehicles with optional search filters
func (r *VehicleRepository) GetAll(params VehicleSearchParams) ([]models.Vehicle, int, error) {
	baseQuery := `SELECT
		v.id, v.user_id, v.status, v.uuid, v.slug, v.title, v.category, v.description, v.price, v.currency, v.negotiable,
		v.person_type_id, pt.name as person_type_name,
		v.brand, v.model, v.engine_capacity, v.power_hp,
		v.fuel_type_id, ft.name as fuel_type_name,
		v.body_type_id, bt.name as body_type_name,
		v.kilometers, v.color, v.year, v.number_of_keys,
		v.condition_id, c.name as condition_name,
		v.transmission_id, t.name as transmission_name,
		v.steering_id, s.name as steering_name,
		v.registered,
		v.city, v.contact_name, v.email, v.phone, v.created_at, v.updated_at
	FROM vehicles v
	LEFT JOIN person_types pt ON v.person_type_id = pt.id
	LEFT JOIN fuel_types ft ON v.fuel_type_id = ft.id
	LEFT JOIN body_types bt ON v.body_type_id = bt.id
	LEFT JOIN conditions c ON v.condition_id = c.id
	LEFT JOIN transmissions t ON v.transmission_id = t.id
	LEFT JOIN steerings s ON v.steering_id = s.id
	WHERE v.status = ?`

	countQuery := `SELECT COUNT(*) FROM vehicles v
	LEFT JOIN fuel_types ft ON v.fuel_type_id = ft.id
	LEFT JOIN body_types bt ON v.body_type_id = bt.id
	LEFT JOIN transmissions t ON v.transmission_id = t.id
	LEFT JOIN conditions c ON v.condition_id = c.id
	WHERE v.status = ?`

	args := []interface{}{models.VehicleStatusActive}
	countArgs := []interface{}{models.VehicleStatusActive}

	// Add search filter
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		baseQuery += " AND (v.title LIKE ? OR v.brand LIKE ? OR v.model LIKE ? OR v.description LIKE ?)"
		countQuery += " AND (v.title LIKE ? OR v.brand LIKE ? OR v.model LIKE ? OR v.description LIKE ?)"
		args = append(args, searchPattern, searchPattern, searchPattern, searchPattern)
		countArgs = append(countArgs, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Add brand filter
	if params.Brand != "" {
		baseQuery += " AND v.brand = ?"
		countQuery += " AND v.brand = ?"
		args = append(args, params.Brand)
		countArgs = append(countArgs, params.Brand)
	}

	// Add model filter
	if params.Model != "" {
		baseQuery += " AND v.model = ?"
		countQuery += " AND v.model = ?"
		args = append(args, params.Model)
		countArgs = append(countArgs, params.Model)
	}

	// Add fuel type filter
	if params.FuelType != "" {
		baseQuery += " AND ft.name = ?"
		countQuery += " AND ft.name = ?"
		args = append(args, params.FuelType)
		countArgs = append(countArgs, params.FuelType)
	}

	// Add body type filter
	if params.BodyType != "" {
		baseQuery += " AND bt.name = ?"
		countQuery += " AND bt.name = ?"
		args = append(args, params.BodyType)
		countArgs = append(countArgs, params.BodyType)
	}

	// Add transmission filter
	if params.Transmission != "" {
		baseQuery += " AND t.name = ?"
		countQuery += " AND t.name = ?"
		args = append(args, params.Transmission)
		countArgs = append(countArgs, params.Transmission)
	}

	// Add condition filter
	if params.Condition != "" {
		baseQuery += " AND c.name = ?"
		countQuery += " AND c.name = ?"
		args = append(args, params.Condition)
		countArgs = append(countArgs, params.Condition)
	}

	// Add price range filter
	if params.MinPrice > 0 {
		baseQuery += " AND v.price >= ?"
		countQuery += " AND v.price >= ?"
		args = append(args, params.MinPrice)
		countArgs = append(countArgs, params.MinPrice)
	}
	if params.MaxPrice > 0 {
		baseQuery += " AND v.price <= ?"
		countQuery += " AND v.price <= ?"
		args = append(args, params.MaxPrice)
		countArgs = append(countArgs, params.MaxPrice)
	}

	// Add year range filter
	if params.MinYear > 0 {
		baseQuery += " AND v.year >= ?"
		countQuery += " AND v.year >= ?"
		args = append(args, params.MinYear)
		countArgs = append(countArgs, params.MinYear)
	}
	if params.MaxYear > 0 {
		baseQuery += " AND v.year <= ?"
		countQuery += " AND v.year <= ?"
		args = append(args, params.MaxYear)
		countArgs = append(countArgs, params.MaxYear)
	}

	// Add city filter
	if params.City != "" {
		baseQuery += " AND v.city = ?"
		countQuery += " AND v.city = ?"
		args = append(args, params.City)
		countArgs = append(countArgs, params.City)
	}

	// Get total count
	var total int
	err := r.db.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Add ordering and pagination
	baseQuery += " ORDER BY v.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, params.Limit, params.Offset)

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var vehicles []models.Vehicle
	for rows.Next() {
		vehicle := models.Vehicle{}
		var personTypeName, fuelTypeName, bodyTypeName, conditionName, transmissionName, steeringName string

		err := rows.Scan(
			&vehicle.ID,
			&vehicle.UserID,
			&vehicle.Status,
			&vehicle.UUID,
			&vehicle.Slug,
			&vehicle.Title,
			&vehicle.Category,
			&vehicle.Description,
			&vehicle.Price,
			&vehicle.Currency,
			&vehicle.Negotiable,
			&vehicle.PersonTypeID,
			&personTypeName,
			&vehicle.Brand,
			&vehicle.Model,
			&vehicle.EngineCapacity,
			&vehicle.PowerHP,
			&vehicle.FuelTypeID,
			&fuelTypeName,
			&vehicle.BodyTypeID,
			&bodyTypeName,
			&vehicle.Kilometers,
			&vehicle.Color,
			&vehicle.Year,
			&vehicle.NumberOfKeys,
			&vehicle.ConditionID,
			&conditionName,
			&vehicle.TransmissionID,
			&transmissionName,
			&vehicle.SteeringID,
			&steeringName,
			&vehicle.Registered,
			&vehicle.City,
			&vehicle.ContactName,
			&vehicle.Email,
			&vehicle.Phone,
			&vehicle.CreatedAt,
			&vehicle.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		// Set the name fields
		vehicle.PersonType = personTypeName
		vehicle.FuelType = fuelTypeName
		vehicle.BodyType = bodyTypeName
		vehicle.Condition = conditionName
		vehicle.Transmission = transmissionName
		vehicle.Steering = steeringName
		vehicle.StatusName = models.GetStatusName(vehicle.Status)

		// Get images for this vehicle
		images, err := r.GetImagesByVehicleID(vehicle.ID)
		if err != nil {
			return nil, 0, err
		}
		vehicle.Images = images

		vehicles = append(vehicles, vehicle)
	}

	return vehicles, total, rows.Err()
}

// GetByUserID retrieves all vehicles for a specific user
func (r *VehicleRepository) GetByUserID(userID uint64) ([]models.Vehicle, error) {
	query := `SELECT
		v.id, v.user_id, v.status, v.uuid, v.slug, v.title, v.category, v.description, v.price, v.currency, v.negotiable,
		v.person_type_id, pt.name as person_type_name,
		v.brand, v.model, v.engine_capacity, v.power_hp,
		v.fuel_type_id, ft.name as fuel_type_name,
		v.body_type_id, bt.name as body_type_name,
		v.kilometers, v.color, v.year, v.number_of_keys,
		v.condition_id, c.name as condition_name,
		v.transmission_id, t.name as transmission_name,
		v.steering_id, s.name as steering_name,
		v.registered,
		v.city, v.contact_name, v.email, v.phone, v.created_at, v.updated_at
	FROM vehicles v
	LEFT JOIN person_types pt ON v.person_type_id = pt.id
	LEFT JOIN fuel_types ft ON v.fuel_type_id = ft.id
	LEFT JOIN body_types bt ON v.body_type_id = bt.id
	LEFT JOIN conditions c ON v.condition_id = c.id
	LEFT JOIN transmissions t ON v.transmission_id = t.id
	LEFT JOIN steerings s ON v.steering_id = s.id
	WHERE v.user_id = ?
	ORDER BY v.created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []models.Vehicle
	for rows.Next() {
		vehicle := models.Vehicle{}
		var personTypeName, fuelTypeName, bodyTypeName, conditionName, transmissionName, steeringName string

		err := rows.Scan(
			&vehicle.ID,
			&vehicle.UserID,
			&vehicle.Status,
			&vehicle.UUID,
			&vehicle.Slug,
			&vehicle.Title,
			&vehicle.Category,
			&vehicle.Description,
			&vehicle.Price,
			&vehicle.Currency,
			&vehicle.Negotiable,
			&vehicle.PersonTypeID,
			&personTypeName,
			&vehicle.Brand,
			&vehicle.Model,
			&vehicle.EngineCapacity,
			&vehicle.PowerHP,
			&vehicle.FuelTypeID,
			&fuelTypeName,
			&vehicle.BodyTypeID,
			&bodyTypeName,
			&vehicle.Kilometers,
			&vehicle.Color,
			&vehicle.Year,
			&vehicle.NumberOfKeys,
			&vehicle.ConditionID,
			&conditionName,
			&vehicle.TransmissionID,
			&transmissionName,
			&vehicle.SteeringID,
			&steeringName,
			&vehicle.Registered,
			&vehicle.City,
			&vehicle.ContactName,
			&vehicle.Email,
			&vehicle.Phone,
			&vehicle.CreatedAt,
			&vehicle.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Set the name fields
		vehicle.PersonType = personTypeName
		vehicle.FuelType = fuelTypeName
		vehicle.BodyType = bodyTypeName
		vehicle.Condition = conditionName
		vehicle.Transmission = transmissionName
		vehicle.Steering = steeringName
		vehicle.StatusName = models.GetStatusName(vehicle.Status)

		// Get images for this vehicle
		images, err := r.GetImagesByVehicleID(vehicle.ID)
		if err != nil {
			return nil, err
		}
		vehicle.Images = images

		vehicles = append(vehicles, vehicle)
	}

	return vehicles, rows.Err()
}
