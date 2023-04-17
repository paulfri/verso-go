package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	"github.com/versolabs/citra/db/query"
)

func Queries() *query.Queries {
	database, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		panic(err)
	}

	return query.New(database)
}
