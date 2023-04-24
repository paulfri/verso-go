package db

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/versolabs/verso/db/query"
	"github.com/versolabs/verso/util"
)

func Init(url string) (*sql.DB, *query.Queries) {
	logger := util.Logger()

	if logger == nil {
		panic("Logger failed to initialize")
	}

	db, err := sql.Open("postgres", url)
	if err != nil {
		logger.Fatal().Err(err)
	}

	// loggerAdapter := zerologadapter.New(*logger)
	// db = sqldblogger.OpenDriver(url, db.Driver(), loggerAdapter)

	err = db.Ping()
	if err != nil {
		logger.Fatal().Err(err)
	}

	return db, query.New(db)
}
