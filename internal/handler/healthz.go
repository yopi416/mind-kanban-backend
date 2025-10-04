package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yopi416/mind-kanban-backend/internal/model"
)

// A HealthzHandler implements health check endpoint.
type HealthzHandler struct{}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	response := &model.HealthzResponse{Message: "OK"}
	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		log.Println("healthz encode error:", err)
		return
	}
}
