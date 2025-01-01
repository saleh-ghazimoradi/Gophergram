package handlers

import (
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/json"
	"net/http"
)

type HealthHandler struct{}

func (h *HealthHandler) Health(w http.ResponseWriter) {
	data := map[string]string{
		"status":  "ok",
		"env":     config.AppConfig.ServerConfig.Env,
		"version": config.AppConfig.ServerConfig.Version,
	}
	if err := json.WriteJSON(w, http.StatusOK, data); err != nil {
		json.WriteJSONError(w, http.StatusInternalServerError, err.Error())
	}
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}
