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
	// CPU Details
	cpuNumActive = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_details_active",
			Help: "Number of active CPUs",
		},
		[]string{"lpar"},
	)
	cpuNumConf = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_details_configured",
			Help: "Number of configured CPUs",
		},
		[]string{"lpar"},
	)
	cpuMHz = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_details_mhz",
			Help: "CPU speed in MHz",
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
	// Memory
	memOnline = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_online",
			Help: "Memory allocated to the system in Bytes",
		},
		[]string{"lpar"},
	)
	memMax = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_max",
			Help: "Memory maximum that can be allocated in Bytes",
		},
		[]string{"lpar"},
	)
	memRealFree = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_real_free",
			Help: "Memory real free in Bytes",
		},
		[]string{"lpar"},
	)
	memRealInUse = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_real_inuse",
			Help: "Memory real in-use in Bytes",
		},
		[]string{"lpar"},
	)
	memRealPinned = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_real_pinned",
			Help: "Memory real pinned in Bytes",
		},
		[]string{"lpar"},
	)
	memRealProcess = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_real_process",
			Help: "Memory real total in Bytes",
		},
		[]string{"lpar"},
	)
	memRealSystem = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_real_system",
			Help: "Memory real total in Bytes",
		},
		[]string{"lpar"},
	)
	memRealTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_real_total",
			Help: "Memory real total in Bytes",
		},
		[]string{"lpar"},
	)
	memRealUser = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_real_user",
			Help: "Memory real total in Bytes",
		},
		[]string{"lpar"},
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
	prometheus.MustRegister(cpuNumActive)
	prometheus.MustRegister(cpuNumConf)
	prometheus.MustRegister(cpuMHz)
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
	prometheus.MustRegister(memMax)
	prometheus.MustRegister(memOnline)
	prometheus.MustRegister(memRealFree)
	prometheus.MustRegister(memRealInUse)
	prometheus.MustRegister(memRealPinned)
	prometheus.MustRegister(memRealProcess)
	prometheus.MustRegister(memRealSystem)
	prometheus.MustRegister(memRealTotal)
	prometheus.MustRegister(memRealUser)
}
