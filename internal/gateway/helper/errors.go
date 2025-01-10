package helper

import (
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/json"
	"github.com/saleh-ghazimoradi/Gophergram/logger"
	"net/http"
)

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	logger.Logger.Error("internal error", "method", r.Method, "path", r.URL.Path, "err", err.Error())
	json.WriteJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	logger.Logger.Warn("bad request", "method", r.Method, "path", r.URL.Path, "err", err.Error())
	json.WriteJSONError(w, http.StatusBadRequest, err.Error())
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	logger.Logger.Warn("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	json.WriteJSONError(w, http.StatusNotFound, err.Error())
}

func ConflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	logger.Logger.Error("conflict response", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	json.WriteJSONError(w, http.StatusConflict, err.Error())
}

func UnauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logger.Logger.Warn("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	json.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func UnauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logger.Logger.Warn("unauthorized basic error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	json.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func ForbiddenResponse(w http.ResponseWriter, r *http.Request) {
	logger.Logger.Warn("forbidden", "method", r.Method, "path", r.URL.Path)
	json.WriteJSONError(w, http.StatusForbidden, "forbidden")
}

func RateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	logger.Logger.Warn("rate limit exceeded", "method", r.Method, "path", r.URL.Path)
	w.Header().Set("Retry_After", retryAfter)
	json.WriteJSONError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}
