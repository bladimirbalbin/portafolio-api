package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Health(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		dbStatus := "disabled"
		if db != nil {
			ctx, cancel := context.WithTimeout(r.Context(), 800*time.Millisecond)
			defer cancel()
			if err := db.Ping(ctx); err != nil {
				dbStatus = "down"
			} else {
				dbStatus = "up"
			}
		}

		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
			"db":     dbStatus,
		})
	}
}
