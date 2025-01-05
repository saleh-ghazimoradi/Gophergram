package handlers

import (
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/helper"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/json"
	"net/http"
)

type HealthHandler struct{}

// Health provides the health status of the application.
//
//	@Summary		Healthcheck
//	@Description	Healthcheck endpoint
//	@Tags			ops
//	@Produce		json
//	@Success		200	{object}	string	"ok"
//	@Router			/v1/health [get]
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     config.AppConfig.ServerConfig.Env,
		"version": config.AppConfig.ServerConfig.Version,
	}
	if err := json.JSONResponse(w, http.StatusOK, data); err != nil {
		helper.InternalServerError(w, r, err)
	}
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}
