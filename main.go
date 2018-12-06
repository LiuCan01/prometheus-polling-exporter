package main

import (
    "encoding/json"
    "fmt"
    "github.com/prometheus-polling-exporter/metric"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "net/http"
    "os"
    "time"
)

type configuration struct {
    Time_interval int
    Port string
}

func main() {
    file, _ := os.Open("/etc/polling_export/polling_export.json")
    defer file.Close()
    decoder := json.NewDecoder(file)
    conf := configuration{}
    err := decoder.Decode(&conf)
    if err != nil {
        fmt.Println("Error", err)
    }

    workerDB := metric.NewClusterManager("db")

    //one minute collect the metrics
    ticker := time.NewTicker(time.Minute * time.Duration(conf.Time_interval))
    go func() {
        for _ = range ticker.C {
            metric.HaConfig()
        }
    }()

    // Since we are dealing with custom Collector implementations, it might
    // be a good idea to try it out with a pedantic registry.
    reg := prometheus.NewPedanticRegistry()
    reg.MustRegister(workerDB)

    http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
    http.ListenAndServe(conf.Port, nil)
}
