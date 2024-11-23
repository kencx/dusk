package request

import (
	"fmt"
	"net/http"
	"net/url"
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

func QueryString(qv url.Values, key string, defaultValue string) string {
	if !qv.Has(key) {
		return defaultValue
	}
	return qv.Get(key)
}

func QueryInt(qv url.Values, key string, defaultValue int) int {
	if !qv.Has(key) {
		return defaultValue
	}
	value := qv.Get(key)

	i, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return i
}

func HasValue(params url.Values, value string) bool {
	if params.Has(value) {
		return params.Get(value) != ""
	}
	return false
}

func HasOptionalValue(params url.Values, value string) bool {
	return params.Has(value)
}
