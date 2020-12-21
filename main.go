package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
)

// Config contains the njmon_exporter configuration data
type Config struct {
	NJmon struct {
		Address string `yaml:"address"`
		Port    string `yaml:"port"`
	} `yaml:"njmon"`
	Exporter struct {
		Address string `yaml:"address"`
		Port    string `yaml:"port"`
	} `yaml:"exporter"`
}

// newConfig imports a yaml formatted config file into a Config struct
func newConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	config := &Config{}
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}

// Create a global configuration
var cfg *Config

func main() {
	var err error
	cfg, err = newConfig("njmon_exporter.yml")
	if err != nil {
		fmt.Println(err)
	}
	go listener()
	http.Handle("/metrics", promhttp.Handler())
	exporter := fmt.Sprintf("%s:%s", cfg.Exporter.Address, cfg.Exporter.Port)
	http.ListenAndServe(exporter, nil)
}
