# njmon_exporter

Prometheus exporter for IBM AIX systems running [njmon](http://nmon.sourceforge.net/pmwiki.php?n=Site.NjmonManualPage).

# WIP - This exporter is currently in development.
- [x] Create a TCP server for inbound njmon connections
- [x] Export initial metrics to Prometheus exporter format and serve them
- [x] Add support for a yaml configuration file
- [x] Register an exporter port in the [Prometheus Wiki](https://github.com/prometheus/prometheus/wiki/Default-port-allocations)
- [x] Support command line flags
- [x] Logging
- [ ] Config validation support
- [ ] Provide decent unit test coverage
- [ ] Declare it ready for use

## Overview
njmon_exporter provides an interface between Prometheus and the AIX njmon tool (http://nmon.sourceforge.net/pmwiki.php?n=Site.Njmon).  It listens for the connections described below:
* njmon instances running on IBM AIX systems (Default port: 8081)
* Promtheus scrapes (Default port: 9772)

## Usage
* Compile the Go code to create an njmon_exporter binary for your system.
* Copy the example njmon_exporter.yml file from the repository to your preferred configuration directory.  E.g. /etc/prometheus/njmon_exporter.yml.
* Test the exporter: `njmon_exporter --config /etc/prometheus/njmon_exporter.yml`