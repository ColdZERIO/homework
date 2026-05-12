package main

import (
	"context"
	handler "homework/internal/handlers"
	"homework/internal/services"
	"homework/internal/store"
	postgres "homework/pkg/db"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	router := chi.NewRouter()

	ctx := context.Background()
	db, err := postgres.Init(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	store := store.NewStore(db)
	srv := services.NewServices(store)
	hand := handler.NewHandler(srv)

	log.Println("Server STARTED")
	router.Get("/ping", hand.Ping)
	router.Post("/create", hand.CreateUser)
	router.Get("/get/{id}", hand.GetUser)
	router.Delete("/delete/{id}", hand.DeleteUser)
	router.Put("/update", hand.UpdateUser)
	router.Get("/list", hand.GetUsersList)

	log.Fatal(http.ListenAndServe(":8080", router))
}
