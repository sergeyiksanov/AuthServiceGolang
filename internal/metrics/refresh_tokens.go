package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc/codes"
	"strconv"
	"time"
)

var requestMetricsRefreshTokens = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "refresh_tokens",
	Subsystem:  "grpc",
	Name:       "request",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
}, []string{"status"})

func ObserveRefreshTokensRequest(d time.Duration, code codes.Code) {
	requestMetricsRefreshTokens.WithLabelValues(strconv.Itoa(MapGRPCCodeToHTTPCode(code))).Observe(d.Seconds())
}
