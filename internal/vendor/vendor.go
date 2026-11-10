// Copyright (c) 2026 ABB Schweiz AG
// Authors: Adam Lange, et al.

// Package vendor talks to the system the app integrates. Replace it with the real vendor API.
package vendor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Device is one thing the vendor reports. Its ID must be stable: it maps the device to its asset.
type Device struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Temperature float64 `json:"temperature"`
}

// Client reads devices from the vendor API.
type Client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

func NewClient(baseURL, apiKey string, timeout time.Duration) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		http:    &http.Client{Timeout: timeout},
	}
}

// Devices calls GET {baseURL}/devices with the API key as bearer token.
func (c *Client) Devices(ctx context.Context) ([]Device, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/devices", nil)
	if err != nil {
		return nil, err
	}
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("vendor request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vendor request: %s", resp.Status)
	}
	var devices []Device
	if err := json.NewDecoder(resp.Body).Decode(&devices); err != nil {
		return nil, fmt.Errorf("vendor response: %w", err)
	}
	return devices, nil
}
