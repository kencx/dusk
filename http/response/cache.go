package response

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"hash"
	"io"
	"log/slog"
	"net/http"
)

type etagResponse struct {
	http.ResponseWriter
	buf  bytes.Buffer
	hash hash.Hash
	w    io.Writer
}

func (e *etagResponse) Write(p []byte) (int, error) {
	return e.w.Write(p)
}

func ETag(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ew := &etagResponse{
			ResponseWriter: rw,
			buf:            bytes.Buffer{},
			hash:           sha1.New(),
		}

		ew.w = io.MultiWriter(&ew.buf, ew.hash)

		next.ServeHTTP(ew, r)

		sum := fmt.Sprintf("%x", ew.hash.Sum(nil))
		rw.Header().Set("ETag", sum)

		if r.Header.Get("If-None-Match") == sum {
			rw.WriteHeader(304)
		} else {
			_, err := ew.buf.WriteTo(rw)
			if err != nil {
				slog.Error("failed to write response", slog.Any("err", err))
			}
		}
	})
}

func SetCache(duration int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", duration))
			next.ServeHTTP(rw, r)
		})
	}
}

func NoCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Cache-Control", "no-cache")
		next.ServeHTTP(rw, r)
	})
}
