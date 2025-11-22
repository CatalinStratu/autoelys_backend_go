package repository

import (
	"autoelys_backend/internal/models"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrDuplicateEmail = errors.New("email already exists")
var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	user.Email = strings.ToLower(user.Email)
	user.UUID = uuid.New().String()

	query := `
		INSERT INTO users (uuid, role_id, first_name, last_name, email, phone, password_hash, active, accepted_terms_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		user.UUID,
		user.RoleID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		user.PasswordHash,
		user.Active,
		user.AcceptedTermsAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "UNIQUE constraint") {
			return ErrDuplicateEmail
		}
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = uint64(id)
	return nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	email = strings.ToLower(email)

	query := `
		SELECT id, uuid, role_id, first_name, last_name, email, phone, password_hash, active, accepted_terms_at, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.UUID,
		&user.RoleID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.Active,
		&user.AcceptedTermsAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByID(id uint64) (*models.User, error) {
	query := `
		SELECT id, uuid, role_id, first_name, last_name, email, phone, password_hash, active, accepted_terms_at, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.UUID,
		&user.RoleID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.Active,
		&user.AcceptedTermsAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	email = strings.ToLower(email)

	query := `SELECT COUNT(*) FROM users WHERE email = ?`
	var count int
	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *UserRepository) GetRoleIDByName(name string) (uint64, error) {
	query := `SELECT id FROM roles WHERE name = ?`
	var roleID uint64
	err := r.db.QueryRow(query, name).Scan(&roleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("role not found")
		}
		return 0, err
	}
	return roleID, nil
}

func (r *UserRepository) UpdatePassword(userID uint64, hashedPassword string) error {
	query := `UPDATE users SET password_hash = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.Exec(query, hashedPassword, userID)
	return err
}

func (r *UserRepository) UpdateProfile(user *models.User) error {
	query := `
		UPDATE users
		SET first_name = ?, last_name = ?, phone = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err := r.db.Exec(query, user.FirstName, user.LastName, user.Phone, user.ID)
	return err
}

func TimePtr(t time.Time) *time.Time {
	return &t
}
