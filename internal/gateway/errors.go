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

func unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logger.Logger.Warnf("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logger.Logger.Warnf("unauthorized basic error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	logger.Logger.Warnw("forbidden", "method", r.Method, "path", r.URL.Path, "error")
	writeJSONError(w, http.StatusForbidden, "forbidden")
}

func rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	logger.Logger.Warnw("rate limit exceeded", "method", r.Method, "path", r.URL.Path)
	w.Header().Set("Retry_After", retryAfter)
	writeJSONError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}
