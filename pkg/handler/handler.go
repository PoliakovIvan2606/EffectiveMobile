package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

const (
	StatusErrResponse = iota
	StatusOkResponse
)

type ApiOkResponse struct {
    Status  string      `json:"status"`           
    Data    interface{} `json:"data,omitempty"`
}

type ApiErrResponse struct {
    Status  string      `json:"status"`           
    Err    string `json:"error,omitempty"`
}

func ErrResponse(w http.ResponseWriter, msg string, err error, http_status int) {
	slog.Debug(err.Error(), )

	MakeResponse(w, map[string]interface{}{
		"status": StatusErrResponse,
		"error": msg,
	}, http_status)
}

func OkResponse(w http.ResponseWriter, out interface{}, http_status int) {
	MakeResponse(w, map[string]interface{}{
		"status": StatusOkResponse,
		"data": out,
	}, http_status)
}


func MakeResponse(w http.ResponseWriter, out interface{}, http_status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http_status)
	json.NewEncoder(w).Encode(out)
}