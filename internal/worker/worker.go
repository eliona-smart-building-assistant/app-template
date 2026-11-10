// Copyright (c) 2026 ABB Schweiz AG
// Authors: Adam Lange, et al.

// Package worker polls the vendor for every enabled configuration and writes the result to BuildingPro Suites.
package worker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"app-name/internal/platform"
	"app-name/internal/store"
	"app-name/internal/vendor"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v3"
)

// AssetType is the name of the asset type in resources/asset-types.
const AssetType = "app_name_device"

// Platform is the part of the platform client the worker needs; tests replace it.
type Platform interface {
	UpsertAsset(projectID, gai, assetType, name string) (int32, error)
	PutData(assetID int32, subtype api.DataSubtype, values map[string]any) error
}

var _ Platform = (*platform.Client)(nil)

type Worker struct {
	store    store.Store
	platform Platform
	log      *slog.Logger
	lastRun  map[int64]time.Time
}

func New(s store.Store, p Platform, log *slog.Logger) *Worker {
	return &Worker{store: s, platform: p, log: log, lastRun: map[int64]time.Time{}}
}

// Run checks every second which configuration is due, until ctx is cancelled.
func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		w.tick(ctx, time.Now())
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (w *Worker) tick(ctx context.Context, now time.Time) {
	configs, err := w.store.ListConfigs(ctx)
	if err != nil {
		w.log.Error("list configurations", "err", err)
		return
	}
	for _, c := range configs {
		if !c.Enable || now.Sub(w.lastRun[c.ID]) < time.Duration(c.RefreshInterval)*time.Second {
			continue
		}
		w.lastRun[c.ID] = now
		if err := w.Sync(ctx, c); err != nil {
			w.log.Error("sync", "config", c.ID, "err", err)
		}
	}
}

// Sync reads the devices of one configuration and creates or updates their assets and data in every project.
func (w *Worker) Sync(ctx context.Context, c store.Config) error {
	client := vendor.NewClient(c.URL, c.APIKey, time.Duration(c.RequestTimeout)*time.Second)
	devices, err := client.Devices(ctx)
	if err != nil {
		return err
	}
	for _, projectID := range c.ProjectIDs {
		for _, d := range devices {
			assetID, err := w.assetID(ctx, c.ID, projectID, d)
			if err != nil {
				return err
			}
			if err := w.platform.PutData(assetID, api.INPUT, map[string]any{"temperature": d.Temperature}); err != nil {
				return fmt.Errorf("data of device %s: %w", d.ID, err)
			}
		}
	}
	w.log.Debug("synced", "config", c.ID, "devices", len(devices), "projects", len(c.ProjectIDs))
	return nil
}

// assetID returns the known asset of a device, or creates it on first sight and remembers the mapping.
func (w *Worker) assetID(ctx context.Context, configID int64, projectID string, d vendor.Device) (int32, error) {
	if id, ok, err := w.store.AssetID(ctx, configID, projectID, d.ID); err != nil || ok {
		return id, err
	}
	gai := fmt.Sprintf("%s_%d_%s", AssetType, configID, d.ID)
	id, err := w.platform.UpsertAsset(projectID, gai, AssetType, d.Name)
	if err != nil {
		return 0, fmt.Errorf("asset of device %s: %w", d.ID, err)
	}
	if err := w.store.SetAssetID(ctx, configID, projectID, d.ID, id); err != nil {
		return 0, err
	}
	w.log.Info("asset created", "config", configID, "project", projectID, "device", d.ID, "asset", id)
	return id, nil
}
