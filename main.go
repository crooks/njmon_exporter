package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	go listener()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
