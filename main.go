package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/olivierlemasle/opsgenie_exporter/pkg/collector"
	"github.com/olivierlemasle/opsgenie_exporter/pkg/config"
	"github.com/olivierlemasle/opsgenie_exporter/pkg/log"
)

const name = "opsgenie_exporter"

var (
	listenAddress = kingpin.Flag(
		"web.listen-address",
		"Address on which to expose metrics and web interface.",
	).Default(":3000").String()
	metricsPath = kingpin.Flag(
		"web.telemetry-path",
		"Path under which to expose metrics.",
	).Default("/metrics").String()
	disableExporterMetrics = kingpin.Flag(
		"web.disable-exporter-metrics",
		"Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).",
	).Bool()
)

func main() {
	// Setup CLI
	log.AddFlags(kingpin.CommandLine)
	config.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print(name))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	// Logs
	logger := log.Logger()

	// Config
	config.LoadConfiguration(logger)

	// Create collector
	collector, err := collector.NewCollector(logger)
	if err != nil {
		logger.Errorf("Cannot create collector: %v", err)
		os.Exit(1)
	}

	// Start HTTP server
	s := createHTTPServer(collector, logger)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		// Error when starting or closing listener
		logger.Errorf("HTTP server ListenAndServe: %v", err)
		os.Exit(1)
	}
}

func createHTTPServer(collector prometheus.Collector, logger *logrus.Logger) *http.Server {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		collector,
		version.NewCollector(name),
	)

	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	if !*disableExporterMetrics {
		// Add metrics about the exporter itself
		registry.MustRegister(
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
			collectors.NewGoCollector(),
		)
		// Instrument the HTTP handler (add promhttp_* metrics)
		http.Handle(*metricsPath, promhttp.InstrumentMetricHandler(
			registry, handler,
		))
	} else {
		// Do not not export metrics about the exporter itself
		http.Handle(*metricsPath, handler)
	}

	// Serve landing page
	landingPage := []byte(`<!DOCTYPE html>
			<html>
			<head><title>Opsgenie Exporter</title></head>
			<body>
			<h1>Opsgenie Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(landingPage)
	})

	// Create HTTP server
	server := &http.Server{Addr: *listenAddress}

	// Handle SIGINT and SIGTERM to gracefully shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		logger.Info("HTTP sever graceful shutdown...")
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Errorf("Error during HTTP server shutdown: %v", err)
		}
		logger.Info("HTTP server closed.")
	}()

	logger.Infof("Server listening on %v", *listenAddress)
	return server
}
