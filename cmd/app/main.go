package main

import (
	_ "EffectiveMobile/docs"
	"EffectiveMobile/internal"
	"EffectiveMobile/internal/config"
	"EffectiveMobile/pkg/db"
	"EffectiveMobile/pkg/logger"
	"log/slog"
	"os"
)

var (
	pathDirMigrations = "./migrations"
	pathConfig = "configs/config.yaml"
)

func main() {
	// Конфиг получаем либо из env если переменных нет получаем из .yaml файла
	cfg, err := config.Init(pathConfig)
	if err != nil {
		slog.Error("получение конфига", "error", err)
	}

	// Иницилизируем логгер
	logger.InitLogger(cfg.LogLevel)

	// Создаём таблицы
	if err := db.RunMigrations(pathDirMigrations, cfg.PostgresDB.GetDSN(), cfg.PostgresDB.MaxAttempts); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// Запускаем приложение
	if err := internal.Run(cfg); err != nil {
		slog.Error("ошибка запуска сервера", "eor", err)
		os.Exit(1)
	}
}