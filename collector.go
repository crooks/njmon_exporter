package main

import (
	"strings"

	"github.com/Masterminds/log-go"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	aixServicePack *prometheus.GaugeVec
	aixTechLevel   *prometheus.GaugeVec
	aixVersion     *prometheus.GaugeVec
	clockDrift     *prometheus.GaugeVec
	cpuLogIdle     *prometheus.GaugeVec
	cpuLogSys      *prometheus.GaugeVec
	cpuLogUser     *prometheus.GaugeVec
	cpuLogWait     *prometheus.GaugeVec
	cpuMHz         *prometheus.GaugeVec
	cpuNumActive   *prometheus.GaugeVec
	cpuNumConf     *prometheus.GaugeVec
	cpuPhysMax     *prometheus.GaugeVec
	cpuPhysOnline  *prometheus.GaugeVec
	cpuTotIdle     *prometheus.GaugeVec
	cpuTotKern     *prometheus.GaugeVec
	cpuTotUser     *prometheus.GaugeVec
	cpuTotWait     *prometheus.GaugeVec
	cpuVirtDesired *prometheus.GaugeVec
	cpuVirtMax     *prometheus.GaugeVec
	cpuVirtMin     *prometheus.GaugeVec
	cpuVirtOnline  *prometheus.GaugeVec
	filesystemFree *prometheus.GaugeVec
	filesystemSize *prometheus.GaugeVec
	hostUp         *prometheus.GaugeVec
	memDesired     *prometheus.GaugeVec
	memMax         *prometheus.GaugeVec
	memMin         *prometheus.GaugeVec
	memOnline      *prometheus.GaugeVec
	memPageFaults  *prometheus.GaugeVec
	memPageIns     *prometheus.GaugeVec
	memPageOuts    *prometheus.GaugeVec
	memPgspFree    *prometheus.GaugeVec
	memPgspRsvd    *prometheus.GaugeVec
	memPgspTotal   *prometheus.GaugeVec
	memRealFree    *prometheus.GaugeVec
	memRealInUse   *prometheus.GaugeVec
	memRealPinned  *prometheus.GaugeVec
	memRealProcess *prometheus.GaugeVec
	memRealSystem  *prometheus.GaugeVec
	memRealTotal   *prometheus.GaugeVec
	memRealUser    *prometheus.GaugeVec
	netBpsRx       *prometheus.GaugeVec
	netBpsTx       *prometheus.GaugeVec
	netPktRxDrp    *prometheus.GaugeVec
	netPktRx       *prometheus.GaugeVec
	netPktTxDrp    *prometheus.GaugeVec
	netPktTx       *prometheus.GaugeVec
	systemUptime   *prometheus.GaugeVec
)

