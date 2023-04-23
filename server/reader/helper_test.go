package reader

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/unrolled/render"
	"github.com/versolabs/verso/db"
	"github.com/versolabs/verso/util"
	"github.com/versolabs/verso/worker"
)

func authenticatedTestRequest(method string, path string, _body io.Reader, token string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/reader/api/0/token", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, ContextAuthTokenKey{}, token)
	req = req.WithContext(ctx)
	return req
}

func initTestContainer() *util.Container {
	config := util.GetConfig()
	db, queries := db.Init(config.DatabaseURL)

	return &util.Container{
		Asynq:   worker.Client(config.RedisURL),
		DB:      db,
		Queries: queries,
		Render:  render.New(),
	}
}

func initTestController() *ReaderController {
	return &ReaderController{
		Container: initTestContainer(),
	}
}