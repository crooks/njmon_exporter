package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {
			listener()
		}
	}()
}

func main() {
	go listener()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
