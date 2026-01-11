package main

import (
	"EffectiveMobile/internal"
	"EffectiveMobile/internal/config"
	"EffectiveMobile/pkg/db"
	"EffectiveMobile/pkg/logger"
	"fmt"
	"log/slog"
	"os"
)

var (
	pathDirMigrations = "./migrations"
	pathConfig = "configs/config.yaml"
)

func main() {
	cfg, err := config.Init(pathConfig)
	if err != nil {
		slog.Error("получение конфига", "error", err)
	}

	wd, _ := os.Getwd()
	slog.Info("workdir", "dir", wd)

	logger.InitLogger(cfg.LogLevel)
	fmt.Println(cfg)

	if err := db.RunMigrations(pathDirMigrations, cfg.PostgresDB.GetDSN(), cfg.PostgresDB.MaxAttempts); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if err := internal.Run(cfg); err != nil {
		slog.Error("ошибка запуска сервера", "eor", err)
		os.Exit(1)
	}
}