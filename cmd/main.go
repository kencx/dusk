package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"github.com/kencx/dusk/file"
	dhttp "github.com/kencx/dusk/http"
	"github.com/kencx/dusk/integration/googlebooks"
	"github.com/kencx/dusk/storage"
)

var version string

const sqliteDB = "library.db"

type config struct {
	port     int
	lib      string
	tlsCert  string
	tlsKey   string
	logLevel string
}

func main() {
	var config config

	flag.IntVar(&config.port, "port", 9090, "Server port")
	flag.StringVar(&config.lib, "lib", "lib", "Path to library directory")
	flag.StringVar(&config.tlsCert, "tlsKey", "", "TLS certificate path")
	flag.StringVar(&config.tlsKey, "tlsCert", "", "TLS key path")
	flag.StringVar(&config.logLevel, "log", "info", "Log level")
	flag.Parse()

	level := new(slog.LevelVar)
	switch strings.ToLower(config.logLevel) {
	case "debug":
		level.Set(slog.LevelDebug)
	case "info":
		level.Set(slog.LevelInfo)
	case "warn":
		level.Set(slog.LevelWarn)
	case "err", "error":
		level.Set(slog.LevelError)
	}

	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(l)

	fw, err := file.NewService(config.lib)
	if err != nil {
		log.Fatal(err)
	}

	// TODO allow multiple fetchers
	fetcher := new(googlebooks.Fetcher)

	dsn := path.Join(config.lib, sqliteDB)
	db, err := storage.Open(dsn)
	if err != nil {
		log.Fatal(err)
	}

	store := storage.New(db)
	err = store.MigrateUp("schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	// err = store.MigrateUp("testdata.sql")
	// if err != nil {
	// 	slog.Error("Migration step failed", slog.Any("err", err))
	// }

	srv := dhttp.New(version, store, fw, fetcher)
	go func() error {
		slog.Info(fmt.Sprintf("Starting server on port %d", config.port))
		err := srv.Run(fmt.Sprintf(":%d", config.port), config.tlsCert, config.tlsKey)
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	slog.Info(fmt.Sprintf("Received signal %s, shutting down...", s.String()))

	if err := store.Close(); err != nil {
		log.Fatal(err)
	}
	slog.Info("Database connection closed")

	if srv != nil {
		if err := srv.Close(); err != nil {
			log.Fatal(err)
		}
		slog.Info("Server connection closed")
	}
	slog.Info("Application gracefully stopped")
}
