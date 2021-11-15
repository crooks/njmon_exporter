package main

import (
	"fmt"
	"net/http"

	"github.com/Masterminds/log-go"
	"github.com/crooks/jlog"
	"github.com/crooks/njmon_exporter/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cfg   *config.Config
	flags *config.Flags
)

func main() {
	var err error
	flags = config.ParseFlags()
	cfg, err = config.ParseConfig(flags.Config)
	if err != nil {
		log.Fatalf("Unable to parse config file: %v", err)
	}

	if jlog.Enabled() {
		loglevel, err := log.Atoi(cfg.Logging.LevelStr)
		if err != nil {
			log.Fatalf("unable to set log level: %v", err)
		}
		log.Current = jlog.NewJournal(loglevel)
	}
	initCollectors()
	go listener()
	http.Handle("/metrics", promhttp.Handler())
	exporter := fmt.Sprintf("%s:%s", cfg.Exporter.Address, cfg.Exporter.Port)
	log.Infof("Listening for prometheus connections on %s", exporter)
	http.ListenAndServe(exporter, nil)
}
