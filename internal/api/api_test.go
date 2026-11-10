// Copyright (c) 2026 ABB Schweiz AG
// Authors: Adam Lange, et al.

package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	appfiles "app-name"
	"app-name/internal/apiserver"
	"app-name/internal/store"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	s, err := store.OpenSQLite(context.Background(), t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	srv, err := NewServer(s, appfiles.OpenAPI)
	if err != nil {
		t.Fatal(err)
	}
	h := apiserver.HandlerWithOptions(apiserver.NewStrictHandler(srv, nil), apiserver.StdHTTPServerOptions{BaseURL: "/v1"})
	ts := httptest.NewServer(h)
	t.Cleanup(ts.Close)
	return ts
}

func do(t *testing.T, method, url, body string) (int, string) {
	t.Helper()
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(b)
}

func TestConfigurationAPI(t *testing.T) {
	ts := newTestServer(t)

	code, body := do(t, http.MethodPost, ts.URL+"/v1/configs", `{"url":"https://vendor.example","apiKey":"secret","projectIDs":["1"]}`)
	if code != http.StatusCreated {
		t.Fatalf("create: %d %s", code, body)
	}
	if strings.Contains(body, "secret") {
		t.Fatalf("the API key must never be returned: %s", body)
	}
	var created apiserver.Configuration
	if err := json.Unmarshal([]byte(body), &created); err != nil {
		t.Fatal(err)
	}
	if *created.RefreshInterval != 60 || !*created.Enable {
		t.Fatalf("defaults not applied: %s", body)
	}

	code, body = do(t, http.MethodPut, ts.URL+"/v1/configs/1", `{"url":"https://vendor.example","enable":false}`)
	if code != http.StatusOK || !strings.Contains(body, `"enable":false`) {
		t.Fatalf("update: %d %s", code, body)
	}

	if code, body = do(t, http.MethodPost, ts.URL+"/v1/configs", `{"url":"not a url"}`); code != http.StatusBadRequest {
		t.Fatalf("invalid url: %d %s", code, body)
	}
	if code, body = do(t, http.MethodPost, ts.URL+"/v1/configs", `{"url":"https://vendor.example","refreshInterval":1}`); code != http.StatusBadRequest {
		t.Fatalf("too short refreshInterval: %d %s", code, body)
	}

	if code, _ = do(t, http.MethodDelete, ts.URL+"/v1/configs/1", ""); code != http.StatusNoContent {
		t.Fatalf("delete: %d", code)
	}
	if code, _ = do(t, http.MethodGet, ts.URL+"/v1/configs/1", ""); code != http.StatusNotFound {
		t.Fatalf("get after delete: %d", code)
	}
}

func TestVersionAndSpec(t *testing.T) {
	ts := newTestServer(t)
	if code, body := do(t, http.MethodGet, ts.URL+"/v1/version", ""); code != http.StatusOK || !strings.Contains(body, "commit") {
		t.Fatalf("version: %d %s", code, body)
	}
	code, body := do(t, http.MethodGet, ts.URL+"/v1/version/openapi.json", "")
	if code != http.StatusOK || !strings.Contains(body, `"openapi":"3.0.3"`) {
		t.Fatalf("openapi.json: %d %.80s", code, body)
	}
}
