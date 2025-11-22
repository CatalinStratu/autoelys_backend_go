package repository

import (
	"autoelys_backend/internal/models"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"
)

var ErrInvalidToken = errors.New("invalid or expired token")
var ErrTokenAlreadyUsed = errors.New("token already used")

type PasswordResetRepository struct {
	db *sql.DB
}

func NewPasswordResetRepository(db *sql.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

func (r *PasswordResetRepository) Create(userID uint64, expirationHours int) (*models.PasswordResetToken, error) {
	token, err := generateSecureToken(32)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(time.Duration(expirationHours) * time.Hour)

	query := `
		INSERT INTO password_reset_tokens (user_id, token, expires_at)
		VALUES (?, ?, ?)
	`

	result, err := r.db.Exec(query, userID, token, expiresAt)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.PasswordResetToken{
		ID:        uint64(id),
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		Used:      false,
		CreatedAt: time.Now(),
	}, nil
}

func (r *PasswordResetRepository) FindByToken(token string) (*models.PasswordResetToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, used, created_at
		FROM password_reset_tokens
		WHERE token = ?
	`

	resetToken := &models.PasswordResetToken{}
	err := r.db.QueryRow(query, token).Scan(
		&resetToken.ID,
		&resetToken.UserID,
		&resetToken.Token,
		&resetToken.ExpiresAt,
		&resetToken.Used,
		&resetToken.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrInvalidToken
	}
	if err != nil {
		return nil, err
	}

	return resetToken, nil
}

func (r *PasswordResetRepository) ValidateToken(token string) (*models.PasswordResetToken, error) {
	resetToken, err := r.FindByToken(token)
	if err != nil {
		return nil, err
	}

	if resetToken.Used {
		return nil, ErrTokenAlreadyUsed
	}

	if time.Now().After(resetToken.ExpiresAt) {
		return nil, ErrInvalidToken
	}

	return resetToken, nil
}

func (r *PasswordResetRepository) MarkAsUsed(id uint64) error {
	query := `UPDATE password_reset_tokens SET used = 1 WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *PasswordResetRepository) DeleteExpiredTokens() error {
	query := `DELETE FROM password_reset_tokens WHERE expires_at < NOW()`
	_, err := r.db.Exec(query)
	return err
}

func (r *PasswordResetRepository) DeleteUserTokens(userID uint64) error {
	query := `DELETE FROM password_reset_tokens WHERE user_id = ?`
	_, err := r.db.Exec(query, userID)
	return err
}

func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
