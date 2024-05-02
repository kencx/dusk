package request

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/kencx/dusk/http/response"

	"github.com/go-chi/chi/v5"
)

func FetchIdFromSlug(rw http.ResponseWriter, r *http.Request) int64 {
	param := chi.URLParam(r, "slug")
	paramSlice := strings.Split(param, "-")

	if len(paramSlice) > 1 {
		idStr := paramSlice[len(paramSlice)-1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.BadRequest(rw, r, fmt.Errorf("unable to fetch id from slug: %w", err))
			return -1
		}
		return int64(id)
	}
	response.BadRequest(rw, r, fmt.Errorf("invalid slug"))
	return -1
}

func HandleInt64(key string, rw http.ResponseWriter, r *http.Request) int64 {
	param := chi.URLParam(r, key)

	id, err := strconv.Atoi(param)
	if err != nil {
		response.BadRequest(rw, r, fmt.Errorf("unable to process id: %w", err))
		return -1
	}
	return int64(id)
}
