package http

import (
	"context"
	"net/http"
	"time"

	"github.com/kencx/dusk"
	"github.com/kencx/dusk/api"
	"github.com/kencx/dusk/file"
	"github.com/kencx/dusk/http/response"
	"github.com/kencx/dusk/integration"
	"github.com/kencx/dusk/ui"
	"github.com/kencx/dusk/util"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	idleTimeout      = time.Minute
	readWriteTimeout = 20 * time.Second
	closeTimeout     = 5 * time.Second
)

type Server struct {
	*http.Server
	db       dusk.Store
	fs       *file.Service
	f        integration.Fetchers
	revision string
}

func New(revision string, db dusk.Store, fs *file.Service, f integration.Fetchers) *Server {
	s := &Server{
		Server: &http.Server{
			IdleTimeout:  idleTimeout,
			ReadTimeout:  readWriteTimeout,
			WriteTimeout: readWriteTimeout,
			Handler:      chi.NewRouter(),
		},
		db:       db,
		fs:       fs,
		f:        f,
		revision: revision,
	}
	s.RegisterRoutes()
	return s
}

func (s *Server) Run(port, tlsCert, tlsKey string) error {
	s.Addr = port

	var err error
	if tlsCert != "" && tlsKey != "" {
		err = s.ListenAndServeTLS(tlsCert, tlsKey)
	} else {
		err = s.ListenAndServe()
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Close() error {
	tc, cancel := context.WithTimeout(context.Background(), closeTimeout)
	defer cancel()
	return s.Shutdown(tc)
}

func (s *Server) RegisterRoutes() {
	r := s.Handler.(*chi.Mux)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(timeoutHandler(2 * readWriteTimeout / 3))

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		res, err := util.ToJSON(response.Envelope{
			"timestamp": time.Now().Unix(),
			"message":   "pong",
			"version":   s.revision,
		})
		if err != nil {
			response.InternalServerError(w, r, err)
			return
		}
		response.OK(w, r, res)
	})
	r.Mount("/api", api.Router(s.revision, s.db, s.fs))
	r.Mount("/", ui.Router(s.revision, s.db, s.fs, s.f))
}

// middleware to add http.TimeoutHandler.
func timeoutHandler(timeout time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, timeout, "Timeout")
	}
}
