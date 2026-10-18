package main

import (
	auth "homework/internal/auth"
	handler "homework/internal/handlers"

	"github.com/go-chi/chi"
)

func routers(handler *handler.UserHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/create", auth.AuthMiddleware(handler.Create)) // Админ создает пользователя
	router.Post("/login", handler.Login)
	router.Get("/get", auth.AuthMiddleware(handler.Find))
	router.Delete("/delete", auth.AuthMiddleware(handler.Delete))
	router.Put("/update", auth.AuthMiddleware(handler.Update))
	router.Get("/list", auth.AuthMiddleware(handler.FindUserList))

	return router
}
