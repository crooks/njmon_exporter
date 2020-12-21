package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-yaml/yaml"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config contains the njmon_exporter configuration data
type Config struct {
	Listen struct {
		address string
		port    int
	} `yaml:"listen"`
}

func readConfig(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	fmt.Println(cfg)
	if err != nil {
		return err
	}
	return nil
}

var cfg Config

func main() {
	err := readConfig("njmon_exporter.yml")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cfg)
	go listener()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
