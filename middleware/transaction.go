package middleware

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/versolabs/verso/db/query"
	"github.com/versolabs/verso/util"
)

func RunInTransaction(db *sql.DB, queries *query.Queries) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			tx, err := db.BeginTx(r.Context(), nil)
			if err != nil {
				panic(err)
			}

			newQueries := queries.WithTx(tx)
			newReq := requestWithQueries(r, newQueries)

			next.ServeHTTP(w, newReq)

			tx.Rollback()
		}

		return http.HandlerFunc(fn)
	}
}

func requestWithQueries(req *http.Request, queries *query.Queries) *http.Request {
	return req.WithContext(
		context.WithValue(
			req.Context(),
			util.ContextDBQueriesKey{},
			queries,
		),
	)
}
