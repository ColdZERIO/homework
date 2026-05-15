package main

import (
	"context"
	handler "homework/internal/handlers"
	"homework/internal/services"
	"homework/internal/store"
	postgres "homework/pkg/db"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	router := chi.NewRouter()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := postgres.Init(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = postgres.MigrationRun(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	store := store.NewStore(db)
	srv := services.NewServices(store)
	hand := handler.NewHandler(srv)

	router.Get("/ping", hand.Ping)
	router.Post("/create", hand.CreateUser)
	router.Get("/get/{id}", hand.GetUser)
	router.Delete("/delete/{id}", hand.DeleteUser)
	router.Put("/update", hand.UpdateUser)
	router.Get("/list", hand.GetUsersList)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Println("Server STARTED")

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

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
