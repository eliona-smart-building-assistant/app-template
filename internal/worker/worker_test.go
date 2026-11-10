// Copyright (c) 2026 ABB Schweiz AG
// Authors: Adam Lange, et al.

package worker

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"app-name/internal/store"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v3"
)

type fakePlatform struct {
	assets map[string]int32 // gai -> asset ID
	data   map[int32]map[string]any
}

func (f *fakePlatform) UpsertAsset(projectID, gai, assetType, name string) (int32, error) {
	if id, ok := f.assets[gai+"@"+projectID]; ok {
		return id, nil
	}
	id := int32(len(f.assets) + 1)
	f.assets[gai+"@"+projectID] = id
	return id, nil
}

func (f *fakePlatform) PutData(assetID int32, subtype api.DataSubtype, values map[string]any) error {
	f.data[assetID] = values
	return nil
}

func TestSyncCreatesAssetsOnceAndWritesData(t *testing.T) {
	ctx := context.Background()
	vendorAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		io.WriteString(w, `[{"id":"a","name":"Sensor A","temperature":21.5},{"id":"b","name":"Sensor B","temperature":19}]`)
	}))
	defer vendorAPI.Close()

	s, err := store.OpenSQLite(ctx, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	c, err := s.CreateConfig(ctx, store.Config{URL: vendorAPI.URL, APIKey: "key", Enable: true, RefreshInterval: 60, RequestTimeout: 5, ProjectIDs: []string{"1", "2"}})
	if err != nil {
		t.Fatal(err)
	}

	p := &fakePlatform{assets: map[string]int32{}, data: map[int32]map[string]any{}}
	w := New(s, p, slog.New(slog.NewTextHandler(io.Discard, nil)))
	for i := 0; i < 2; i++ {
		if err := w.Sync(ctx, c); err != nil {
			t.Fatalf("sync %d: %v", i, err)
		}
	}

	if len(p.assets) != 4 {
		t.Fatalf("want 4 assets (2 devices x 2 projects), got %d", len(p.assets))
	}
	id, ok, err := s.AssetID(ctx, c.ID, "2", "a")
	if err != nil || !ok {
		t.Fatalf("mapping of device a in project 2: ok=%v err=%v", ok, err)
	}
	if got := p.data[id]["temperature"]; got != 21.5 {
		t.Fatalf("temperature of device a: %v", got)
	}
}

func TestSyncReportsVendorError(t *testing.T) {
	ctx := context.Background()
	vendorAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer vendorAPI.Close()

	s, err := store.OpenSQLite(ctx, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	p := &fakePlatform{assets: map[string]int32{}, data: map[int32]map[string]any{}}
	w := New(s, p, slog.New(slog.NewTextHandler(io.Discard, nil)))

	if err := w.Sync(ctx, store.Config{ID: 1, URL: vendorAPI.URL, RequestTimeout: 5, ProjectIDs: []string{"1"}}); err == nil {
		t.Fatal("want an error for a rejected vendor request")
	}
	if len(p.assets) != 0 {
		t.Fatalf("no asset may be created on error, got %d", len(p.assets))
	}
}
