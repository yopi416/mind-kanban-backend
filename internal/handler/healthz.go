package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/yopi416/mind-kanban-backend/api"
)

func (s *Server) GetHealthz(w http.ResponseWriter, r *http.Request) {
	// 追加属性は必要に応じてWith
	lg := slog.Default().With("handler", "GetHealthz")

	// レスポンスヘッダを一応明示
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := api.Healthz{Message: "health check OK"}
	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		lg.Error("healthz encode error", "err", err)
		return
	}

	// 正常時は Info ログを出力
	lg.Info("health check ok")
}
