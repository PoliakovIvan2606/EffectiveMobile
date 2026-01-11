package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)


func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        slog.Info(
            fmt.Sprintf("%s %s%s",
            r.Method,
            r.RemoteAddr,
            r.RequestURI),
        )
        next.ServeHTTP(w, r)
    })
}

func RecoverMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if rec := recover(); rec != nil {
                // Логируем панику
                slog.Error("%v\n%s", rec, debug.Stack())
                // Возвращаем 500
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusInternalServerError)
                w.Write([]byte(`{"error": "internal server error"}`))
            }
        }()

        next.ServeHTTP(w, r)
    })
}