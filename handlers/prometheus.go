package handlers

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type Prometheus interface {
	GetClient(metricName string) *push.Pusher
}

type promClient struct {
	URL string
}

func (p promClient) GetClient(metricName string) *push.Pusher {
	return push.New(p.URL, metricName)
}

func NewPrometheus(url string) Prometheus {
	return &promClient{URL: url}
}

func PushToPrometheus(gatewayURL string, metricsName string) (func(err error), error) {
	successGauge := prometheus.NewGauge(prometheus.GaugeOpts{Name: "success"})
	failureGauge := prometheus.NewGauge(prometheus.GaugeOpts{Name: "failure"})
	client := NewPrometheus(gatewayURL).GetClient(metricsName).Collector(successGauge).Collector(failureGauge)

	return func(err error) {
		if err != nil {
			failureGauge.Inc()
		}
		successGauge.Inc()
		go func() {
			_ = client.Push()
		}()
	}, nil
}
