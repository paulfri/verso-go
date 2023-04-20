package rainier

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
)

func authenticatedTestRequest(method string, path string, _body io.Reader, token string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/reader/api/0/token", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, ContextAuthTokenKey{}, token)
	req = req.WithContext(ctx)
	return req
}
