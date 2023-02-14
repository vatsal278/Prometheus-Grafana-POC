package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "endpoint"},
	)
)

func main() {
	// Register the httpRequestsTotal counter with Prometheus.
	prometheus.MustRegister(httpRequestsTotal)

	// Start the HTTP server.
	http.HandleFunc("/", handleRequest)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Increment the httpRequestsTotal counter for this endpoint.
	httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path).Inc()

	// Return a response to the client.
	fmt.Fprintf(w, "Hello, world!")
}
