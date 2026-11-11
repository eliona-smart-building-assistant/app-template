// Copyright (c) 2026 ABB Schweiz AG
// Authors: Adam Lange, et al.

package store

import (
	"context"
	"errors"
	"testing"
)

func openTest(t *testing.T) *SQLite {
	t.Helper()
	s, err := OpenSQLite(context.Background(), t.TempDir())
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestConfigLifecycle(t *testing.T) {
	ctx := context.Background()
	s := openTest(t)

	created, err := s.CreateConfig(ctx, Config{URL: "http://vendor", APIKey: "secret", Enable: true, RefreshInterval: 60, RequestTimeout: 30})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("create returned no ID")
	}

	got, err := s.GetConfig(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.URL != "http://vendor" || got.APIKey != "secret" {
		t.Fatalf("get returned %+v", got)
	}

	got.Enable = false
	if _, err := s.UpdateConfig(ctx, got); err != nil {
		t.Fatalf("update: %v", err)
	}
	all, err := s.ListConfigs(ctx)
	if err != nil || len(all) != 1 || all[0].Enable {
		t.Fatalf("list after update: %+v, %v", all, err)
	}

	if err := s.DeleteConfig(ctx, created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := s.GetConfig(ctx, created.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("get after delete: %v", err)
	}
	if err := s.DeleteConfig(ctx, created.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("second delete: %v", err)
	}
}

func TestMigrationsAreIdempotent(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	for i := 0; i < 2; i++ {
		s, err := OpenSQLite(ctx, dir)
		if err != nil {
			t.Fatalf("open %d: %v", i, err)
		}
		s.Close()
	}
}
