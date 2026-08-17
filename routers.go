package main

import (
	auth "homework/internal/auth"
	handler "homework/internal/handlers"

	"github.com/go-chi/chi"
)

func routers(handler *handler.UserHandler, tokens *auth.TokenManager) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/login", handler.Login)
	router.Post("/refresh", handler.Refresh)
	router.Post("/logout", tokens.JWTMiddleware(handler.Logout))

	router.Post("/create", tokens.JWTMiddleware(auth.AdminMiddleware(handler.Create)))
	router.Get("/list", tokens.JWTMiddleware(auth.AdminMiddleware(handler.FindUserList)))

	router.Get("/get", tokens.JWTMiddleware(handler.Find))
	router.Delete("/delete", tokens.JWTMiddleware(handler.Delete))
	router.Put("/update", tokens.JWTMiddleware(handler.Update))

	return router
}
