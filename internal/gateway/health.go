package gateway

import (
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"net/http"
)

const Version = ""

//	 healthCheckHandler godoc
//
//		@Summary		Healthcheck
//		@Description	Healthcheck endpoint
//		@Tags			ops
//		@Produce		json
//		@Success		200	{object}	string	"ok"
//		@Router			/health [get]
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     config.AppConfig.Env.Env,
		"version": Version,
	}
	if err := jsonResponse(w, http.StatusOK, data); err != nil {
		internalServerError(w, r, err)
	}
}
