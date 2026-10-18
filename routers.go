package main

import (
	auth "homework/internal/auth"
	handler "homework/internal/handlers"

	"github.com/go-chi/chi"
)

func routers(handler *handler.UserHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/create", auth.CheckAccessTokenMiddleware(handler.Create)) // Админ создает пользователя
	router.Post("/login", handler.Login)
	router.Get("/get", auth.CheckAccessTokenMiddleware(handler.Find))
	router.Delete("/delete", auth.CheckAccessTokenMiddleware(handler.Delete))
	router.Put("/update", auth.CheckAccessTokenMiddleware(handler.Update))
	router.Get("/list", auth.CheckAccessTokenMiddleware(handler.FindUserList))

	return router
}
