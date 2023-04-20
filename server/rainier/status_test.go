package rainier

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/unrolled/render"
	"github.com/versolabs/citra/db"
	"github.com/versolabs/citra/server/utils"
	"github.com/versolabs/citra/tasks"
)

func initTestContainer() *utils.Container {
	db, queries := db.Init()

	return &utils.Container{
		Asynq:   tasks.Client(),
		Render:  render.New(),
		DB:      db,
		Queries: queries,
	}
}

func initTestController() *RainierController {
	return &RainierController{
		Container: initTestContainer(),
	}
}

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
