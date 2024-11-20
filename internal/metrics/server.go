package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Listen(host string) error {
	log.Printf("Run metrics on: %v", host)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(host, mux)
}
