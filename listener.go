package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/tidwall/gjson"
)

const (
	connHost = "localhost"
	connPort = "3333"
	connType = "tcp"
	mb       = float64(1024 ^ 2)
)

func listener() {
	// Listen for incoming connections.
	l, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + connHost + ":" + connPort)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleConnection(conn)
	}
}

func jFloat(result gjson.Result, f string) (float64, error) {
	value := result.Get(f)
	if !value.Exists() {
		return 0, fmt.Errorf("%s: Field not found", f)
	}
	return value.Float(), nil
}

func filesystems(hostname string, result gjson.Result) {
	for _, f := range result.Map() {
		device := f.Get("device").String()
		mount := f.Get("mount").String()

		size, err := jFloat(f, "size_mb")
		if err != nil {
			fmt.Println(err)
		}
		free, err := jFloat(f, "free_mb")
		if err != nil {
			fmt.Println(err)
		}
		filesystemSize.WithLabelValues(hostname, device, mount).Set(size * mb)
		filesystemFree.WithLabelValues(hostname, device, mount).Set(free * mb)
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
		panic(err)
	}
	// Close the connection when you're done with it.
	conn.Close()

	jp := gjson.ParseBytes(buf)

	jvalue := jp.Get("identity.hostname")
	if !jvalue.Exists() {
		fmt.Println("Unable to register hostname")
		return
	}
	hostname := jvalue.String()
	clockDrift.WithLabelValues(hostname).Set(clockDiff(jp.Get("timestamp.UTC").String()))

	// server
	aixVersion.WithLabelValues(hostname).Set(jp.Get("server.aix_version").Float())
	aixTechLevel.WithLabelValues(hostname).Set(jp.Get("server.aix_technology_level").Float())
	aixServicePack.WithLabelValues(hostname).Set(jp.Get("server.aix_service_pack").Float())
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
