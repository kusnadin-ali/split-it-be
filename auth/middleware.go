package auth

import (
	"context"
	"github.com/kusnadin-ali/split-it-be/utils"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.WriteError(w, http.StatusUnauthorized, errInvalidToken)
				return
			}

			// format: Bearer <token>
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.WriteError(w, http.StatusUnauthorized, errInvalidToken)
				return
			}

			tokenStr := parts[1]

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				utils.WriteError(w, http.StatusUnauthorized, errInvalidToken)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				utils.WriteError(w, http.StatusUnauthorized, errInvalidToken)
				return
			}

			// ambil userID dari claim (pastikan lu set ini pas login)
			userID, ok := claims["user_id"].(string)
			if !ok || userID == "" {
				utils.WriteError(w, http.StatusUnauthorized, errInvalidToken)
				return
			}

			// 🔥 INI INTINYA: inject ke context
			ctx := context.WithValue(r.Context(), contextKeyUserID, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
