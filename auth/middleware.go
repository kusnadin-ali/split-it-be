package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/kusnadin-ali/split-it-be/utils"
)

// JWTMiddleware adalah satu-satunya exported function di file ini.
// Dipasang di main.go untuk protected routes.
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteError(w, http.StatusUnauthorized, &appError{
				Code:    "missing_token",
				Message: "authorization header is required",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.WriteError(w, http.StatusUnauthorized, &appError{
				Code:    "invalid_token_format",
				Message: "format must be: Bearer {token}",
			})
			return
		}

		claims, err := parseToken(parts[1])
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, &appError{
				Code:    "invalid_token",
				Message: "invalid token",
			})
			return
		}

		// Inject claims ke context — handler ambil via r.Context().Value(...)
		ctx := context.WithValue(r.Context(), utils.ContextKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, utils.ContextKeyEmail, claims.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
