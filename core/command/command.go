package command

import (
	"database/sql"

	"github.com/hibiken/asynq"
	"github.com/versolabs/verso/db/query"
)

type Command struct {
	Asynq   *asynq.Client
	DB      *sql.DB
	Queries *query.Queries
}

// Yield a Queries struct that operates within a transaction.
func (c Command) QueryTransaction(callback func(*query.Queries) error) error {
	tx, err := c.DB.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	withTx := c.Queries.WithTx(tx)

	err = callback(withTx)

	if err != nil {
		return err
	}

	err = tx.Commit()

	return err
}