func initCollectors() {
	defaultLabels := []string{"instance", cfg.InstanceLabel.Name}
	log.Debugf("Default labels: %v", strings.Join(defaultLabels, ","))

	// AIX section
	aixVersion = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_aix_version",
			Help: "AIX version number",
		},
		defaultLabels,
	)

	aixTechLevel = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_aix_techlevel",
			Help: "AIX technology level",
		},
		defaultLabels,
	)
	aixServicePack = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_aix_servicepack",
			Help: "AIX service pack number",
		},
		defaultLabels,
	)
	// Timestamp
	clockDrift = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_clock_drift",
			Help: "Difference between remote UTC and local UTC in seconds",
		},
		defaultLabels,
	)
	hostUp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "up",
			Help: "Returns 0 if a known host has stopped submitting metrics",
		},
		defaultLabels,
	)
	// CPU Details
	cpuNumActive = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_details_active",
			Help: "Number of active CPUs",
		},
		defaultLabels,
	)
	cpuNumConf = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_details_configured",
			Help: "Number of configured CPUs",
		},
		defaultLabels,
	)
	cpuMHz = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_details_mhz",
			Help: "CPU speed in MHz",
		},
		defaultLabels,
	)
	// CPU Utilization
	cpuTotIdle = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_total_idle",
			Help: "Percentage of CPU cycles spent idle",
		},
		defaultLabels,
	)
	cpuTotKern = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_total_kern",
			Help: "Percentage of CPU cycles spent on kernel",
		},
		defaultLabels,
	)
	cpuTotUser = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_total_user",
			Help: "Percentage of CPU cycles spent on user",
		},
		defaultLabels,
	)
	cpuTotWait = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_total_wait",
			Help: "Percentage of CPU cycles spent waiting",
		},
		defaultLabels,
	)
	// CPU Logical
	cpuLogIdle = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_logical_idle",
			Help: "Percentage of logical CPU cycles spent idle",
		},
		append(defaultLabels, "cpu"),
	)
	cpuLogSys = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_logical_sys",
			Help: "Percentage of logical CPU cycles spent on sys",
		},
		append(defaultLabels, "cpu"),
	)
	cpuLogUser = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_logical_user",
			Help: "Percentage of logical CPU cycles spent on user",
		},
		append(defaultLabels, "cpu"),
	)
	cpuLogWait = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_logical_wait",
			Help: "Percentage of logical CPU cycles spent waiting",
		},
		append(defaultLabels, "cpu"),
	)
	// CPU Physical
	cpuPhysMax = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_physical_max",
			Help: "Number of physical CPUs installed",
		},
		defaultLabels,
	)
	cpuPhysOnline = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_physical_online",
			Help: "Number of physical CPUs online",
		},
		defaultLabels,
	)
	// CPU Virtual
	cpuVirtMin = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_virtual_min",
			Help: "Minimum number of virtual CPUs in the LPAR",
		},
		defaultLabels,
	)
	cpuVirtMax = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_virtual_max",
			Help: "Maximum number of virtual CPUs in the LPAR",
		},
		defaultLabels,
	)
	cpuVirtDesired = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_virtual_desired",
			Help: "Desired number of virtual CPUs in the LPAR",
		},
		defaultLabels,
	)
	cpuVirtOnline = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_cpu_virtual_online",
			Help: "Online number of virtual CPUs in the LPAR",
		},
		defaultLabels,
	)
	// Filesystems
	filesystemSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_filesystem_size",
			Help: "Size of the filesystem in Bytes",
		},
		append(defaultLabels, "device", "mountpoint"),
	)
	filesystemFree = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_filesystem_free",
			Help: "Available filesystem space in Bytes",
		},
		append(defaultLabels, "device", "mountpoint"),
	)
	// Memory
	memDesired = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_desired",
			Help: "Desired memory allocated in Bytes",
		},
		defaultLabels,
	)
	memMax = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_max",
			Help: "Memory maximum that can be allocated in Bytes",
		},
		defaultLabels,
	)
	memMin = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_min",
			Help: "Minimum maximum that can be allocated in Bytes",
		},
		defaultLabels,
	)
	memOnline = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_online",
			Help: "Memory allocated to the system in Bytes",
		},
		defaultLabels,
	)
	memPgspFree = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_pgsp_free",
			Help: "Memory pagespace free in Bytes",
		},
		defaultLabels,
	)
	memPgspRsvd = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_pgsp_rsvd",
			Help: "Memory pagespace reserved in Bytes",
		},
		defaultLabels,
	)
	memPgspTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_pgsp_total",
			Help: "Memory pagespace total in Bytes",
		},
		defaultLabels,
	)
	memRealFree = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_real_free",
			Help: "Memory real free in Bytes",
		},
		defaultLabels,
	)
	memRealInUse = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_real_inuse",
			Help: "Memory real in-use in Bytes",
		},
		defaultLabels,
	)
	memRealPinned = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_real_pinned",
			Help: "Memory real pinned in Bytes",
		},
		defaultLabels,
	)
	memRealProcess = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_real_process",
			Help: "Memory real total in Bytes",
		},
		defaultLabels,
	)
	memRealSystem = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_real_system",
			Help: "Memory real total in Bytes",
		},
		defaultLabels,
	)
	memRealTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_real_total",
			Help: "Memory real total in Bytes",
		},
		defaultLabels,
	)
	memRealUser = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_real_user",
			Help: "Memory real total in Bytes",
		},
		defaultLabels,
	)
	// Memory Paging
	memPageFaults = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_page_faults",
			Help: "Number of page faults",
		},
		append(defaultLabels, "psize"),
	)
	memPageIns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_page_ins",
			Help: "Number of pages paged in",
		},
		append(defaultLabels, "psize"),
	)
	memPageOuts = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_mem_page_outs",
			Help: "Number of pages paged out",
		},
		append(defaultLabels, "psize"),
	)
	netPktRx = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_net_pkt_rx_total",
			Help: "Network packet receive rate",
		},
		append(defaultLabels, "interface"),
	)
	netPktRxDrp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_net_pkt_rx_drop",
			Help: "Network total RX packets dropped",
		},
		append(defaultLabels, "interface"),
	)
	netPktTx = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_net_pkt_tx_total",
			Help: "Network packet transmit rate",
		},
		append(defaultLabels, "interface"),
	)
	netPktTxDrp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_net_pkt_tx_drop",
			Help: "Network total TX packets dropped",
		},
		append(defaultLabels, "interface"),
	)
	netBpsRx = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_net_bps_rx",
			Help: "Network bytes/s receive rate",
		},
		append(defaultLabels, "interface"),
	)
	netBpsTx = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_net_bps_tx",
			Help: "Network bytes/s transmit rate",
		},
		append(defaultLabels, "interface"),
	)
	systemUptime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "njmon_system_uptime",
			Help: "System uptime in seconds",
		},
		defaultLabels,
	)

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
	prometheus.MustRegister(cpuPhysMax)
	prometheus.MustRegister(cpuPhysOnline)
	prometheus.MustRegister(cpuVirtDesired)
	prometheus.MustRegister(cpuVirtMax)
	prometheus.MustRegister(cpuVirtMin)
	prometheus.MustRegister(cpuVirtOnline)
	prometheus.MustRegister(filesystemSize)
	prometheus.MustRegister(filesystemFree)
	prometheus.MustRegister(hostUp)
	prometheus.MustRegister(memDesired)
	prometheus.MustRegister(memMax)
	prometheus.MustRegister(memMin)
	prometheus.MustRegister(memOnline)
	prometheus.MustRegister(memPageFaults)
	prometheus.MustRegister(memPageIns)
	prometheus.MustRegister(memPageOuts)
	prometheus.MustRegister(memPgspFree)
	prometheus.MustRegister(memPgspRsvd)
	prometheus.MustRegister(memPgspTotal)
	prometheus.MustRegister(memRealFree)
	prometheus.MustRegister(memRealInUse)
	prometheus.MustRegister(memRealPinned)
	prometheus.MustRegister(memRealProcess)
	prometheus.MustRegister(memRealSystem)
	prometheus.MustRegister(memRealTotal)
	prometheus.MustRegister(memRealUser)
	prometheus.MustRegister(netBpsRx)
	prometheus.MustRegister(netBpsTx)
	prometheus.MustRegister(netPktRx)
	prometheus.MustRegister(netPktRxDrp)
	prometheus.MustRegister(netPktTx)
	prometheus.MustRegister(netPktTxDrp)
	prometheus.MustRegister(systemUptime)
}
