package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc/codes"
	"strconv"
	"time"
)

var requestMetricsSignUp = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "signup",
	Subsystem:  "grpc",
	Name:       "request",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
}, []string{"status"})

func ObserveSignUpRequest(d time.Duration, code codes.Code) {
	requestMetricsSignUp.WithLabelValues(strconv.Itoa(MapGRPCCodeToHTTPCode(code))).Observe(d.Seconds())
}
