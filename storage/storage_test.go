package storage

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
)

// test globals
var (
	ts      *Store
	testdb  *sqlx.DB
	testDSN string

	testSchemaPath  = "migrations/schema.sql"
	testDataPath    = "migrations/testdata.sql"
	resetSchemaPath = "migrations/reset.sql"
)

func TestMain(m *testing.M) {
	os.Exit(testMainWrapper(m))
}

func testMainWrapper(m *testing.M) int {
	testDSN = "./test.db"

	defer func() {
		if err := ts.Close(); err != nil {
			log.Fatal(err)
		}
		if err := os.Remove(testDSN); err != nil {
			log.Fatalf("could not remove test db: %v", err)
		}
	}()

	db, err := Open(testDSN)
	if err != nil {
		log.Fatal(err)
	}
	ts = New(db)

	err = ts.MigrateUp(testSchemaPath)
	if err != nil {
		log.Fatal(err)
	}

	// seed test data
	if err := ts.MigrateUp(testDataPath); err != nil {
		log.Fatalf("failed to seed test data: %v", err)
	}
	return m.Run()
}

func TestOpen(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, err := Open(testDSN)
		if err != nil {
			t.Error(err)
		}
		defer db.Close()
	})

	t.Run("no DSN", func(t *testing.T) {
		_, err := Open("")
		if err == nil {
			t.Error("expected error: connection string required")
		}
	})
}

func resetDB() {
	if err := ts.MigrateUp(resetSchemaPath); err != nil {
		log.Fatal(err)
	}
	if err := ts.MigrateUp(testDataPath); err != nil {
		log.Fatal(err)
	}
}

// pretty prints structs for readability
func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func contains(s []string, a string) bool {
	for _, b := range s {
		if a == b {
			return true
		}
	}
	return false
}
