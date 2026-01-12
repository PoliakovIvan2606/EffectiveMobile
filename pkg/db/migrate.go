package db

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(pathDirMigrations, dsn string, maxAttempts int) error {
	var lastErr error

	for i := 0; i < maxAttempts; i++ {
		lastErr = func() error {
			m, err := migrate.New(
				fmt.Sprintf("file://%s", pathDirMigrations),
				dsn,
			)
			if err != nil {
				return fmt.Errorf("migrate init: %w", err)
			}

			if err := m.Up(); err != nil {
				if err == migrate.ErrNoChange {
					return nil
				}
				return fmt.Errorf("migrate up: %w", err)
			}

			slog.Info("database migrations applied successfully")
			return nil
		}()

		if lastErr == nil {
			return nil
		}

		slog.Warn("migration attempt failed, retrying...", "attempt", i+1, "error", lastErr)
		time.Sleep(time.Second)
	}

	return fmt.Errorf("all migration attempts failed: %w", lastErr)
}