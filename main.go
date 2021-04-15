package main

import (
	"database/sql"
	"log"

	"github.com/ashishkhuraishy/go_forum/api"
	db "github.com/ashishkhuraishy/go_forum/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	driver  = "postgres"
	connUrl = "postgresql://root:secret@localhost:5432/go_forum?sslmode=disable"

	addreess = "localhost:8080"
)

func main() {
	conn, err := sql.Open(driver, connUrl)
	if err != nil {
		log.Fatal("Cannot connect to DB ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(store)
	if err != nil {
		panic(err)
	}

	if err := server.Start(addreess); err != nil {
		log.Fatal(err)
	}
}
