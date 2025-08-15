package metric

import "github.com/prometheus/client_golang/prometheus"

type Metric struct {
	RequestCount prometheus.Gauge
	ValidRequestCount prometheus.Gauge
}

func NewMetric(reg prometheus.Registerer) *Metric {
	m := &Metric{
		RequestCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "queue",
			Name: "requests_made",
			Help: "number of requests received by the cimri queue microservice",
		}),	
		ValidRequestCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "queue",
			Name: "valid_requests_made",
			Help: "number of valid requests sent to the cimri workerpool microservice",
		}),	

	}
	reg.MustRegister(m.RequestCount)
	reg.MustRegister(m.ValidRequestCount)
	return m
}

func (m Metric) IncrementRequestCount() {
	m.RequestCount.Add(1)
}

func (m Metric) IncrementValidRequestCount() {
	m.ValidRequestCount.Add(1)
}