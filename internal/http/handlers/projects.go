package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bladimirbalbin/portafolio-api/internal/repository/postgres"
	"github.com/go-chi/chi/v5"
)

func ListProjects(repo *postgres.ProjectRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var featured *bool
		if v := r.URL.Query().Get("featured"); v != "" {
			b, err := strconv.ParseBool(v)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]any{"error": "invalid featured parameter (use true/false)"})
				return
			}
			featured = &b
		}

		var tag *string
		if v := r.URL.Query().Get("tag"); v != "" {
			tag = &v
		}

		limit := 20
		if v := r.URL.Query().Get("limit"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]any{"error": "invalid limit (integer)"})
				return
			}
			limit = n
		}

		offset := 0
		if v := r.URL.Query().Get("offset"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]any{"error": "invalid offset (integer)"})
				return
			}
			offset = n
		}

		sort := r.URL.Query().Get("sort") // recent|featured

		items, err := repo.List(r.Context(), postgres.ProjectListOpts{
			Featured: featured,
			Tag:      tag,
			Limit:    limit,
			Offset:   offset,
			Sort:     sort,
		})
		if err != nil {
			// Error de sort inválido lo devolvemos como 400
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
			return
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"data":   items,
			"limit":  clamp(limit, 1, 100),
			"offset": max(0, offset),
			"sort":   defaultSort(sort),
		})
	}
}

func GetProjectBySlug(repo *postgres.ProjectRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		slug := chi.URLParam(r, "slug")
		if slug == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]any{"error": "missing slug"})
			return
		}

		p, err := repo.GetBySlug(r.Context(), slug)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]any{"error": "internal error"})
			return
		}
		if p == nil {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]any{"error": "project not found"})
			return
		}

		_ = json.NewEncoder(w).Encode(map[string]any{"data": p})
	}
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func defaultSort(s string) string {
	if s == "" {
		return "featured"
	}
	return s
}
