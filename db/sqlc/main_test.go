package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	driver  = "postgres"
	connUrl = "postgresql://root:secret@localhost:5432/go_forum?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	testDB, err = sql.Open(driver, connUrl)
	if err != nil {
		log.Fatal("Cannot connect to DB ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
