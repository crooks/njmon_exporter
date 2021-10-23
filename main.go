package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/crooks/njmon_exporter/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cfg   *config.Config
	flags *config.Flags
)

type logWriter struct{}

// writer creates a new Write method for the logger
func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(string(bytes))
}

// initLogging initiates the custom logWriter if the debug flag isn't set.  This makes logging more systemd friendly by
// chopping out the date/time prefix.
func initLogging() {
	if !flags.Debug {
		log.SetFlags(0)
		log.SetOutput(new(logWriter))
	}
}

func main() {
	var err error
	flags = config.ParseFlags()
	cfg, err = config.ParseConfig(flags.Config)
	if err != nil {
		log.Fatalf("Unable to parse config file: %v", err)
	}
	initLogging()

	initCollectors()
	go listener()
	http.Handle("/metrics", promhttp.Handler())
	exporter := fmt.Sprintf("%s:%s", cfg.Exporter.Address, cfg.Exporter.Port)
	log.Printf("Listening for prometheus connections on %s", exporter)
	http.ListenAndServe(exporter, nil)
}
