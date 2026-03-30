package auth

import "github.com/jackc/pgx/v5/pgxpool"

// pool adalah package-level variable — sengaja unexported.
// Hanya bisa di-set dari main via SetPool().
// Pattern ini dari MDA: simple, tanpa dependency injection framework.
var pool *pgxpool.Pool

// jwtSecret juga package-level, di-set dari main.
var jwtSecret string

// SetPool dipanggil sekali di main() sebelum server start.
func SetPool(p *pgxpool.Pool) {
	pool = p
}

// SetJWTSecret dipanggil sekali di main() sebelum server start.
func SetJWTSecret(s string) {
	jwtSecret = s
}
