package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bladimirbalbin/portafolio-api/internal/domain"
	"github.com/bladimirbalbin/portafolio-api/internal/repository/postgres"
	"github.com/go-chi/chi/v5"
)

type projectPayload struct {
	Slug        string   `json:"slug"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	RepoURL     *string  `json:"repo_url"`
	DemoURL     *string  `json:"demo_url"`
	Tags        []string `json:"tags"`
	Featured    bool     `json:"featured"`
}

func CreateProject(repo *postgres.ProjectRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req projectPayload
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if req.Slug == "" || req.Name == "" || req.Description == "" {
			http.Error(w, "slug, name, description are required", http.StatusBadRequest)
			return
		}
		if req.Tags == nil {
			req.Tags = []string{}
		}

		created, err := repo.Create(r.Context(), domain.Project{
			Slug:        req.Slug,
			Name:        req.Name,
			Description: req.Description,
			RepoURL:     req.RepoURL,
			DemoURL:     req.DemoURL,
			Tags:        req.Tags,
			Featured:    req.Featured,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(created)
	}
}

func UpdateProject(repo *postgres.ProjectRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		if slug == "" {
			http.Error(w, "missing slug", http.StatusBadRequest)
			return
		}

		var req projectPayload
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if req.Name == "" || req.Description == "" {
			http.Error(w, "name and description are required", http.StatusBadRequest)
			return
		}
		if req.Tags == nil {
			req.Tags = []string{}
		}

		updated, err := repo.UpdateBySlug(r.Context(), slug, domain.Project{
			Name:        req.Name,
			Description: req.Description,
			RepoURL:     req.RepoURL,
			DemoURL:     req.DemoURL,
			Tags:        req.Tags,
			Featured:    req.Featured,
		})
		if err != nil {
			if err == postgres.ErrNotFound {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(updated)
	}
}

func DeleteProject(repo *postgres.ProjectRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		if slug == "" {
			http.Error(w, "missing slug", http.StatusBadRequest)
			return
		}

		if err := repo.DeleteBySlug(r.Context(), slug); err != nil {
			if err == postgres.ErrNotFound {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
