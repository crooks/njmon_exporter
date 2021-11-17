package main

import (
	"sync"
	"time"

	"github.com/Masterminds/log-go"
)

// hostInfo contains the fields to be recorded against each discovered
// hostname.
type hostInfo struct {
	lastSeen   time.Time
	labelVal   string
	alertState bool
}

// hostInfoMap declares struct containing a map of structs (hostInfo), keyed by a string (hostname).
// It also includes a RWMutex mu to prevent read/write race conditions on the map.
type hostInfoMap struct {
	hosts map[string]*hostInfo
	mu    sync.Mutex
}

func newHostInfoMap() *hostInfoMap {
	h := new(hostInfoMap)
	h.hosts = make(map[string]*hostInfo)
	return h
}

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
func (h *hostInfoMap) upTest() {
	// If the AliveTimeout is too short, it will mark hosts dead in between
	// njmon checkins.  Even 60 seconds is a bit bonkers.
	if cfg.Thresholds.AliveTimeout < 60 {
		cfg.Thresholds.AliveTimeout = 60
		log.Warnf("Setting a sane Alive Timeout of %d seconds", cfg.Thresholds.AliveTimeout)
	}
	timeout := time.Second * time.Duration(cfg.Thresholds.AliveTimeout)
	log.Infof("Host considered dead if not seen for %d seconds", int(timeout.Seconds()))
	// Wait a little while on startup to give hosts a chance to check in.
	for {
		now := time.Now().UTC()
		h.mu.Lock()
		for hostname, t := range h.hosts {
			if now.Sub(t.lastSeen) > timeout {
				// Host is considered down
				hostUp.WithLabelValues(hostname, t.labelVal).Set(float64(0))
				// If the host was previously not alerting, this is a state change
				if !t.alertState {
					log.Warnf("%s: host is now dead", hostname)
					t.alertState = true
				}
			} else {
				// Host is up
				hostUp.WithLabelValues(hostname, t.labelVal).Set(float64(1))
				// If the host was previously alerting, this is a state change
				if t.alertState {
					log.Infof("%s: host has transitioned to up", hostname)
					t.alertState = false
				}
			}
		} // End of hosts loop
		h.mu.Unlock()
		time.Sleep(60 * time.Second)
	} // Endless loop
}

// registerHost takes a hostname and returns a Hit/Miss label.  If the
// hostname is known, the time it was last seen is updated to Now.  If it's
// unknown, the new host is registered and its Hit/Miss status recorded.
func (h *hostInfoMap) registerHost(hostname string) string {
	var n *hostInfo
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, seen := h.hosts[hostname]; !seen {
		n = new(hostInfo)
		if strContains(cfg.InstanceLabel.Instances, hostname) {
			n.labelVal = cfg.InstanceLabel.Hit
		} else {
			n.labelVal = cfg.InstanceLabel.Miss
		}
		log.Infof("New host discovered: hostname=%s, %s=%s", hostname, cfg.InstanceLabel.Name, n.labelVal)
	} else {
		n = h.hosts[hostname]
	}
	n.lastSeen = time.Now().UTC()
	h.hosts[hostname] = n
	return n.labelVal
}
