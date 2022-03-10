package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/Masterminds/log-go"
	"github.com/tidwall/gjson"
)

const (
	mb   = float64(1048576)
	page = float64(4096)
)

// listener listens for connections from njmon.  It forks handleConnection() for each connection.
func listener() {
	njmonListen := fmt.Sprintf("%s:%s", cfg.NJmon.Address, cfg.NJmon.Port)
	// Listen for incoming connections.
	l, err := net.Listen("tcp", njmonListen)
	if err != nil {
		log.Fatalf("Unable to start njmon listener: %v", err)
	}
	// Close the listener when the application closes.
	defer l.Close()
	// hosts provides a means to populate a simple "up" metric to indicate if a
	// previously seen host has stopped submitting data.
	h := newHostInfoMap()
	go h.upTest()

	// An endless loop of listening for incoming njmon connections
	log.Infof("Listening for njmon connections on %s\n", njmonListen)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Warnf("Unable to accept njmon connection: %v", err)
			continue
		}
		// Handle connections in a new goroutine.
		go h.connHandle(conn)
	}
}

func disks(hostname, instanceLabel string, result gjson.Result) {
	for device, k := range result.Map() {
		// vg = Volume Group
		vg := k.Get("vg").String()
		// The following metrics need converting from MB to Bytes
		size := k.Get("size_mb").Float() * mb
		free := k.Get("free_mb").Float() * mb
		read := k.Get("read_mbps").Float() * mb
		write := k.Get("write_mbps").Float() * mb
		diskBlockSize.WithLabelValues(hostname, instanceLabel, device, vg).Set(k.Get("blocksize").Float())
		diskBusy.WithLabelValues(hostname, instanceLabel, device, vg).Set(k.Get("busy").Float())
		diskFree.WithLabelValues(hostname, instanceLabel, device, vg).Set(free)
		diskRead.WithLabelValues(hostname, instanceLabel, device, vg).Set(read)
		diskSize.WithLabelValues(hostname, instanceLabel, device, vg).Set(size)
		diskWrite.WithLabelValues(hostname, instanceLabel, device, vg).Set(write)
	}
}

// filesystems iterates over the content of the njmon filesystems data and
// produces a set of metrics for each filesystem.
func filesystems(hostname, instanceLabel string, result gjson.Result) {
	for _, f := range result.Map() {
		// Required labels for filesystems
		device := f.Get("device").String()
		mount := f.Get("mount").String()
		// Filesystem metrics
		size := f.Get("size_mb").Float() * mb
		free := f.Get("free_mb").Float() * mb
		filesystemSize.WithLabelValues(hostname, instanceLabel, device, mount).Set(size)
		filesystemFree.WithLabelValues(hostname, instanceLabel, device, mount).Set(free)
	}
}

// netAdapters iterates over the content of the njmon network_adapters data and
// produces a set of metrics for each interface.
func netAdapters(hostname, instanceLabel string, result gjson.Result) {
	for i, f := range result.Map() {
		// Undecided whether to use Bits or Bytes for network interfaces.
		// Bytes will do for now.
		netBpsRx.WithLabelValues(hostname, instanceLabel, i).Set(f.Get("rx_bytes").Float())
		netBpsTx.WithLabelValues(hostname, instanceLabel, i).Set(f.Get("tx_bytes").Float())
		netPktRx.WithLabelValues(hostname, instanceLabel, i).Set(f.Get("rx_packets").Float())
		netPktTx.WithLabelValues(hostname, instanceLabel, i).Set(f.Get("tx_packets").Float())
		// These are counter type values; only ever increasing.  They're
		// represented as Gauges as they don't increase in known, linear increments.
		netPktRxDrp.WithLabelValues(hostname, instanceLabel, i).Set(f.Get("rx_packets_dropped").Float())
		netPktTxDrp.WithLabelValues(hostname, instanceLabel, i).Set(f.Get("tx_packets_dropped").Float())
	}
}

// cpuLogical iterates over the content of the njmon.cpu_logical data, writing
// a set of metrics for each logical CPU.
func cpuLogical(hostname, instanceLabel string, result gjson.Result) {
	for cpuNum, f := range result.Map() {
		// Divide these by 100 to express percentages as 0-1.
		cpuLogIdle.WithLabelValues(hostname, instanceLabel, cpuNum).Set(f.Get("idle").Float() / 100)
		cpuLogSys.WithLabelValues(hostname, instanceLabel, cpuNum).Set(f.Get("sys").Float() / 100)
		cpuLogUser.WithLabelValues(hostname, instanceLabel, cpuNum).Set(f.Get("user").Float() / 100)
		cpuLogWait.WithLabelValues(hostname, instanceLabel, cpuNum).Set(f.Get("wait").Float() / 100)
	}
}

