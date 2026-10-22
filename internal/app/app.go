package app

import (
	auth "homework/internal/auth"
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

	router.Post("/create", auth.CheckAccessTokenMiddleware(handler.Create)) // Админ создает пользователя
	router.Post("/login", handler.Login)
	router.Get("/get", auth.CheckAccessTokenMiddleware(handler.Find))
	router.Delete("/delete", auth.CheckAccessTokenMiddleware(handler.Delete))
	router.Put("/update", auth.CheckAccessTokenMiddleware(handler.Update))
	router.Get("/list", auth.CheckAccessTokenMiddleware(handler.FindUserList))

	return router
}
