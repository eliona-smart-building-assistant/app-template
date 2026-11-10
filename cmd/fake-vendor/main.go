// Copyright (c) 2026 ABB Schweiz AG
// Authors: Adam Lange, et al.

// fake-vendor stands in for the vendor API during local development: GET /devices returns three sensors.
package main

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
)

func main() {
	addr := ":8080"
	if a := os.Getenv("FAKE_VENDOR_ADDR"); a != "" {
		addr = a
	}
	http.HandleFunc("GET /devices", func(w http.ResponseWriter, _ *http.Request) {
		devices := []map[string]any{
			{"id": "sensor-1", "name": "Meeting room", "temperature": 21 + rand.Float64()},
			{"id": "sensor-2", "name": "Open space", "temperature": 22 + rand.Float64()},
			{"id": "sensor-3", "name": "Server room", "temperature": 18 + rand.Float64()},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(devices)
	})
	log.Printf("fake vendor listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
