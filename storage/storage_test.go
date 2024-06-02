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

	testSchemaPath  = "schema.sql"
	resetSchemaPath = "reset.sql"
)

func TestMain(m *testing.M) {
	os.Exit(testMainWrapper(m))
}

func testMainWrapper(m *testing.M) int {
	testDSN = "./test.db"

	defer func() {
		if err := ts.Close(); err != nil {
			log.Print(err)
		}
		if err := os.Remove(testDSN); err != nil {
			log.Fatalf("could not remove test db: %v", err)
		}
	}()

	db, err := Open(testDSN)
	if err != nil {
		log.Print(err)
	}
	ts = New(db)

	err = ts.MigrateUp(testSchemaPath)
	if err != nil {
		log.Print(err)
	}

	// seed test data
	if err := seedData(); err != nil {
		log.Printf("failed to seed test data: %v", err)
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
		log.Print(err)
	}

	if err := seedData(); err != nil {
		log.Print(err)
	}
}

// pretty prints structs for readability
func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
