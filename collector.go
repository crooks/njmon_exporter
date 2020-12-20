package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// AIX section
	aixVersion = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_aix_version",
			Help: "AIX version number",
		},
		[]string{"lpar"},
	)
	aixTechLevel = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_aix_techlevel",
			Help: "AIX technology level",
		},
		[]string{"lpar"},
	)
	aixServicePack = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_aix_servicepack",
			Help: "AIX service pack number",
		},
		[]string{"lpar"},
	)
	// CPU Utilization
	cpuTotIdle = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_total_idle",
			Help: "Percentage of CPU cycles spent idle",
		},
		[]string{"lpar"},
	)
	cpuTotKern = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_total_kern",
			Help: "Percentage of CPU cycles spent on kernel",
		},
		[]string{"lpar"},
	)
	cpuTotUser = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_total_user",
			Help: "Percentage of CPU cycles spent on user",
		},
		[]string{"lpar"},
	)
	cpuTotWait = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_total_wait",
			Help: "Percentage of CPU cycles spent waiting",
		},
		[]string{"lpar"},
	)
	// CPU Logical
	cpuLogIdle = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_logical_idle",
			Help: "Percentage of logical CPU cycles spent idle",
		},
		[]string{"lpar", "cpu"},
	)
	cpuLogSys = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_logical_sys",
			Help: "Percentage of logical CPU cycles spent on sys",
		},
		[]string{"lpar", "cpu"},
	)
	cpuLogUser = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_logical_user",
			Help: "Percentage of logical CPU cycles spent on user",
		},
		[]string{"lpar", "cpu"},
	)
	cpuLogWait = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_logical_wait",
			Help: "Percentage of logical CPU cycles spent waiting",
		},
		[]string{"lpar", "cpu"},
	)
	// Filesystems
	filesystemSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_filesystem_size",
			Help: "Size of the filesystem in Bytes",
		},
		[]string{"lpar", "device", "mount"},
	)
	filesystemFree = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_filesystem_free",
			Help: "Available filesystem space in Bytes",
		},
		[]string{"lpar", "device", "mount"},
	)
	// Timestamp
	clockDrift = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_time_drift",
			Help: "Difference between remote UTC and local UTC in seconds",
		},
		[]string{"lpar"},
	)
)

func init() {
	prometheus.MustRegister(aixVersion)
	prometheus.MustRegister(aixTechLevel)
	prometheus.MustRegister(aixServicePack)
	prometheus.MustRegister(clockDrift)
	prometheus.MustRegister(cpuTotIdle)
	prometheus.MustRegister(cpuTotKern)
	prometheus.MustRegister(cpuTotUser)
	prometheus.MustRegister(cpuTotWait)
	prometheus.MustRegister(cpuLogIdle)
	prometheus.MustRegister(cpuLogSys)
	prometheus.MustRegister(cpuLogUser)
	prometheus.MustRegister(cpuLogWait)
	prometheus.MustRegister(filesystemSize)
	prometheus.MustRegister(filesystemFree)
}
