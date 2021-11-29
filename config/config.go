package config

import (
	"flag"
	"os"

	"gopkg.in/yaml.v2"
)

// Flags are the command line Flags
type Flags struct {
	Config string
	Debug  bool
}

// Config contains the njmon_exporter configuration data
type Config struct {
	NJmon struct {
		Address string `yaml:"address"`
		Port    string `yaml:"port"`
	} `yaml:"njmon"`
	Logging struct {
		Journal  bool   `yaml:"journal"`
		LevelStr string `yaml:"level"`
	} `yaml:"logging"`
	Exporter struct {
		Address string `yaml:"address"`
		Port    string `yaml:"port"`
	} `yaml:"exporter"`
	InstanceLabel struct {
		Name      string   `yaml:"label_name"`
		Hit       string   `yaml:"label_hit"`
		Miss      string   `yaml:"label_miss"`
		Instances []string `yaml:"hit_instances"`
	} `yaml:"instance_label"`
	Thresholds struct {
		// AliveTimeout is the period of time (in seconds) a host is considered alive without checking in.
		AliveTimeout int `yaml:"alive_timeout"`
		// ConnectionTimeout is the maximum time (in seconds) to wait for data.
		ConnectionTimout int `yaml:"connection_timeout"`
	} `yaml:"thresholds"`
}

// ParseConfig imports a yaml formatted config file into a Config struct
func ParseConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}

// parseFlags processes arguments passed on the command line in the format
// standard format: --foo=bar
func ParseFlags() *Flags {
	f := new(Flags)
	flag.StringVar(&f.Config, "config", "njmon_exporter.yml", "Path to njmon_exporter configuration file")
	flag.BoolVar(&f.Debug, "debug", false, "Expand logging with Debug level messaging and format")
	flag.Parse()
	return f
}
