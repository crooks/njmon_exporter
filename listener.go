package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/tidwall/gjson"
)

const (
	mb   = float64(1048576)
	page = float64(4096)
)

// listener listens for connections from njmon.  It forks handleConnectinon() for each connection.
func listener() {
	njmonListen := fmt.Sprintf("%s:%s", cfg.NJmon.Address, cfg.NJmon.Port)
	// Listen for incoming connections.
	l, err := net.Listen("tcp", njmonListen)
	if err != nil {
		log.Fatalf("Unable to start njmon listener: %v", err)
	}
	// Close the listener when the application closes.
	defer l.Close()
	log.Printf("Listening for njmon connections on %s\n", njmonListen)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Unable to accept njmon connection: %v", err)
			continue
		}
		// Handle connections in a new goroutine.
		go handleConnection(conn)
	}
}

func filesystems(hostname string, result gjson.Result) {
	for _, f := range result.Map() {
		// Required labels for filesystems
		device := f.Get("device").String()
		mount := f.Get("mount").String()
		// Filesystem metrics
		size := f.Get("size_mb").Float() * mb
		free := f.Get("free_mb").Float() * mb
		filesystemSize.WithLabelValues(hostname, device, mount).Set(size)
		filesystemFree.WithLabelValues(hostname, device, mount).Set(free)
	}
}

func cpuLogical(hostname string, result gjson.Result) {
	for cpuNum, f := range result.Map() {
		cpuLogIdle.WithLabelValues(hostname, cpuNum).Set(f.Get("idle").Float())
		cpuLogSys.WithLabelValues(hostname, cpuNum).Set(f.Get("sys").Float())
		cpuLogUser.WithLabelValues(hostname, cpuNum).Set(f.Get("user").Float())
		cpuLogWait.WithLabelValues(hostname, cpuNum).Set(f.Get("wait").Float())
	}
}

// clockDiff returns the difference (in seconds) between a supplied timestamp and local UTC
func clockDiff(timestamp string) float64 {
	format := "2006-01-02T15:05:05"
	t1, _ := time.Parse(format, timestamp)
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

func handleConnection(conn net.Conn) {
	// Make a buffer to hold incoming data.
	reader := bufio.NewReader(conn)
	buf, err := reader.ReadBytes('\x0a')
	if err != nil {
		log.Printf("Error reading njmon data: %v", err)
		return
	}
	// Close the connection when you're done with it.
	conn.Close()

	jp := gjson.ParseBytes(buf)

	jvalue := jp.Get("identity.hostname")
	if !jvalue.Exists() {
		log.Println("Unable to read hostname from njmon json")
		return
	}
	hostname := jvalue.String()
	clockDrift.WithLabelValues(hostname).Set(clockDiff(jp.Get("timestamp.UTC").String()))

	// server
	aixVersion.WithLabelValues(hostname).Set(jp.Get("server.aix_version").Float())
	aixTechLevel.WithLabelValues(hostname).Set(jp.Get("server.aix_technology_level").Float())
	aixServicePack.WithLabelValues(hostname).Set(jp.Get("server.aix_service_pack").Float())
	// config
	memOnline.WithLabelValues(hostname).Set(jp.Get("config.mem_online").Float() * mb)
	memMax.WithLabelValues(hostname).Set(jp.Get("config.mem_max").Float() * mb)
	// memory
	memRealFree.WithLabelValues(hostname).Set(jp.Get("memory.real_free").Float() * page)
	memRealInUse.WithLabelValues(hostname).Set(jp.Get("memory.real_inuse").Float() * page)
	memRealPinned.WithLabelValues(hostname).Set(jp.Get("memory.real_pinned").Float() * page)
	memRealProcess.WithLabelValues(hostname).Set(jp.Get("memory.real_process").Float() * page)
	memRealSystem.WithLabelValues(hostname).Set(jp.Get("memory.real_system").Float() * page)
	memRealTotal.WithLabelValues(hostname).Set(jp.Get("memory.real_total").Float() * page)
	memRealUser.WithLabelValues(hostname).Set(jp.Get("memory.real_user").Float() * page)
	// cpu_details
	cpuNumActive.WithLabelValues(hostname).Set(jp.Get("cpu_details.cpus_active").Float())
	cpuNumConf.WithLabelValues(hostname).Set(jp.Get("cpu_details.cpus_configured").Float())
	cpuMHz.WithLabelValues(hostname).Set(jp.Get("cpu_details.mhz").Float())
	// cpu_util
	cpuTotIdle.WithLabelValues(hostname).Set(jp.Get("cpu_util.idle_pct").Float())
	cpuTotKern.WithLabelValues(hostname).Set(jp.Get("cpu_util.kern_pct").Float())
	cpuTotUser.WithLabelValues(hostname).Set(jp.Get("cpu_util.user_pct").Float())
	cpuTotWait.WithLabelValues(hostname).Set(jp.Get("cpu_util.wait_pct").Float())
	// cpu_logical
	cpuLogical(hostname, jp.Get("cpu_logical"))
	// filesystems
	filesystems(hostname, jp.Get("filesystems"))
}
