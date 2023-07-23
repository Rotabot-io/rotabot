package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var RequestDurationBuckets = []float64{
	0.025,
	0.05,
	0.1,
	0.25,
	0.5,
	1,
	2.5,
	5,
	10,
	30,
}

var RequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "rotabot_requests_total",
	Help: "Number of requests received by the server",
}, []string{"endpoint"})

var RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "rotabot_endpoint_duration_seconds",
	Help:    "Duration for request to complete",
	Buckets: RequestDurationBuckets,
}, []string{"endpoint", "status"})

var ResponsesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "rotabot_responses_total",
	Help: "Number of responses sent by the server",
}, []string{"endpoint", "status"})

var PanicsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "rotabot_panics_total",
	Help: "Number of panics recovered",
}, []string{"endpoint"})

var AppTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "rotabot_app_total",
	Help: "Number of apps being run with a given version",
}, []string{"app_name", "sha"})
