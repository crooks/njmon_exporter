package config

import (
	"testing"
)

func TestConfig(t *testing.T) {
	cfg, err := ParseConfig("njmon_exporter.yml")
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

func TestFlags(t *testing.T) {
	f := ParseFlags()
	expectingConfig := "njmon_exporter.yml"
	if f.Config != expectingConfig {
		t.Fatalf("Unexpected config flag: Expected=%s, Got=%s", expectingConfig, f.Config)
	}
	if f.Debug {
		t.Fatal("Unexpected debug flag: Expected=false, Got=true")
	}
}
