module github.com/crooks/njmon_exporter

go 1.16

require (
	github.com/Masterminds/log-go v1.0.0
	github.com/crooks/jlog v0.0.0-20220702135307-b00406788daa
	github.com/prometheus/client_golang v1.12.2
	github.com/prometheus/common v0.35.0 // indirect
	github.com/tidwall/gjson v1.14.1
	golang.org/x/sys v0.0.0-20220702020025-31831981b65f // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/Masterminds/log-go v1.0.0 => github.com/crooks/log-go v0.4.2
