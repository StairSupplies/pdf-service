package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

// Auth enforces bearer-token authentication when PDF_SERVICE_API_KEY is set.
// When the env var is empty or unset, all requests pass through unchanged.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := os.Getenv("PDF_SERVICE_API_KEY")
		if key == "" {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		token, found := strings.CutPrefix(authHeader, "Bearer ")
		if !found || token != key {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"}) //nolint:errcheck
			return
		}

		next.ServeHTTP(w, r)
	})
}
