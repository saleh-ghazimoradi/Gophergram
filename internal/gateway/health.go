package gateway

import (
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"net/http"
)

const version = "0.0.1"

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     config.AppConfig.Env.Env,
		"version": version,
	}
	if err := writeJSON(w, http.StatusOK, data); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
