package rainier

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/unrolled/render"
	"github.com/versolabs/citra/db"
	"github.com/versolabs/citra/util"
	"github.com/versolabs/citra/worker"
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
	db, queries := db.Init(config.DatabaseUrl)

	return &util.Container{
		Asynq:   worker.Client(config.RedisUrl),
		DB:      db,
		Queries: queries,
		Render:  render.New(),
	}
}

func initTestController() *RainierController {
	return &RainierController{
		Container: initTestContainer(),
	}
}
