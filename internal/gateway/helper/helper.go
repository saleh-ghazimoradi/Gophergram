package helper

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func ReadIdParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}
