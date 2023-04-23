package reader

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetaStatus(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	w := httptest.NewRecorder()

	testController := initTestController()
	testController.MetaStatus(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Errorf("%v", err)
	}

	if string(data) != "{\"status\":\"ok\"}" {
		t.Errorf("Invalid response: %v", string(data))
	}
}

func TestMetaPing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	testController := initTestController()
	testController.MetaPing(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Errorf("%v", err)
	}

	if string(data) != "OK" {
		t.Errorf("Invalid response: %v", string(data))
	}
}
