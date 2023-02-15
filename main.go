package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var myCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "my_counter",
		Help: "This is my example counter.",
	})

func init() {
	prometheus.MustRegister(myCounter)
}

func handleRequest() {
	// your code here
	myCounter.Inc()
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		handleRequest()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!"))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
