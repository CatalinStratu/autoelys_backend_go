package repository

import (
	"autoelys_backend/internal/models"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var ErrServiceNotFound = errors.New("service not found")

type ServiceRepository struct {
	db *sql.DB
}

func NewServiceRepository(db *sql.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) Create(service *models.Service) error {
	service.UUID = uuid.New().String()

	query := `
		INSERT INTO services (uuid, title, description, price, currency, duration_minutes, active)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		service.UUID,
		service.Title,
		service.Description,
		service.Price,
		service.Currency,
		service.DurationMinutes,
		service.Active,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	service.ID = uint64(id)
	return nil
}

func (r *ServiceRepository) FindByUUID(uuid string) (*models.Service, error) {
	query := `
		SELECT id, uuid, title, description, price, currency, duration_minutes, active, created_at, updated_at
		FROM services
		WHERE uuid = ?
	`

	service := &models.Service{}
	err := r.db.QueryRow(query, uuid).Scan(
		&service.ID,
		&service.UUID,
		&service.Title,
		&service.Description,
		&service.Price,
		&service.Currency,
		&service.DurationMinutes,
		&service.Active,
		&service.CreatedAt,
		&service.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrServiceNotFound
	}
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (r *ServiceRepository) Update(service *models.Service) error {
	query := `
		UPDATE services
		SET title = ?, description = ?, price = ?, currency = ?, duration_minutes = ?, active = ?, updated_at = NOW()
		WHERE id = ?
	`
	result, err := r.db.Exec(query, service.Title, service.Description, service.Price, service.Currency, service.DurationMinutes, service.Active, service.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrServiceNotFound
	}

	return nil
}

func (r *ServiceRepository) Delete(id uint64) error {
	query := `DELETE FROM services WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrServiceNotFound
	}

	return nil
}

func (r *ServiceRepository) GetAll(limit, offset int, activeOnly bool) ([]models.Service, int, error) {
	var services []models.Service
	var total int

	countQuery := `SELECT COUNT(*) FROM services`
	dataQuery := `
		SELECT id, uuid, title, description, price, currency, duration_minutes, active, created_at, updated_at
		FROM services
	`

	var args []interface{}
	if activeOnly {
		whereClause := ` WHERE active = true`
		countQuery += whereClause
		dataQuery += whereClause
	}

	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	dataQuery += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.db.Query(dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var service models.Service
		if err := rows.Scan(
			&service.ID,
			&service.UUID,
			&service.Title,
			&service.Description,
			&service.Price,
			&service.Currency,
			&service.DurationMinutes,
			&service.Active,
			&service.CreatedAt,
			&service.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		services = append(services, service)
	}

	return services, total, nil
}
