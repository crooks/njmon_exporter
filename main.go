package main

import (
	"flag"
	"fmt"
	"log"
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
	AliveTimeout int `yaml:"alive_timeout"`
}

var (
	cfg            *Config
	flagConfigFile string // Fully-qualified path to config file
)

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

// parseFlags processes arguments passed on the command line in the format
// standard format: --foo=bar
func parseFlags() {
	flag.StringVar(
		&flagConfigFile,
		"config",
		"njmon_exporter.yml",
		"Path to njmon_exporter configuration file",
	)
	flag.Parse()
}

func main() {
	var err error
	parseFlags()
	cfg, err = newConfig(flagConfigFile)
	if err != nil {
		log.Fatalf("Unable to parse config file: %v", err)
	}
	go listener()
	http.Handle("/metrics", promhttp.Handler())
	exporter := fmt.Sprintf("%s:%s", cfg.Exporter.Address, cfg.Exporter.Port)
	log.Printf("Listening for prometheus connections on %s", exporter)
	http.ListenAndServe(exporter, nil)
}
