// Copyright (c) 2026 ABB Schweiz AG
// Authors: Adam Lange, et al.

// Package platform writes to BuildingPro Suites through API v2. The app has no other access to the platform.
package platform

import (
	"encoding/json"
	"fmt"
	"io/fs"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v3"
	"github.com/eliona-smart-building-assistant/go-eliona/v2/asset"
)

// Client holds the API endpoint and the app's API key.
type Client struct {
	endpoint string
	apiKey   string
}

func NewClient(endpoint, apiKey string) *Client {
	return &Client{endpoint: endpoint, apiKey: apiKey}
}

// InitAssetTypes creates or updates every asset type found as JSON file in files.
func (c *Client) InitAssetTypes(files fs.FS) error {
	paths, err := fs.Glob(files, "resources/asset-types/*.json")
	if err != nil {
		return err
	}
	for _, p := range paths {
		raw, err := fs.ReadFile(files, p)
		if err != nil {
			return err
		}
		var t api.AssetType
		if err := json.Unmarshal(raw, &t); err != nil {
			return fmt.Errorf("asset type %s: %w", p, err)
		}
		if err := asset.UpsertAssetType(c.endpoint, c.apiKey, t); err != nil {
			return fmt.Errorf("asset type %s: %w", p, err)
		}
	}
	return nil
}

// UpsertAsset creates the asset or updates it by its global asset identifier and returns its ID. The asset
// belongs to the tenant of the API key.
func (c *Client) UpsertAsset(gai, assetType, name string) (int32, error) {
	a := api.NewAsset(gai, assetType)
	a.SetName(name)
	id, err := asset.UpsertAsset(c.endpoint, c.apiKey, *a)
	if err != nil {
		return 0, err
	}
	if id == nil {
		return 0, fmt.Errorf("asset %s: API returned no ID", gai)
	}
	return *id, nil
}

// PutData writes the current values of an asset.
func (c *Client) PutData(assetID int32, subtype api.DataSubtype, values map[string]any) error {
	return asset.UpsertData(c.endpoint, c.apiKey, *api.NewData(assetID, subtype, values))
}
