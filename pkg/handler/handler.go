package handler

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
)

const (
	StatusErrResponse = iota
	StatusOkResponse
)

type ApiOkResponse struct {
    Status  string      `json:"status"`            // "ok" / "error"
    Data    interface{} `json:"data,omitempty"`
}

type ApiErrResponse struct {
    Status  string      `json:"status"`            // "ok" / "error"
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

func NewSqlNullString(str string) sql.NullString {
	if len(str) == 0 {
		return sql.NullString{}
	}

	return sql.NullString{String: str, Valid: true}
}

func NewSqlNullInt64(num int) sql.NullInt64 {
	if num == 0 {
		return sql.NullInt64{}
	}

	return sql.NullInt64{Int64: int64(num), Valid: true}
}