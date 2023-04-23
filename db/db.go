package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/versolabs/verso/db/query"
)

func Init(url string) (*sql.DB, *query.Queries) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
	}

	return db, query.New(db)
}
