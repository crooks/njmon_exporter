package main

import (
	"log"
	"time"
)

// hostInfo contains the fields to be recorded against each discovered
// hostname.
type hostInfo struct {
	lastSeen time.Time
	labelVal string
}

// hostInfoMap declares a map of structs (hostInfo), keyed by a string
// (hostname).
type hostInfoMap map[string]*hostInfo

//strContains returns true if string s is a member of slice l
func strContains(l []string, s string) bool {
	for _, a := range l {
		if a == s {
			return true
		}
	}
	return false
}

// upTest runs an endless loop, iterating over all the known hosts and
// populating an "up" metric.  The metric returns 1 if a host has been seen
// within an acceptable period of time or 0 if it hasn't.
func (h hostInfoMap) upTest() {
	if cfg.AliveTimeout < 60 {
		cfg.AliveTimeout = 60
		log.Printf("Setting a sane Alive Timeout of %d seconds", cfg.AliveTimeout)
	}
	timeout := time.Second * time.Duration(cfg.AliveTimeout)
	log.Printf("Host considered dead if not seen for %d seconds", int(timeout.Seconds()))
	// Wait a little while on startup to give hosts a chance to check in.
	time.Sleep(120 * time.Second)
	for {
		now := time.Now().UTC()
		for hostname, t := range h {
			if now.Sub(t.lastSeen) > timeout {
				// Host is considered down
				hostUp.WithLabelValues(hostname, t.labelVal).Set(float64(0))
			} else {
				// Host is up
				log.Printf("Labelling %s as %s", hostname, t.labelVal)
				hostUp.WithLabelValues(hostname, t.labelVal).Set(float64(1))
			}
		} // End of hosts loop
		time.Sleep(60 * time.Second)
	} // Endless loop
}

// registerHost takes a hostname and returns a Hit/Miss label.  If the
// hostname is known, the time it was last seen is updated to Now.  If it's
// unknown, the new host is registered and its Hit/Miss status recorded.
func (h hostInfoMap) registerHost(hostname string) string {
	var n *hostInfo
	if _, seen := h[hostname]; !seen {
		log.Printf("New host discovered: %s", hostname)
		n = new(hostInfo)
		if strContains(cfg.InstanceLabel.Instances, hostname) {
			n.labelVal = cfg.InstanceLabel.Hit
		} else {
			n.labelVal = cfg.InstanceLabel.Miss
		}
	} else {
		n = h[hostname]
	}
	h[hostname] = n
	n.lastSeen = time.Now().UTC()
	return n.labelVal
}
