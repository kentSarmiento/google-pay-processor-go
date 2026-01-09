package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/kentSarmiento/google-pay-processor-go/internal/domain"
)

// PrometheusCollector implements the domain.MetricsCollector interface
type PrometheusCollector struct {
	counters   map[string]*prometheus.CounterVec
	histograms map[string]*prometheus.HistogramVec
	gauges     map[string]*prometheus.GaugeVec
}

// NewPrometheusCollector creates a new Prometheus metrics collector
func NewPrometheusCollector() domain.MetricsCollector {
	return &PrometheusCollector{
		counters:   make(map[string]*prometheus.CounterVec),
		histograms: make(map[string]*prometheus.HistogramVec),
		gauges:     make(map[string]*prometheus.GaugeVec),
	}
}

func (p *PrometheusCollector) IncrementCounter(name string, labels map[string]string) {
	counter := p.getOrCreateCounter(name, getLabelKeys(labels))
	counter.With(prometheus.Labels(labels)).Inc()
}

func (p *PrometheusCollector) RecordDuration(name string, duration float64, labels map[string]string) {
	histogram := p.getOrCreateHistogram(name, getLabelKeys(labels))
	histogram.With(prometheus.Labels(labels)).Observe(duration)
}

func (p *PrometheusCollector) SetGauge(name string, value float64, labels map[string]string) {
	gauge := p.getOrCreateGauge(name, getLabelKeys(labels))
	gauge.With(prometheus.Labels(labels)).Set(value)
}

func (p *PrometheusCollector) getOrCreateCounter(name string, labels []string) *prometheus.CounterVec {
	if counter, exists := p.counters[name]; exists {
		return counter
	}

	counter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: "Counter metric for " + name,
		},
		labels,
	)
	p.counters[name] = counter
	return counter
}

func (p *PrometheusCollector) getOrCreateHistogram(name string, labels []string) *prometheus.HistogramVec {
	if histogram, exists := p.histograms[name]; exists {
		return histogram
	}

	histogram := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    name,
			Help:    "Histogram metric for " + name,
			Buckets: prometheus.DefBuckets,
		},
		labels,
	)
	p.histograms[name] = histogram
	return histogram
}

func (p *PrometheusCollector) getOrCreateGauge(name string, labels []string) *prometheus.GaugeVec {
	if gauge, exists := p.gauges[name]; exists {
		return gauge
	}

	gauge := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: "Gauge metric for " + name,
		},
		labels,
	)
	p.gauges[name] = gauge
	return gauge
}

func getLabelKeys(labels map[string]string) []string {
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	return keys
}
