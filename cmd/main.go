package main

import (
	dhttp "dusk/http"
	"dusk/storage"
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
	port int
	dsn  string
}

func main() {
	var config config

	flag.IntVar(&config.port, "port", 9090, "Server Port")
	flag.StringVar(&config.dsn, "dsn", "library.db", "sqlite DSN")
	flag.Parse()

	db, err := storage.Open(config.dsn)
	if err != nil {
		log.Fatal(err)
	}

	store := storage.New(db)
	err = store.MigrateUp("../storage/migrations/schema.sql")
	if err != nil {
		log.Fatal(err)
	}

	srv := dhttp.New(store)
	go func() error {
		srv.InfoLog.Printf("Starting server on :%d", config.port)
		err := srv.Run(fmt.Sprintf(":%d", config.port))
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	srv.InfoLog.Printf("Received signal %s, shutting down...", s.String())

	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
	srv.InfoLog.Println("Database connection closed")

	if srv != nil {
		if err := srv.Close(); err != nil {
			log.Fatal(err)
		}
		srv.InfoLog.Println("Server connection closed")
	}
	srv.InfoLog.Println("Application gracefully stopped")
}
