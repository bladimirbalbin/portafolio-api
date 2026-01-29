package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/bladimirbalbin/portafolio-api/internal/auth"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	if req.Username == os.Getenv("ADMIN_USER") &&
		req.Password == os.Getenv("ADMIN_PASS") {

		token, _ := auth.GenerateToken(req.Username)
		json.NewEncoder(w).Encode(map[string]string{"token": token})
		return
	}

	http.Error(w, "invalid credentials", http.StatusUnauthorized)
}
