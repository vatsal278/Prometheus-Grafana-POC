package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var (
	histogramVec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "my_request_status_code",
			Help:    "Histogram of HTTP request status codes by method and path.",
			Buckets: []float64{200, 300, 400, 500},
		},
		[]string{"method", "path", "status_code"},
	)

	counterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "my_counter",
			Help: "Counter of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	gauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "queue",
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
		histogramVec.WithLabelValues(r.Method, r.URL.Path, "200").Observe(1)
		counterVec.WithLabelValues(r.Method, r.URL.Path, "200").Inc()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!"))
	})
	http.HandleFunc("/bad", func(w http.ResponseWriter, r *http.Request) {
		histogramVec.WithLabelValues(r.Method, r.URL.Path, "400").Observe(1)
		counterVec.WithLabelValues(r.Method, r.URL.Path, "400").Inc()
		gauge.Inc()
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad, world!"))
	})
	http.HandleFunc("/internal", func(w http.ResponseWriter, r *http.Request) {
		histogramVec.WithLabelValues(r.Method, r.URL.Path, "500").Observe(1)
		counterVec.WithLabelValues(r.Method, r.URL.Path, "500").Inc()
		gauge.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Error"))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
