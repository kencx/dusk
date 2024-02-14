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
	"syscall"

	"github.com/kencx/dusk/file"
	dhttp "github.com/kencx/dusk/http"
	"github.com/kencx/dusk/storage"
)

type config struct {
	port    int
	dsn     string
	tlsCert string
	tlsKey  string
	dataDir string
}

func main() {
	var config config

	flag.IntVar(&config.port, "port", 9090, "Server Port")
	flag.StringVar(&config.dsn, "dsn", "library.db", "sqlite DSN")
	flag.StringVar(&config.tlsCert, "tlsKey", "", "TLS Certificate path")
	flag.StringVar(&config.tlsKey, "tlsCert", "", "TLS Key path")
	flag.StringVar(&config.dataDir, "dataDir", "dusk_data", "Data directory")
	flag.Parse()

	db, err := storage.Open(config.dsn)
	if err != nil {
		log.Fatal(err)
	}

	store := storage.New(db)
	err = store.MigrateUp("schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	err = store.MigrateUp("testdata.sql")
	if err != nil {
		slog.Error("Migration step failed", slog.Any("err", err))
	}

	fw, err := file.NewWorker(config.dataDir)
	if err != nil {
		log.Fatal(err)
	}

	srv := dhttp.New(store, fw)
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