// memPages iterates through the njmon memory_pages and writes a set of metrics for each page size.
func memPages(hostname, instanceLabel string, result gjson.Result) {
	for psize, f := range result.Map() {
		memPageFaults.WithLabelValues(hostname, instanceLabel, psize).Set(f.Get("pgexct").Float())
		memPageIns.WithLabelValues(hostname, instanceLabel, psize).Set(f.Get("pgins").Float())
		memPageOuts.WithLabelValues(hostname, instanceLabel, psize).Set(f.Get("pgouts").Float())
	}
}

// clockDiff returns the difference (in seconds) between a supplied timestamp and local UTC
func clockDiff(timestamp string) float64 {
	/*
		The clock difference will be skewed by the latency between njmon
		creating and the exporter receiving the metric. For alerting purposes,
		a difference of several seconds should be tolerated.
	*/
	format := "2006-01-02T15:04:05"
	t1, err := time.Parse(format, timestamp)
	if err != nil {
		log.Warnf("Cannot parse timestamp: %s", timestamp)
		// return 0 in the absense of anything more meaningful
		return 0
	}
	t2 := time.Now().UTC()
	// We always want a positive diff, regardless of which clock is ahead.
	var diff time.Duration
	if t1.After(t2) {
		diff = t1.Sub(t2)
	} else {
		diff = t2.Sub(t1)
	}
	return diff.Seconds()
}

// connHandle processes each incoming TCP connection.  It reads a byte string from the connection and, if it's deemed
// valid, hands it to the Prometheus metric parser.
func (h *hostInfoMap) connHandle(conn net.Conn) {
	// Close the connection when we're done with it.
	defer conn.Close()
	// Extract the remote address from the connection string
	remoteAddr := strings.Split(conn.RemoteAddr().String(), ":")[0]
	log.Debugf("Processing connection from: %s", remoteAddr)
	// Set a connection timeout duration
	conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(cfg.Thresholds.ConnectionTimout)))
	reader := bufio.NewReader(conn)
	b, err := reader.ReadBytes('\x0a')
	if err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			log.Infof("%s: Connection timeout", remoteAddr)
		} else {
			log.Warnf("%s: Error reading njmon data: %v", remoteAddr, err)
		}
		return
	}
	h.parseNJmonJSON(gjson.ParseBytes(b))
}

