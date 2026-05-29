package internal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	Bytes *prometheus.GaugeVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		Bytes: promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
			Name: "garage_bytes_total",
			Help: "The total number of processed events",
		}, []string{"bucket"}),
	}
	return m
}
