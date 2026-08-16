package main

import (
	"context"
	handler "homework/internal/handlers"
	"homework/internal/services"
	"homework/internal/storage"
	postgres "homework/pkg/db"
	"log"
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
// Добавить логи (в формате json)
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

	db, err := postgres.Init(ctx)
	if err != nil {
		//
	}

	err = postgres.MigrationRun()
	if err != nil {
		//

	}

	// Перенести storage и services в handler
	storage := storage.NewUserStorage(db, &cache)
	service := services.NewUserServices(storage)
	handler := handler.NewUserHandler(service)

	router := routers(handler)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// log start

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
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
		log.Println(err)
	}

	// log stop
}
