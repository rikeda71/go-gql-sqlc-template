package metrics

import (
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
)

// Client is a struct that manages metrics
// key: metric name, value: metric
type Client struct {
	counterVecMap   map[string]*prometheus.CounterVec
	gaugeVecMap     map[string]*prometheus.GaugeVec
	histogramVecMap map[string]*prometheus.HistogramVec
}

// NewClient is a constructor for Client
func NewClient() *Client {
	return &Client{
		counterVecMap:   make(map[string]*prometheus.CounterVec),
		gaugeVecMap:     make(map[string]*prometheus.GaugeVec),
		histogramVecMap: make(map[string]*prometheus.HistogramVec),
	}
}

// RegisterCounter is registers counter metrics
func (m *Client) RegisterCounter(name string, help string, labels ...string) {
	c := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
		}, labels,
	)
	prometheus.MustRegister(c)
	m.counterVecMap[name] = c
}

// RegisterGauge is registers gauge metrics
func (m *Client) RegisterGauge(name string, help string, labels ...string) {
	g := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		}, labels,
	)
	prometheus.MustRegister(g)
	m.gaugeVecMap[name] = g
}

// RegisterHistogram is registers histogram metrics
func (m *Client) RegisterHistogram(name string, help string, buckets []float64, labels ...string) {
	h := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    name,
			Help:    help,
			Buckets: buckets,
		}, labels,
	)
	prometheus.MustRegister(h)
	m.histogramVecMap[name] = h
}

// Count is increments counter metrics
func (m *Client) Count(name string, value float64, labels ...string) {
	cv, ok := m.counterVecMap[name]
	if !ok {
		slog.Warn("counter not found", "name", name)
		return
	}

	counter, err := cv.GetMetricWithLabelValues(labels...)
	if err != nil {
		slog.Warn("counter not found", "name", name, "labels", labels)
		return
	}
	counter.Add(value)
}

// SetGauge is sets gauge metrics
func (m *Client) SetGauge(name string, value float64, labels ...string) {
	gv, ok := m.gaugeVecMap[name]
	if !ok {
		slog.Warn("gauge not found", "name", name)
		return
	}

	gauge, err := gv.GetMetricWithLabelValues(labels...)
	if err != nil {
		slog.Warn("gauge not found", "name", name, "labels", labels)
		return
	}
	gauge.Set(value)
}

// ObserveHistogram is observes histogram metrics
func (m *Client) ObserveHistogram(name string, value float64, labels ...string) {
	hv, ok := m.histogramVecMap[name]
	if !ok {
		slog.Warn("histogram not found", "name", name)
		return
	}

	histogram, err := hv.GetMetricWithLabelValues(labels...)
	if err != nil {
		slog.Warn("histogram not found", "name", name, "labels", labels)
		return
	}
	histogram.Observe(value)
}
