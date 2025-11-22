package repository

import (
	"autoelys_backend/internal/models"
	"database/sql"
)

type AutomobileRepository struct {
	db *sql.DB
}

func NewAutomobileRepository(db *sql.DB) *AutomobileRepository {
	return &AutomobileRepository{db: db}
}

// GetByBrandID retrieves all automobiles for a specific brand with brand details
func (r *AutomobileRepository) GetByBrandID(brandID uint64) ([]models.Automobile, error) {
	query := `SELECT
	            a.id, a.name
	          FROM automobiles a
	          INNER JOIN brands b ON a.brand_id = b.id
	          WHERE a.brand_id = ? AND b.deleted_at IS NULL`

	rows, err := r.db.Query(query, brandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var automobiles []models.Automobile
	for rows.Next() {
		var automobile models.Automobile

		err := rows.Scan(
			&automobile.ID,
			&automobile.Name,
		)
		if err != nil {
			return nil, err
		}

		automobiles = append(automobiles, automobile)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return automobiles, nil
}
