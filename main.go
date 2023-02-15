package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	histogramVec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of HTTP request durations.",
			Buckets: []float64{0.1, 0.2, 0.5, 1, 2, 5},
		},
		[]string{"method", "path"},
	)

	counterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "Counter of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	gauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "queue_size",
			Help: "Number of items in the queue.",
		},
	)
)

func init() {
	prometheus.MustRegister(histogramVec)
	prometheus.MustRegister(gauge)
	prometheus.MustRegister(counterVec)
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		status := rand.Intn(1000)

		histogramVec.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())
		counterVec.WithLabelValues(r.Method, r.URL.Path, string(status)).Inc()
		gauge.Inc()

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!"))
	})
	http.HandleFunc("/bad", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusBadRequest) })
	http.HandleFunc("/internal", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusInternalServerError) })

	log.Fatal(http.ListenAndServe(":8080", nil))
}
