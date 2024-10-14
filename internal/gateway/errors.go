package gateway

import (
	"github.com/saleh-ghazimoradi/Gophergram/logger"
	"net/http"
)

func internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	logger.Logger.Errorf("internal server error: %s path: %s error:%s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusInternalServerError, "the server encountered an error")
}

func badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	logger.Logger.Errorf("bad request error: %s path: %s error:%s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	logger.Logger.Errorf("not found error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusNotFound, "not found")
}

func conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	logger.Logger.Errorf("conflict error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusConflict, err.Error())
}
