package auth

import "time"

// User adalah domain object — state yang di-persist ke database.
// Ini murni data, tidak ada method business logic di sini.
type User struct {
	ID           string    `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	AvatarURL    string    `db:"avatar_url"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// UserJSON adalah representasi JSON dari User.
// Dipisah dari domain struct supaya perubahan serialisasi
// tidak menyentuh domain object. PasswordHash tidak pernah keluar.
type UserJSON struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func toUserJSON(u User) UserJSON {
	return UserJSON{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		AvatarURL: u.AvatarURL,
		CreatedAt: u.CreatedAt,
	}
}

// --- Request types (input dari HTTP layer) ---

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// --- Response types (output ke HTTP layer) ---

type authResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	User         UserJSON `json:"user"`
}

// tokenClaims adalah read model — hanya untuk dibaca, tidak dimodifikasi.
type tokenClaims struct {
	UserID string
	Email  string
}

// --- Errors ---

// Sentinel errors — handler pakai errors.Is() untuk cek ini
// dan tentukan HTTP status code yang tepat.
var (
	errEmailExists       = &appError{Code: "email_exists", Message: "email already registered"}
	errInvalidCredential = &appError{Code: "invalid_credentials", Message: "invalid email or password"}
	errUserNotFound      = &appError{Code: "user_not_found", Message: "user not found"}
	errInvalidToken      = &appError{Code: "invalid_token", Message: "token is invalid or expired"}
)

type appError struct {
	Code    string `json:"error"`
	Message string `json:"message"`
}

func (e *appError) Error() string { return e.Message }
