package gateway

import (
	"net/http"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("status: available"))
}
