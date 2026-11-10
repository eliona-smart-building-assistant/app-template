// Copyright (c) 2026 ABB Schweiz AG
// Authors: Adam Lange, et al.

// Package config reads the environment the app is started with.
package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

// Env holds the settings from the environment variables listed in metadata.json.
type Env struct {
	// ApiEndpoint is the BuildingPro Suites API v2 base URL.
	ApiEndpoint string
	// ApiToken is the API key the app runs with. It belongs to exactly one tenant.
	ApiToken string
	// ApiServerPort is the port of the app's own API.
	ApiServerPort int
	// DataDir holds the SQLite database; mount a volume here.
	DataDir string
	// LogLevel is one of debug, info, warn, error.
	LogLevel slog.Level
}

// FromEnv reads the environment and fails on missing or invalid values.
func FromEnv() (Env, error) {
	env := Env{
		ApiEndpoint: os.Getenv("API_ENDPOINT"),
		ApiToken:    os.Getenv("API_TOKEN"),
		DataDir:     getenv("DATA_DIR", "/data"),
	}
	if env.ApiEndpoint == "" {
		return Env{}, fmt.Errorf("API_ENDPOINT is not set")
	}
	if env.ApiToken == "" {
		return Env{}, fmt.Errorf("API_TOKEN is not set")
	}

	port, err := strconv.Atoi(getenv("API_SERVER_PORT", "3000"))
	if err != nil || port <= 0 || port > 65535 {
		return Env{}, fmt.Errorf("API_SERVER_PORT is not a valid port")
	}
	env.ApiServerPort = port

	if err := env.LogLevel.UnmarshalText([]byte(strings.ToUpper(getenv("LOG_LEVEL", "info")))); err != nil {
		return Env{}, fmt.Errorf("LOG_LEVEL: %w", err)
	}
	return env, nil
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
