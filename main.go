package main

import (
	"context"
	handler "homework/internal/handlers"
	"homework/internal/logger"
	"homework/internal/services"
	"homework/internal/storage"
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

// Дополнительно проверка в мидлваре UUID
// JWT Прописать refresh, добавить роль в claims (проверка ролей и доступа)
// done Добавить логи (в формате json)
// Хеширование поменять на bcrypt (добавить соль)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("load env file: %w", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cache := cache.Cache{} // Перенос в Ручки

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

	// Перенести storage и services в handler
	storage := storage.NewUserStorage(db, &cache, appLogger)
	service := services.NewUserServices(storage, appLogger)
	handler := handler.NewUserHandler(service, appLogger)

	router := routers(handler)

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
