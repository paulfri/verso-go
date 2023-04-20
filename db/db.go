package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	"github.com/versolabs/citra/db/query"
)

func Init() (*sql.DB, *query.Queries) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		panic(err)
	}

	return db, query.New(db)
}
