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
)

// Объединить функцию Persist и Update
// Найти информацию об id в DB UUID
// Создать на каждую ручку свою стрктуру
// Создать 3 роли При создании пользователей
// Переименовать название интерфейсов (более информативнее)

// ТЕОРИЯ
// Изучить Статус Коды
// HTTP и HTTPS как работает
// Патерны/антипатерны архитектуры, ооп
// Индексы в DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("load env file: %w", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := postgres.Init(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = postgres.MigrationRun()
	if err != nil {
		log.Fatal(err)
		return
	}

	// перенести стор и сервисы в хэндлер
	storage := storage.UserStorage(db)
	service := services.UserServices(storage)
	handler := handler.UserHandler(service)

	router := routers(handler)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Server STARTED")

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Println(err)
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
		log.Fatal(err)
	}

	log.Println("Server STOPPED")
}


