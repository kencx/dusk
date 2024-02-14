package request

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/kencx/dusk/http/response"

	"github.com/go-chi/chi/v5"
)

func HandleInt64(key string, rw http.ResponseWriter, r *http.Request) int64 {
	param := chi.URLParam(r, key)

	id, err := strconv.Atoi(param)
	if err != nil {
		response.BadRequest(rw, r, fmt.Errorf("unable to process id: %w", err))
		return -1
	}
	return int64(id)
}
