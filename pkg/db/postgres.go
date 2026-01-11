package db

import (
	"EffectiveMobile/internal/config"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// TODO подумать надо ли разделять от конфига
func NewConnectPostgres(configDB config.PostgresDB, maxAttempts int) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		configDB.User,
		configDB.Pass,
		configDB.Host,
		configDB.Port,
		configDB.Name,
	)
	
	var db *sql.DB
	var lastErr error

	for i := 0; i < maxAttempts; i++ {
		db, lastErr = sql.Open("pgx", dsn)
		if lastErr != nil {
			slog.Warn("failed to open db connection", "attempt", i+1, "error", lastErr)
			time.Sleep(time.Second)
			continue
		}

		// проверяем подключение
		lastErr = db.Ping()
		if lastErr == nil {
			slog.Info("connected to Postgres successfully")
			return db, nil
		}

		slog.Warn("failed to ping db, retrying...", "attempt", i+1, "error", lastErr)
		time.Sleep(time.Second)
	}

	return nil, fmt.Errorf("could not connect to Postgres after %d attempts: %w", maxAttempts, lastErr)
}