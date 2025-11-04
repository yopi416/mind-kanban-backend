package middleware

import (
	"net/http"

	"github.com/yopi416/mind-kanban-backend/configs"
)

func ApplyCORS(next http.Handler, cfg *configs.ConfigList) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", cfg.CorsAllowOrigins) //フロントエンドURL
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, Accept, Origin, Authorization")
		// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // Cookie許可

		// Preflight (OPTIONS) 対応
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
