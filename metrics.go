package main

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"strconv"
)

var (
	packetsReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "shoveler_packets_received",
		Help: "The total number of packets received",
	})

	validationsFailed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "shoveler_validations_failed",
		Help: "The total number of packets that failed validation",
	})

	rabbitmqReconnects = promauto.NewCounter(prometheus.CounterOpts{
		Name: "shoveler_rabbitmq_reconnects",
		Help: "The total number of reconnections to rabbitmq bus",
	})
)

func StartMetrics() {
	// Start the metrics
	if !viper.GetBool("metrics.enable") {
		return
	}

	// Get the configuration for the metrics listening port
	metricsPort := viper.GetInt("metrics.port")

	// Listen to the metrics requests in a separate thread
	go func() {
		listenAddress := ":" + strconv.Itoa(metricsPort)
		log.Debugln("Starting metrics at " + listenAddress + "/metrics")
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(listenAddress, nil)
	}()

}