// parseNJmonJSON takes a gjson result object and parses it into Prometheus metrics.
func (h *hostInfoMap) parseNJmonJSON(jp gjson.Result) {
	jvalue := jp.Get("identity.hostname")
	if !jvalue.Exists() {
		log.Warn("Unable to read hostname from njmon json")
		return
	}
	hostname := jvalue.String()
	log.Debugf("Extracted hostname=%s from njmon data", hostname)
	instanceLabel := h.registerHost(hostname)

	// Compare the local clock with the timestamp provided by njmon.  It's important to be aware that these two clocks
	// will not be perfectly aligned!  There is latency introduced between NJmon capturing a timestamp and
	// njmon_exporter comparing it to the local clock.  The purpose of this metric is to provide a means to generate
	// an alert if a timestamp is off by several seconds, not several milliseconds!
	driftSecs := clockDiff(jp.Get("timestamp.UTC").String())
	log.Debugf("%s: Clock difference=%f seconds", hostname, driftSecs)
	clockDrift.WithLabelValues(hostname, instanceLabel).Set(driftSecs)

	// Uptime has only minute level granularity but we convert it to seconds for metric consistency.
	uptimeDays := jp.Get("uptime.days").Float()
	uptimeHours := jp.Get("uptime.hours").Float()
	uptimeMins := jp.Get("uptime.minutes").Float()
	uptimeSecs := (uptimeDays * 24 * 60 * 60) + (uptimeHours * 60 * 60) + (uptimeMins * 60)
	systemUptime.WithLabelValues(hostname, instanceLabel).Set(uptimeSecs)
	// server
	aixVersion.WithLabelValues(hostname, instanceLabel).Set(jp.Get("server.aix_version").Float())
	aixTechLevel.WithLabelValues(hostname, instanceLabel).Set(jp.Get("server.aix_technology_level").Float())
	aixServicePack.WithLabelValues(hostname, instanceLabel).Set(jp.Get("server.aix_service_pack").Float())
	// config
	memDesired.WithLabelValues(hostname, instanceLabel).Set(jp.Get("config.mem_desired").Float() * mb)
	memMax.WithLabelValues(hostname, instanceLabel).Set(jp.Get("config.mem_max").Float() * mb)
	memMin.WithLabelValues(hostname, instanceLabel).Set(jp.Get("config.mem_min").Float() * mb)
	memOnline.WithLabelValues(hostname, instanceLabel).Set(jp.Get("config.mem_online").Float() * mb)
	cpuPhysMax.WithLabelValues(hostname, instanceLabel).Set(jp.Get("config.pcpu_max").Float())
	cpuPhysOnline.WithLabelValues(hostname, instanceLabel).Set(jp.Get("config.pcpu_online").Float())
	cpuVirtDesired.WithLabelValues(hostname, instanceLabel).Set(jp.Get("config.vcpus_desired").Float())
	cpuVirtMax.WithLabelValues(hostname, instanceLabel).Set(jp.Get("config.vcpus_max").Float())
	cpuVirtMin.WithLabelValues(hostname, instanceLabel).Set(jp.Get("config.vcpus_min").Float())
	cpuVirtOnline.WithLabelValues(hostname, instanceLabel).Set(jp.Get("config.vcpus_online").Float())
	// memory
	memPgspFree.WithLabelValues(hostname, instanceLabel).Set(jp.Get("memory.pgsp_free").Float() * page)
	memPgspRsvd.WithLabelValues(hostname, instanceLabel).Set(jp.Get("memory.pgsp_rsvd").Float() * page)
	memPgspTotal.WithLabelValues(hostname, instanceLabel).Set(jp.Get("memory.pgsp_total").Float() * page)
	memRealFree.WithLabelValues(hostname, instanceLabel).Set(jp.Get("memory.real_free").Float() * page)
	memRealInUse.WithLabelValues(hostname, instanceLabel).Set(jp.Get("memory.real_inuse").Float() * page)
	memRealPinned.WithLabelValues(hostname, instanceLabel).Set(jp.Get("memory.real_pinned").Float() * page)
	memRealProcess.WithLabelValues(hostname, instanceLabel).Set(jp.Get("memory.real_process").Float() * page)
	memRealSystem.WithLabelValues(hostname, instanceLabel).Set(jp.Get("memory.real_system").Float() * page)
	memRealTotal.WithLabelValues(hostname, instanceLabel).Set(jp.Get("memory.real_total").Float() * page)
	memRealUser.WithLabelValues(hostname, instanceLabel).Set(jp.Get("memory.real_user").Float() * page)
	// cpu_details
	cpuNumActive.WithLabelValues(hostname, instanceLabel).Set(jp.Get("cpu_details.cpus_active").Float())
	cpuNumConf.WithLabelValues(hostname, instanceLabel).Set(jp.Get("cpu_details.cpus_configured").Float())
	cpuMHz.WithLabelValues(hostname, instanceLabel).Set(jp.Get("cpu_details.mhz").Float())
	// cpu_util
	cpuTotIdle.WithLabelValues(hostname, instanceLabel).Set(jp.Get("cpu_util.idle_pct").Float() / 100)
	cpuTotKern.WithLabelValues(hostname, instanceLabel).Set(jp.Get("cpu_util.kern_pct").Float() / 100)
	cpuTotUser.WithLabelValues(hostname, instanceLabel).Set(jp.Get("cpu_util.user_pct").Float() / 100)
	cpuTotWait.WithLabelValues(hostname, instanceLabel).Set(jp.Get("cpu_util.wait_pct").Float() / 100)
	// cpu_logical
	cpuLogical(hostname, instanceLabel, jp.Get("cpu_logical"))
	// disks
	disks(hostname, instanceLabel, jp.Get("disks"))
	// filesystems
	filesystems(hostname, instanceLabel, jp.Get("filesystems"))
	// network_adapters
	netAdapters(hostname, instanceLabel, jp.Get("network_adapters"))
	// memory_pages
	memPages(hostname, instanceLabel, jp.Get("memory_page"))
}
