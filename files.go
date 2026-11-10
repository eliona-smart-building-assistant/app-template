// Copyright (c) 2026 ABB Schweiz AG
// Authors: Adam Lange, et al.

// Package appfiles embeds the files the app needs at runtime, so the binary runs without them on disk.
package appfiles

import "embed"

//go:embed metadata.json
var Metadata []byte

//go:embed openapi.yaml
var OpenAPI []byte

//go:embed resources/asset-types/*.json
var AssetTypes embed.FS
