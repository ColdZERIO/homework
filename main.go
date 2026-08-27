package main

import (
	"context"
	"homework/internal/auth"
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

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

// JWT Прописать refresh, добавить роль в claims (проверка ролей и доступа), сделать отдельный файл ролей и алиас на тип данных Роль (ВЕРНУТЬ КАК БЫЛО)
// Git flow, изучить. 

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
		return
	}

	err = postgres.MigrationRun()
	if err != nil {
		appLogger.ErrorContext(
			ctx, "migration error database",
			slog.Any("error", err),
		)
		return
	}

	tokenManager, err := auth.NewTokenManager(
		os.Getenv("SECRET_KEY"),
		accessTokenTTL,
		refreshTokenTTL,
	)
	if err != nil {
		appLogger.Error("error init token manager", slog.Any("error", err))
		return
	}

	// Перенести storage и services в application
	storage := storage.NewUserStorage(db, &cache, appLogger)
	service := services.NewUserServices(storage, appLogger)
	handler := handler.NewUserHandler(service, tokenManager, appLogger)

	router := routers(handler, tokenManager)

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
