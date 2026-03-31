package ping

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kusnadin-ali/split-it-be/utils"
	"github.com/rs/zerolog/log"
)

type pongResponse struct {
	Message string `json:"message"`
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(utils.ContextKeyUserID).(string)
	if !ok {
		log.Error().Msg("userID not found in context")
		utils.CommonResponse(w, http.StatusUnauthorized, pongResponse{
			Message: "Unauthorized",
		})
		return
	}
	email, ok := r.Context().Value(utils.ContextKeyEmail).(string)
	if !ok {
		log.Error().Msg("email not found in context")
		utils.CommonResponse(w, http.StatusUnauthorized, pongResponse{
			Message: "Unauthorized",
		})
		return
	}

	utils.CommonResponse(w, http.StatusOK, pongResponse{
		Message: "pong! userID: " + userId + ", email: " + email,
	})
}

func Router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", handlePing)
	return r
}
