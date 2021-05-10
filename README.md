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
- [x] Add a method for appending labels to selected source hosts (E.g. Dev vs. Prod)
- [ ] Declare it ready for use

## Overview
njmon_exporter provides an interface between Prometheus and the AIX njmon tool (http://nmon.sourceforge.net/pmwiki.php?n=Site.Njmon).  It listens for the connections described below:
* njmon instances running on IBM AIX systems (Default port: 8086)
* Promtheus scrapes (Default port: 9772)

## Installation
### Install and run the exporter
* Compile the Go code to create an njmon_exporter binary for your system.
* Choose a Linux VM to host the exporter
* Copy over the compiled binary to somewhere sane. E.g. /usr/local/bin/njmon_exporter
* Copy the example njmon_exporter.yml file from the repository to your preferred configuration directory.  E.g. /etc/prometheus/njmon_exporter.yml.
* Test the exporter: `/usr/local/bin/njmon_exporter --config /etc/prometheus/njmon_exporter.yml`
### Configure your prometheus.yml
* Create a new section for njmon with content something like this:
```
- job_name: 'njmon'
    static_configs:
      - targets:
          - <exporter_host>:9772
    scrape_interval: 60s
    honor_labels: true
```
* Reload Prometheus to pick up the new config.
### Configure njmon
* njmon can be configured to run from a crontab entry.  I find this ideal but YMMV.
* Select a user account to run njmon from
* Edit the crontab (`crontab -e`)
* Insert a line like this: `4 * * * * /usr/local/bin/njmon -k -s 60 -i <exporter_host> -p 8086`.
* Run njmon to avoid waiting for crontab to start it `/usr/local/bin/njmon -k -s 60 -i <exporter_host> -p 8086`
* Repeat this section for each AIX host you want to monitor
