package rainier

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestToken(t *testing.T) {
	token := "foo"

	req := authenticatedTestRequest(http.MethodGet, "/reader/api/0/token", nil, token)
	w := httptest.NewRecorder()

	testController := initTestController()
	testController.UserToken(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Errorf("%v", err)
	}

	if string(data) != token {
		t.Errorf("Invalid response: %v", string(data))
	}
}

func TestUserPreferences(t *testing.T) {
	req := authenticatedTestRequest(http.MethodGet, "/reader/api/0/preference/list", nil, "token")
	w := httptest.NewRecorder()

	testController := initTestController()
	testController.UserPreferences(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Errorf("%v", err)
	}

	if string(data) != UserPreferencesResponse {
		t.Errorf("Invalid response: %v", string(data))
	}
}

func TestUserStreamPreferences(t *testing.T) {
	req := authenticatedTestRequest(http.MethodGet, "/reader/api/0/preference/stream/list", nil, "token")
	w := httptest.NewRecorder()

	testController := initTestController()
	testController.UserStreamPreferences(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Errorf("%v", err)
	}

	if string(data) != UserStreamPreferencesResponse {
		t.Errorf("Invalid response: %v", string(data))
	}
}
