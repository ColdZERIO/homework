package main

import (
	"context"
	"homework/internal/app"
	"homework/internal/logger"
	postgres "homework/pkg/db"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
)

// done JWT Прописать refresh, добавить роль в claims (проверка ролей и доступа)
// done Удалить логику с куки и вернуть json ответ refresh and access
// done Добавить в колонки user_token Время жизни токена и отозван bool
// done Логику обновления токена в БД поменять на Добавление
// done Дописать Refresh ручку
// done Logger объявить глобально (Не рекомендуют)
// done Переименовать обратно интерфейс слоев, структуры перенести в application
// done Пересмотреть структуру, роутер перенести в internal

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("load env file: %w", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cache := cache.Cache{}

	appLogger := logger.New(logger.Config{
		Level:     slog.LevelInfo,
		AddSource: false,
	})

	db, err := postgres.Init(ctx)
	if err != nil {
		appLogger.ErrorContext(
			ctx, "error init database",
			slog.Any("error", err),
		)
	}

	err = postgres.MigrationRun()
	if err != nil {
		appLogger.ErrorContext(
			ctx, "migration error database",
			slog.Any("error", err),
		)
	}

	router := app.NewApp(db, &cache, appLogger)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	appLogger.Info(
		"Server started",
		slog.String("address", "localhost"),
		slog.String("port", ":8080"),
	)

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		appLogger.ErrorContext(
			ctx, "error listenAndServe",
			slog.Any("error", err),
		)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	<-ctx.Done()

	shutDownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutDownCtx); err != nil {
		appLogger.ErrorContext(
			ctx, "error shutdown server",
			slog.Any("error", err),
		)
	}

	appLogger.Info(
		"server stopped",
		slog.String("port", ":8080"), // ???
	)
}
