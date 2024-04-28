package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redvant/lenslocked/rand"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID     int
	UserID int
	// Token is only set when creating a new PasswordReset.
	// When look up a passwordReset this will be left empty.
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes to use when
	// generating each password reset token. If this value is not set
	// or is less than the MinBytesPerToken const it will be ignored
	// and MinBytesPerToken will be used.
	BytesPerToken int
	// Duration is the amount of time that a PasswordReset is valid for.
	// Default to DefaultResetDuration.
	Duration time.Duration
}

func (prs *PasswordResetService) Create(email string) (*PasswordReset, error) {
	// Verify email for user, retrieve userID
	email = strings.ToLower(email)
	var userID int
	row := prs.DB.QueryRow(`
		SELECT id
		FROM users
		WHERE email = $1;
	`, email)
	err := row.Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEmailNotFound
		}
		return nil, fmt.Errorf("create: %w", err)
	}

	// Build the password Reset
	bytesPerToken := prs.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	duration := prs.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}
	pwReset := PasswordReset{
		UserID:    userID,
		Token:     token,
		TokenHash: prs.hash(token),
		ExpiresAt: time.Now().Add(duration),
	}

	// Insert or update password reset into DB
	row = prs.DB.QueryRow(`
		INSERT INTO password_resets (user_id, token_hash, expires_at)
		VALUES($1, $2, $3) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2, expires_at = $3
		RETURNING id;
	`, pwReset.UserID, pwReset.TokenHash, pwReset.ExpiresAt)
	err = row.Scan(&pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &pwReset, nil
}

func (prs *PasswordResetService) Consume(token string) (*User, error) {
	tokenHash := prs.hash(token)
	var user User
	var pwReset PasswordReset
	row := prs.DB.QueryRow(`
		SELECT u.id, u.email, u.password_hash, pr.id, pr.expires_at
		FROM password_resets as pr
		JOIN users as u
		ON pr.user_id = u.id
		WHERE pr.token_hash = $1;
	`, tokenHash)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash,
		&pwReset.ID, &pwReset.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}

	// Delete password reset token
	err = prs.delete(pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}

	// Check password token expiration
	if time.Now().After(pwReset.ExpiresAt) {
		return nil, fmt.Errorf("token expired: %v", token)
	}
	return &user, nil
}

func (prs *PasswordResetService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

func (prs *PasswordResetService) delete(pwResetID int) error {
	_, err := prs.DB.Exec(`
		DELETE FROM password_resets
		WHERE id = $1;
	`, pwResetID)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}
