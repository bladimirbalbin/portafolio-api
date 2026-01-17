package http

import (
	"net/http"
	"time"

	"github.com/bladimirbalbin/portafolio-api/internal/config"
	"github.com/bladimirbalbin/portafolio-api/internal/http/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/bladimirbalbin/portafolio-api/internal/repository/postgres"
)

func NewRouter(cfg config.Config, db *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))

	r.Get("/health", handlers.Health(db))
	projectRepo := postgres.NewProjectRepo(db)
	r.Get("/projects", handlers.ListProjects(projectRepo))

	return r
}
