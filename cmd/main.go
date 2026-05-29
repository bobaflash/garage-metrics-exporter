package main

import (
	"fmt"
	"garage-metrics-exporter/internal"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func getEnv(env, defaultValue string) string {
	result := os.Getenv(env)

	if result == "" {
		if len(defaultValue) == 0 {
			err := fmt.Errorf("ENV not set: %s", env)
			slog.Error(err.Error())
			os.Exit(1)
		} else {
			slog.Warn("ENV not set, use default value", "env", env, "default", defaultValue)
			return defaultValue
		}
	}

	return result
}

func main() {

	var token = fmt.Sprintf("Bearer %s", getEnv("TOKEN", ""))
	var listenSocket = getEnv("LISTEN_SOCKET", ":3905")
	var garageBaseUrl = getEnv("GARAGE_BASE_URL", "http://localhost:3903")
	updateIntervallSec, err := strconv.Atoi(getEnv("UPDATE_INTERVAL_SECONDS", "30"))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	requester := internal.NewGarageCollector(garageBaseUrl, token)

	prometheusReg := prometheus.NewRegistry()
	metrics := internal.NewMetrics(prometheusReg)

	go startWebserver(prometheusReg, listenSocket)

	for {
		cycle(*requester, metrics)
		time.Sleep(time.Duration(updateIntervallSec) * time.Second)
	}

}

func cycle(gc internal.GarageCollector, metrics *internal.Metrics) {
	ids, err := gc.GetBucketIds()
	if err != nil {
		slog.Error(err.Error())
	}

	for _, id := range ids {
		err := gc.GetBucketIinfo(id, metrics)

		if err != nil {
			slog.Error(err.Error())
		}
	}
}

func startWebserver(prometheusReg *prometheus.Registry, listenSocket string) {

	prometheusReg.MustRegister()

	slog.Info("Starting webserver...", "Socket", listenSocket)
	http.Handle("/metrics", promhttp.HandlerFor(prometheusReg, promhttp.HandlerOpts{}))
	if err := http.ListenAndServe(listenSocket, nil); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
