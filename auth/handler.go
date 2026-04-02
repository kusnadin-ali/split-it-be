package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kusnadin-ali/split-it-be/utils"
	"github.com/rs/zerolog/log"
)

// handleRegister menangani POST /api/v1/auth/register
func handleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("decode register request")
		utils.WriteError(w, http.StatusBadRequest, &appError{
			Code:    "bad_request",
			Message: "invalid request body",
		})
		return
	}

	resp, err := register(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("register")
		var appErr *appError
		if errors.As(err, &appErr) {
			switch appErr {
			case errEmailExists:
				utils.WriteError(w, http.StatusConflict, appErr)
			default:
				utils.WriteError(w, http.StatusBadRequest, appErr)
			}
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, &appError{
			Code:    "internal_error",
			Message: "something went wrong",
		})
		return
	}

	utils.CommonResponse(w, http.StatusCreated, resp)
}

// handleLogin menangani POST /api/v1/auth/login
func handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, &appError{
			Code:    "bad_request",
			Message: "invalid request body",
		})
		return
	}

	resp, err := login(r.Context(), req)
	if err != nil {
		var appErr *appError
		if errors.As(err, &appErr) && appErr == errInvalidCredential {
			utils.WriteError(w, http.StatusUnauthorized, appErr)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, &appError{
			Code:    "internal_error",
			Message: "something went wrong",
		})
		return
	}

	utils.CommonResponse(w, http.StatusOK, resp)
}

func handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("decode update user request")
		utils.WriteError(w, http.StatusBadRequest, &appError{
			Code:    "bad_request",
			Message: "invalid request body",
		})
		return
	}

	resp, err := updateAfterRegister(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("update user")
		var appErr *appError
		if errors.As(err, &appErr) {
			switch appErr {
			case errUserNotFound:
				utils.WriteError(w, http.StatusNotFound, appErr)
			default:
				utils.WriteError(w, http.StatusBadRequest, appErr)
			}
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, &appError{
			Code:    "internal_error",
			Message: "something went wrong",
		})
		return
	}

	utils.CommonResponse(w, http.StatusOK, resp)
}

// handleMe menangani GET /api/v1/me — hanya bisa diakses dengan JWT valid
func handleMe(w http.ResponseWriter, r *http.Request) {
	// userID diambil dari context, diset oleh JWTMiddleware
	userID, ok := r.Context().Value(utils.ContextKeyUserID).(string)
	if !ok || userID == "" {
		utils.WriteError(w, http.StatusUnauthorized, errInvalidToken)
		return
	}

	user, err := getMe(r.Context(), userID)
	if err != nil {
		var appErr *appError
		if errors.As(err, &appErr) && appErr == errUserNotFound {
			utils.WriteError(w, http.StatusNotFound, appErr)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, &appError{
			Code:    "internal_error",
			Message: "something went wrong",
		})
		return
	}

	utils.CommonResponse(w, http.StatusOK, user)
}
