package app

import (
	"homework/internal/auth"
	handler "homework/internal/handlers"
	"homework/internal/services"
	"homework/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

type Application struct {
	UserHandler *handler.UserHandler
}

func NewApp(db *gorm.DB, cache *cache.Cache, logger *slog.Logger) http.Handler {
	storage := storage.NewUserStorage(db, cache, logger)
	services := services.NewUserServices(storage, logger)
	handler := handler.NewUserHandler(services)

	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Post("/login", handler.Login)
		r.Post("/refresh", handler.RefreshToken)
	})

	router.Group(func(r chi.Router) {
		r.Use(auth.CheckAccessTokenMiddleware)

		r.Post("/create", auth.RequireAdmin(handler.Create))
		r.Get("/get", handler.Find)
		r.Delete("/delete", handler.Delete)
		r.Put("/update", handler.Update)
		r.Get("/list", auth.RequireAdmin(handler.FindUserList))
	})

	return router
}
