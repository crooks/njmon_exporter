package main

import (
	"log"
	"time"
)

// lastSeen keeps a tally of known hosts and when they were last seen.
type lastSeen map[string]time.Time

// upTest runs an endless loop, iterating over all the known hosts and
// populating an "up" metric.  The metric returns 1 if a host has been seen
// within an acceptable period of time or 0 is it hasn't.
func (ls lastSeen) upTest() {
	timeout := time.Second * time.Duration(cfg.AliveTimeout)
	for {
		now := time.Now().UTC()
		for hostname, t := range ls {
			if now.Sub(t) > timeout {
				// Host is considered down
				hostUp.WithLabelValues(hostname).Set(float64(0))
			} else {
				// Host is up
				hostUp.WithLabelValues(hostname).Set(float64(1))
			}
		} // End of hosts loop
		time.Sleep(60 * time.Second)
	} // Endless loop
}

func (ls lastSeen) registerHost(hostname string) {
	if _, seen := ls[hostname]; !seen {
		// This host hasn't been seen before.
		log.Printf("%s: New host discovered", hostname)
	}
	ls[hostname] = time.Now().UTC()
}
