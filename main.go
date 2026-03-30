package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/kusnadin-ali/split-it-be/auth"
	"github.com/kusnadin-ali/split-it-be/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func main() {
	utils.LoggerInit()
	var configFileName string
	flag.StringVar(&configFileName, "c", "config.yml", "Config file name")
	flag.Parse()

	cfg := defaultConfig()
	cfg.loadFromEnv()

	if configFileName != "" {
		if err := loadConfigFromFile(configFileName, &cfg); err != nil {
			log.Warn().Str("file", configFileName).Err(err).Msg("cannot load config file, using defaults")
		}
	}

	log.Debug().Any("config", cfg).Msg("config loaded")

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DBConfig.ConnStr())
	if err != nil {
		log.Fatal().Err(err).Msg("unable to connect to database")
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal().Err(err).Msg("database is not reachable")
	}

	log.Info().Msg("database connected")

	// Inject pool ke setiap modul — setiap modul punya SetPool() sendiri
	auth.SetPool(pool)
	auth.SetJWTSecret(cfg.JWTSecret)

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)

	// Mount modul — setiap modul export Router() *chi.Mux
	r.Mount("/api/v1/auth", auth.Router())

	// Protected routes — semua di bawah ini butuh JWT
	//r.Group(func(r chi.Router) {
	//	r.Use(auth.JWTMiddleware)
	//	r.Mount("/api/v1/me", auth.MeRouter())
	//	// nanti: r.Mount("/api/v1/groups", group.Router())
	//	// nanti: r.Mount("/api/v1/splits", split.Router())
	//})

	log.Info().Str("addr", cfg.Listen.Addr()).Msg("starting server")

	fmt.Println("   ___              ________  ________  ___       ___  _________               ___  _________               ________  _______      \n _|\\  \\__          |\\   ____\\|\\   __  \\|\\  \\     |\\  \\|\\___   ___\\            |\\  \\|\\___   ___\\            |\\   __  \\|\\  ___ \\     \n|\\   ____\\         \\ \\  \\___|\\ \\  \\|\\  \\ \\  \\    \\ \\  \\|___ \\  \\_|____________\\ \\  \\|___ \\  \\_|____________\\ \\  \\|\\ /\\ \\   __/|    \n\\ \\  \\___|_         \\ \\_____  \\ \\   ____\\ \\  \\    \\ \\  \\   \\ \\  \\|\\____________\\ \\  \\   \\ \\  \\|\\____________\\ \\   __  \\ \\  \\_|/__  \n \\ \\_____  \\         \\|____|\\  \\ \\  \\___|\\ \\  \\____\\ \\  \\   \\ \\  \\|____________|\\ \\  \\   \\ \\  \\|____________|\\ \\  \\|\\  \\ \\  \\_|\\ \\ \n  \\|____|\\  \\          ____\\_\\  \\ \\__\\    \\ \\_______\\ \\__\\   \\ \\__\\              \\ \\__\\   \\ \\__\\              \\ \\_______\\ \\_______\\\n    ____\\_\\  \\        |\\_________\\|__|     \\|_______|\\|__|    \\|__|               \\|__|    \\|__|               \\|_______|\\|_______|\n   |\\___    __\\       \\|_________|                                                                                                 \n   \\|___|\\__\\_|                                                                                                                    \n        \\|__|                                                                                                                      ")
	fmt.Println("starting server...")
	if err := http.ListenAndServe(cfg.Listen.Addr(), r); err != nil {
		log.Fatal().Err(err).Msg("server stopped")
	}
}

// corsMiddleware diletakkan di main karena ini cross-cutting concern global,
// bukan milik satu modul tertentu.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
