package main

import (
	handler "homework/internal/handlers"

	"github.com/go-chi/chi"
)

func routers(hand *handler.Handler) *chi.Mux {
	rout := chi.NewRouter()

	rout.Get("/ping", hand.Ping)
	rout.Post("/persist", hand.Persist)
	rout.Get("/get", hand.Find)
	rout.Delete("/delete", hand.Delete)
	rout.Put("/update", hand.Update)
	rout.Get("/list", hand.GetList)

	return rout
}
