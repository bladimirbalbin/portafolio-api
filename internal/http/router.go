package http

import (
	"net/http"
	"time"

	"github.com/bladimirbalbin/portafolio-api/internal/config"
	"github.com/bladimirbalbin/portafolio-api/internal/http/handlers"
	appmw "github.com/bladimirbalbin/portafolio-api/internal/http/middleware"
	"github.com/bladimirbalbin/portafolio-api/internal/repository/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(cfg config.Config, db *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))
	// --------PUBLIC --------
	r.Get("/health", handlers.Health(db))

	projectRepo := postgres.NewProjectRepo(db)
	r.Get("/projects", handlers.ListProjects(projectRepo))
	r.Get("/projects/", handlers.ListProjects(projectRepo))
	r.Get("/projects/{slug}", handlers.GetProjectBySlug(projectRepo))
	r.Get("/docs/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/openapi.json")
	})

	r.Post("/auth/login", handlers.Login)

	// -------- PROTECTED --------
	r.Group(func(pr chi.Router) {
		pr.Use(appmw.JWT)

		pr.Post("/projects", handlers.CreateProject(projectRepo))
		pr.Put("/projects/{slug}", handlers.UpdateProject(projectRepo))
		pr.Delete("/projects/{slug}", handlers.DeleteProject(projectRepo))
	})

	return r
}
