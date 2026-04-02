package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// register menangani business logic registrasi user baru.
// Fungsi murni — tidak tahu apapun tentang HTTP.
func register(ctx context.Context, req registerRequest) (authResponse, error) {
	if err := validateRegisterRequest(req); err != nil {
		return authResponse{}, err
	}

	existing, err := findUserByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Msg("find user by email")
		return authResponse{}, err
	}
	if existing != nil {
		log.Error().Str("email", req.Email).Msg("email already exists")
		return authResponse{}, errEmailExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		log.Error().Err(err).Msg("hash password")
		return authResponse{}, fmt.Errorf("hash password: %w", err)
	}

	now := time.Now()
	user := User{
		ID:           uuid.NewString(),
		Email:        req.Email,
		PasswordHash: string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := insertUser(ctx, user); err != nil {
		log.Error().Err(err).Msg("insert user")
		return authResponse{}, err
	}

	return buildAuthResponse(user)
}

// login memvalidasi kredensial dan return token pair.
func login(ctx context.Context, req loginRequest) (authResponse, error) {
	if err := validateLoginRequest(req); err != nil {
		return authResponse{}, err
	}

	user, err := findUserByEmail(ctx, req.Email)
	if err != nil {
		return authResponse{}, err
	}
	// Sengaja return error yang sama untuk email tidak ditemukan maupun
	// password salah — supaya attacker tidak bisa enumerate email.
	if user == nil {
		return authResponse{}, errInvalidCredential
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return authResponse{}, errInvalidCredential
	}

	return buildAuthResponse(*user)
}

// getMe mengambil data user yang sedang login berdasarkan ID dari JWT claims.
func getMe(ctx context.Context, userID string) (UserJSON, error) {
	user, err := findUserByID(ctx, userID)
	if err != nil {
		return UserJSON{}, err
	}
	if user == nil {
		return UserJSON{}, errUserNotFound
	}
	return toUserJSON(*user), nil
}

// --- Validation ---

func validateRegisterRequest(req registerRequest) error {
	if req.Email == "" {
		return &appError{Code: "validation_error", Message: "email is required"}
	}
	if len(req.Password) < 8 {
		return &appError{Code: "validation_error", Message: "password must be at least 8 characters"}
	}
	return nil
}

func validateLoginRequest(req loginRequest) error {
	if req.Email == "" {
		return &appError{Code: "validation_error", Message: "email is required"}
	}
	if req.Password == "" {
		return &appError{Code: "validation_error", Message: "password is required"}
	}
	return nil
}

// --- Token helpers ---

func buildAuthResponse(user User) (authResponse, error) {
	access, err := generateToken(user, 15*time.Minute)
	if err != nil {
		return authResponse{}, err
	}
	refresh, err := generateToken(user, 7*24*time.Hour)
	if err != nil {
		return authResponse{}, err
	}
	return authResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		User:         toUserJSON(user),
	}, nil
}

func generateToken(user User, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(duration).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}

func parseToken(tokenString string) (tokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return tokenClaims{}, errInvalidToken
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return tokenClaims{}, errInvalidToken
	}
	return tokenClaims{
		UserID: claims["user_id"].(string),
		Email:  claims["email"].(string),
	}, nil
}

func updateAfterRegister(ctx context.Context, req updateUserRequest) (UserDetailJSON, error) {
	user, err := findUserByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Msg("find user by email")
		return UserDetailJSON{}, err
	}
	if user == nil {
		log.Error().Str("email", req.Email).Msg("user not found")
		return UserDetailJSON{}, errUserNotFound
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.AvatarURL = req.AvatarURL
	user.UpdatedAt = time.Now()

	if err := editUser(ctx, *user); err != nil {
		log.Error().Err(err).Msg("update user")
		return UserDetailJSON{}, err
	}
	return toUserDetailJSON(*user), nil
}
