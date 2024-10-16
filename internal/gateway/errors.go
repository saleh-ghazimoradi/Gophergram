package gateway

import (
	"github.com/saleh-ghazimoradi/Gophergram/logger"
	"net/http"
)

func internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	logger.Logger.Errorw("internal error", "method", r.Method, "path", r.URL.Path, "err", err.Error())
	writeJSONError(w, http.StatusInternalServerError, "the server encountered an error")
}

func badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	logger.Logger.Warnw("bad request", "method", r.Method, "path", r.URL.Path, "err", err.Error())
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	logger.Logger.Warnf("not found error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusNotFound, "not found")
}

func conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	logger.Logger.Errorf("conflict response: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusConflict, err.Error())
}
