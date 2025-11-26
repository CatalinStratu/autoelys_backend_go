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

func (r *UserRepository) FindByUUID(uuid string) (*models.User, error) {
	query := `
		SELECT id, uuid, role_id, first_name, last_name, email, phone, password_hash, active, accepted_terms_at, created_at, updated_at
		FROM users
		WHERE uuid = ?
	`

	user := &models.User{}
	err := r.db.QueryRow(query, uuid).Scan(
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

func (r *UserRepository) AdminUpdateUser(user *models.User) error {
	query := `
		UPDATE users
		SET role_id = ?, first_name = ?, last_name = ?, email = ?, phone = ?, active = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err := r.db.Exec(query, user.RoleID, user.FirstName, user.LastName, strings.ToLower(user.Email), user.Phone, user.Active, user.ID)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "UNIQUE constraint") {
			return ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUser(id uint64) error {
	query := `DELETE FROM users WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) GetAllUsers(limit, offset int, search string) ([]models.User, int, error) {
	var users []models.User
	var total int

	countQuery := `SELECT COUNT(*) FROM users`
	dataQuery := `
		SELECT id, uuid, role_id, first_name, last_name, email, phone, active, accepted_terms_at, created_at, updated_at
		FROM users
	`

	var args []interface{}
	if search != "" {
		searchPattern := "%" + search + "%"
		whereClause := ` WHERE first_name LIKE ? OR last_name LIKE ? OR email LIKE ?`
		countQuery += whereClause
		dataQuery += whereClause
		args = append(args, searchPattern, searchPattern, searchPattern)
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
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.UUID,
			&user.RoleID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Phone,
			&user.Active,
			&user.AcceptedTermsAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

func TimePtr(t time.Time) *time.Time {
	return &t
}
