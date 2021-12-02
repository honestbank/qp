package handlers

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type Prometheus interface {
	getClient(metricName string) *push.Pusher
}

type promClient struct {
	URL string
}

func (p promClient) getClient(metricName string) *push.Pusher {
	return push.New(p.URL, metricName)
}

func newPrometheus(url string) Prometheus {
	return &promClient{URL: url}
}

func PushToPrometheus(gatewayURL string, metricsName string) (func(err error), error) {
	registry := prometheus.NewRegistry()
	successGauge := prometheus.NewCounter(prometheus.CounterOpts{Name: "success"})
	failureGauge := prometheus.NewCounter(prometheus.CounterOpts{Name: "failure"})
	client := newPrometheus(gatewayURL).getClient(metricsName).Grouping("service", "qp").Collector(successGauge).Collector(failureGauge)
	if err := registry.Register(successGauge); err != nil {
		return nil, err
	}
	if err := registry.Register(failureGauge); err != nil {
		return nil, err
	}

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
