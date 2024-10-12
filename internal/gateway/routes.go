package gateway

import (
	"github.com/justinas/alice"
	"net/http"
)

func Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/health", healthCheckHandler)

	standard := alice.New(recoverPanic, logRequest, commonHeaders)

	return standard.Then(mux)
}
