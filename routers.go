package main

import (
	middleware "homework/internal/Middleware"
	handler "homework/internal/handlers"

	"github.com/go-chi/chi"
)

func routers(hand *handler.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)

		r.Get("/get", hand.Find)
		r.Put("/update", hand.Update)
		r.Get("/list", hand.GetList)
		r.Delete("/delete", hand.Delete)
	})

	r.Get("/ping", hand.Ping)
	r.Post("/persist", hand.Persist)
	// Make login

	return r
}
