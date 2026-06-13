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
	// hand -> handler
	handler := handler.UserHandler(service)

	router := routers(handler)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Server STARTED")

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
		// поменять на Print
		log.Println(err)
	}

	log.Println("Server STOPPED")
}
