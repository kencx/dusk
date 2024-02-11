package storage

import (
	"embed"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var migrationFs embed.FS

type Store struct {
	db *sqlx.DB
}

func Open(path string) (*sqlx.DB, error) {
	if path == "" {
		return nil, fmt.Errorf("db: connection string required")
	}

	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("db: failed to open: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("db: failed to connect: %w", err)
	}
	return db, nil
}

func New(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return errors.New("db: database does not exist")
}

func (s *Store) MigrateUp(filePath string) error {
	schema, err := migrationFs.ReadFile(fmt.Sprintf("migrations/%s", filePath))
	if err != nil {
		return fmt.Errorf("db: cannot read sql file %q: %w", filePath, err)
	}
	if _, err := s.db.Exec(string(schema)); err != nil {
		return fmt.Errorf("db: failed to execute sql file %q: %w", filePath, err)
	}

	slog.Info("Database schema loaded", "schema", filePath)
	return nil
}
