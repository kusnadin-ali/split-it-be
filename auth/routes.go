package auth

import (
	"github.com/go-chi/chi/v5"
)

// Router return subrouter untuk public auth endpoints.
// Di-mount di main.go: r.Mount("/api/v1/auth", auth.Router())
func Router() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/register", handleRegister)
	r.Post("/login", handleLogin)
	return r
}

// MeRouter return subrouter untuk endpoint /me.
// Dipasang di dalam protected group yang sudah ada JWTMiddleware.
func MeRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", handleMe)
	return r
}
