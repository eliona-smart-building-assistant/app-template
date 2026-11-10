// Copyright (c) 2026 ABB Schweiz AG
// Authors: Adam Lange, et al.

// Package store keeps the app's own data locally under DATA_DIR: its configurations. The app never connects to the
// BuildingPro Suites database; everything in BuildingPro Suites goes through API v2, including the mapping of
// vendor devices to assets (global asset identifier).
package store

import (
	"context"
	"errors"
)

// ErrNotFound is returned when a configuration does not exist.
var ErrNotFound = errors.New("not found")

// Config is one connection to a vendor account.
type Config struct {
	ID              int64
	URL             string
	APIKey          string
	Enable          bool
	RefreshInterval int32 // seconds
	RequestTimeout  int32 // seconds
	ProjectIDs      []string
}

// Store is everything the app persists. Swap the implementation without touching the rest of the app.
type Store interface {
	ListConfigs(ctx context.Context) ([]Config, error)
	GetConfig(ctx context.Context, id int64) (Config, error)
	CreateConfig(ctx context.Context, c Config) (Config, error)
	UpdateConfig(ctx context.Context, c Config) (Config, error)
	DeleteConfig(ctx context.Context, id int64) error

	Close() error
}
