package http

import (
	"dusk/http/response"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// handle int64 path parameter
func HandleInt64(key string, rw http.ResponseWriter, r *http.Request) int64 {
	param := chi.URLParam(r, key)

	id, err := strconv.Atoi(param)
	if err != nil {
		response.BadRequest(rw, r, fmt.Errorf("unable to process id: %v", err))
		return -1
	}
	return int64(id)
}
