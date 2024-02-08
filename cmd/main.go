package main

import (
	dhttp "dusk/http"
	"dusk/storage"
	"dusk/worker"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	flag.StringVar(&config.tlsCert, "TLS cert", "", "TLS Certificate path")
	flag.StringVar(&config.tlsKey, "TLS key", "", "TLS Key path")
	flag.StringVar(&config.dataDir, "data_dir", "data", "Data directory")
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
		log.Print(err)
	}

	fw, err := worker.NewFileWorker(config.dataDir)
	if err != nil {
		log.Fatal(err)
	}

	srv := dhttp.New(store, fw)
	go func() error {
		// srv.InfoLog.Printf("Starting server on :%d", config.port)
		err := srv.Run(fmt.Sprintf(":%d", config.port), config.tlsCert, config.tlsKey)
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	_ = <-sig
	// srv.InfoLog.Printf("Received signal %s, shutting down...", s.String())

	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
	// srv.InfoLog.Println("Database connection closed")

	if srv != nil {
		if err := srv.Close(); err != nil {
			log.Fatal(err)
		}
		// srv.InfoLog.Println("Server connection closed")
	}
	// srv.InfoLog.Println("Application gracefully stopped")
}
