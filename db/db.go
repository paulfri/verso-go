package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/versolabs/citra/db/query"
)

func Queries() *query.Queries {
	database, err := sql.Open("postgres", "user=citra dbname=citra_dev sslmode=disable")

	if err != nil {
		fmt.Println(err)
	}

	return query.New(database)
}
