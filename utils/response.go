package utils

import (
	"encoding/json"
	"net/http"
)

func CommonResponse(w http.ResponseWriter, status int, data any) {
	var body struct {
		Code string `json:"code"`
		Msg  string `json:"message"`
		Data any    `json:"data"`
	}

	body.Msg = "success"
	body.Data = data
	body.Code = http.StatusText(status)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(body)
}

func CommonErrorResponse(w http.ResponseWriter, status int, message string) {
	var body struct {
		Code string `json:"code"`
		Msg  string `json:"message"`
	}

	body.Msg = message
	body.Code = http.StatusText(status)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(body)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	CommonErrorResponse(w, status, err.Error())
}
