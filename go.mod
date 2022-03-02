module github.com/crooks/njmon_exporter

go 1.16

require (
	github.com/Masterminds/log-go v0.4.0
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/crooks/jlog v0.0.0-20220209112859-daa09de7149a
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/tidwall/gjson v1.11.0
	golang.org/x/sys v0.0.0-20211117180635-dee7805ff2e1 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/Masterminds/log-go v0.4.0 => github.com/crooks/log-go v0.4.1
