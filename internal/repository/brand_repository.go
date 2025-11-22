package repository

import (
	"autoelys_backend/internal/models"
	"database/sql"
)

type BrandRepository struct {
	db *sql.DB
}

func NewBrandRepository(db *sql.DB) *BrandRepository {
	return &BrandRepository{db: db}
}

// GetAll retrieves all brands from the database
func (r *BrandRepository) GetAll() ([]models.Brand, error) {
	query := `SELECT id, name
	          FROM brands
	          WHERE deleted_at IS NULL
	          ORDER BY name ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var brands []models.Brand
	for rows.Next() {
		var brand models.Brand
		err := rows.Scan(
			&brand.ID,
			&brand.Name,
		)
		if err != nil {
			return nil, err
		}
		brands = append(brands, brand)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return brands, nil
}

// FindByID retrieves a brand by ID
func (r *BrandRepository) FindByID(id uint64) (*models.Brand, error) {
	query := `SELECT id, url_hash, url, name, logo, deleted_at, created_at, updated_at
	          FROM brands
	          WHERE id = ? AND deleted_at IS NULL`

	var brand models.Brand
	err := r.db.QueryRow(query, id).Scan(
		&brand.ID,
		&brand.URLHash,
		&brand.URL,
		&brand.Name,
		&brand.Logo,
		&brand.DeletedAt,
		&brand.CreatedAt,
		&brand.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &brand, nil
}
