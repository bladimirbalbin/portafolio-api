package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bladimirbalbin/portafolio-api/internal/repository/postgres"
)

func ListProjects(repo *postgres.ProjectRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var featured *bool
		if v := r.URL.Query().Get("featured"); v != "" {
			b, err := strconv.ParseBool(v)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"error": "invalid featured parameter (use true/false)",
				})
				return
			}
			featured = &b
		}

		var tag *string
		if v := r.URL.Query().Get("tag"); v != "" {
			tag = &v
		}

		items, err := repo.List(r.Context(), featured, tag)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]any{"error": "internal error"})
			return
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": items,
		})
	}
}
