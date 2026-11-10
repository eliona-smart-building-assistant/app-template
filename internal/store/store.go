// Copyright (c) 2026 ABB Schweiz AG
// Authors: Adam Lange, et al.

// Package store keeps the app's own data: its configurations and the mapping of vendor IDs to BuildingPro Suites assets.
// The app never connects to the BuildingPro Suites database. Everything in BuildingPro Suites goes through API v2.
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

	// AssetID returns the BuildingPro Suites asset of a vendor device in a project; ok is false if it is not known yet.
	AssetID(ctx context.Context, configID int64, projectID, vendorID string) (assetID int32, ok bool, err error)
	SetAssetID(ctx context.Context, configID int64, projectID, vendorID string, assetID int32) error

	Close() error
}
