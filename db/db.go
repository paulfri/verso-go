package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/versolabs/citra/db/query"
)

func Init(databaseUrl string) (*sql.DB, *query.Queries) {
	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	return db, query.New(db)
}
