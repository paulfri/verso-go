package reader

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"regexp"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/samber/lo"
	"github.com/unrolled/render"
	"github.com/versolabs/verso/db"
	"github.com/versolabs/verso/util"
	"github.com/versolabs/verso/worker"
)

const snapshotDirectory = "__snapshots__"
const testConfigFile = "./snapshot.toml"

type request struct {
	Method string `toml:"method"`
	Path   string `toml:"path"`
	Auth   bool   `toml:"auth"`
	Body   string `toml:"body"`
	// Query string `toml:"query"` // TODO
}

type test struct {
	Name     string    `toml:"name"`
	Requests []request `toml:"requests"`
}

type config struct {
	Tests []test `toml:"tests"`
}

// Snapshot regression testing for the Reader API.
func TestSnapshot(t *testing.T) {
	var conf config
	_, err := toml.DecodeFile(testConfigFile, &conf)

	if err != nil {
		t.Fatal(err)
	}

	container := initTestContainer()
	router := Router(container)
	server := httptest.NewServer(router)
	defer server.Close()

	makeReq := func(req request) *http.Request {
		r, err := http.NewRequest(req.Method, server.URL+req.Path, strings.NewReader(req.Body))

		if err != nil {
			t.Fatal(err)
		}

		if req.Auth {
			r.Header.Add("Authorization", "GoogleLogin auth=F2vwA2wKSHISLXT7slqt")
		}

		return r
	}

	for _, tt := range conf.Tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Generate requests from configuration.
			reqs := lo.Map(tt.Requests, func(r request, i int) *http.Request {
				return makeReq(r)
			})

			// Execute the configured requests.
			dumps := lo.Map(reqs, dumpRequest)

			// Join the requests into a single snapshot.
			snapshot := strings.Join(dumps, "\n\n")

			// Configure snapshotter with subdirectory.
			snap := cupaloy.New(cupaloy.SnapshotSubdirectory(snapshotDirectory))

			// Create snapshot.
			err = snap.SnapshotWithName(tt.Name, snapshot)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func initTestContainer() *util.Container {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}

	config := util.GetConfig()

	db, queries := db.Init(config.DatabaseURL, false)

	return &util.Container{
		Asynq:     worker.Client(config.RedisURL),
		Config:    &config,
		DB:        db,
		Logger:    util.Logger(),
		Queries:   queries,
		Render:    render.New(),
		Validator: validator.New(),
	}
}

func dumpRequest(req *http.Request, index int) string {
	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	resDump, err := httputil.DumpResponse(res, true)
	if err != nil {
		panic(err)
	}

	// Dump request and response to snapshot output.
	out := fmt.Sprintf("%v\n\n%v", string(reqDump), string(resDump))

	// Stabilize the snapshot.
	out = regexp.MustCompile("Host: .*").ReplaceAllString(out, "Host: localhost:8080")
	out = regexp.MustCompile("Date: .*").ReplaceAllString(out, "Date: Thu, 1 Jan 1970 00:00:00 GMT")

	return out
}
