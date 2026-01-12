package internal

import (
	"EffectiveMobile/internal/config"
	repository "EffectiveMobile/internal/repository/subscription"
	subscriptionRouter "EffectiveMobile/internal/routers/subscription"
	usecase "EffectiveMobile/internal/usecases/subscription"
	"EffectiveMobile/pkg/db"
	"EffectiveMobile/pkg/middleware"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
)

func Run(cfg *config.Config) error {
	// подключение к БД
	connectDB, err := db.NewConnectPostgres(cfg.PostgresDB, cfg.PostgresDB.MaxAttempts)
	if err != nil {
		return err
	}
	defer connectDB.Close()
	slog.Info("Подключение к БД")

	// роутер
	r := mux.NewRouter()

	// Добавляем middleware для логирования входящий запросов и обработки паники
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.RecoverMiddleware)

	// Иницилизируем роутер для логерв
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
        httpSwagger.URL("/swagger/doc.json"), 
    ))

	// Роутер обработки запросов подписок
	repoSubscription := repository.NewRepositorySubscription(connectDB)
	usecaseSubscription := usecase.NewUseCaseSubscription(repoSubscription)
	subscriptionRouter.InitRouter(r, usecaseSubscription)

	// Поднимаем приложение на указанном порту и хосте указанной в конфиге
	slog.Info("Сервер поднялся на порту", "port", cfg.Server.Port)
	if err := http.ListenAndServe(cfg.Server.Host + cfg.Server.Port, r); err != nil {
		return err
	}
	return nil
}
