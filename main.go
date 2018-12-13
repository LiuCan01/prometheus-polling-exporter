package main

import (
	"github.com/prometheus-polling-exporter/metric"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)


func main() {
	workerDB := metric.NewClusterManager()

	//one minute collect the metrics
	ticker := time.NewTicker(time.Minute * 1)
	go func() {
		for _ = range ticker.C {
			metric.Get_ceph_health_metric()
		}

	}()

	// Since we are dealing with custom Collector implementations, it might
	// be a good idea to try it out with a pedantic registry.
	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(workerDB)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.ListenAndServe(":9105", nil)
}
