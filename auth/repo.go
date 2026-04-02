package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// findUserByEmail mengambil user dari DB berdasarkan email.
// Return (nil, nil) kalau tidak ditemukan — bukan error.
func findUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, avatar_url, created_at, updated_at
		FROM users WHERE email = $1
	`
	u := &User{}
	err := pool.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Email,
		&u.PasswordHash, &u.AvatarURL,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("findUserByEmail: %w", err)
	}
	return u, nil
}

// findUserByID mengambil user dari DB berdasarkan ID.
func findUserByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT id, first_name, last_name, email, password_hash, avatar_url, created_at, updated_at
		FROM users WHERE id = $1
	`
	u := &User{}
	err := pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.FirstName, &u.LastName, &u.Email,
		&u.PasswordHash, &u.AvatarURL,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("findUserByID: %w", err)
	}
	return u, nil
}

// insertUser menyimpan user baru ke database.
func insertUser(ctx context.Context, u User) error {
	query := `
		INSERT INTO users (id, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := pool.Exec(ctx, query,
		u.ID, u.Email,
		u.PasswordHash,
		u.CreatedAt, u.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insertUser: %w", err)
	}
	return nil
}

func editUser(ctx context.Context, user User) error {
	query := `
		UPDATE users SET first_name = $2, last_name = $3, avatar_url = $4, updated_at = $5
		WHERE id = $1
	`
	_, err := pool.Exec(ctx, query,
		user.ID, user.FirstName, user.LastName, user.AvatarURL, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("updateUser: %w", err)
	}
	return nil
}
