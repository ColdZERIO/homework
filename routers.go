package main

import (
	auth "homework/internal/auth"
	handler "homework/internal/handlers"

	"github.com/go-chi/chi"
)

func routers(handler *handler.UserHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/create", auth.JWTMiddleware(handler.Create)) // Админ создает пользователя
	router.Post("/login", handler.Login)
	router.Get("/get", auth.JWTMiddleware(handler.Find))
	router.Delete("/delete", auth.JWTMiddleware(handler.Delete))
	router.Put("/update", auth.JWTMiddleware(handler.Update))
	router.Get("/list", auth.JWTMiddleware(handler.FindUserList))

	return router
}
