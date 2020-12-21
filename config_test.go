package main

import (
	"testing"
)

func TestConfig(t *testing.T) {
	cfg, err := newConfig("njmon_exporter.yml")
	if err != nil {
		t.Fatalf("Failed with: %v", err)
	}
	if cfg.Exporter.Address != "127.0.0.1" {
		t.Fatalf("Expected=127.0.0.1, Got=%s", cfg.Exporter.Address)
	}
	if cfg.Exporter.Port != "9772" {
		t.Fatalf("Expected=9772, Got=%s", cfg.Exporter.Port)
	}
}
