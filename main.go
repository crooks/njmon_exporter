package main

import (
	"fmt"
	"net/http"

	"github.com/Masterminds/log-go"
	"github.com/crooks/jlog"
	loglevel "github.com/crooks/log-go-level"
	"github.com/crooks/njmon_exporter/config"
	"github.com/crooks/strmatch"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cfg     *config.Config
	flags   *config.Flags
	exclude *strmatch.Matcher
)

func initExclude() {
	exclude = strmatch.NewMatcher()
	for _, s := range cfg.Exclude.Str {
		err := exclude.SAdd(s)
		if err != nil {
			log.Warnf("%s: Unable to add exclusion string: %v", s, err)
			continue
		}
	}
	for _, r := range cfg.Exclude.Regex {
		err := exclude.RAdd(r)
		if err != nil {
			log.Warnf("%s: Unable to add exclusion Regular Expression: %v", r, err)
			continue
		}
	}
}

func main() {
	var err error
	flags = config.ParseFlags()
	cfg, err = config.ParseConfig(flags.Config)
	if err != nil {
		log.Fatalf("Unable to parse config file: %v", err)
	}

	// Define logging level and method
	loglev, err := loglevel.ParseLevel(cfg.Logging.LevelStr)
	if err != nil {
		log.Fatalf("unable to set log level: %v", err)
	}
	if cfg.Logging.Journal && jlog.Enabled() {
		log.Current = jlog.NewJournal(loglev)
	} else {
		log.Current = log.StdLogger{Level: loglev}
	}

	initExclude()
	initCollectors()
	go listener()
	http.Handle("/metrics", promhttp.Handler())
	exporter := fmt.Sprintf("%s:%s", cfg.Exporter.Address, cfg.Exporter.Port)
	log.Infof("Listening for prometheus connections on %s", exporter)
	http.ListenAndServe(exporter, nil)
}
