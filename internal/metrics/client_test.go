package metrics

import (
	"fmt"
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	dto "github.com/prometheus/client_model/go"
)

const tolerance = 1e-6

type expectedRecord struct {
	name   string
	labels []string
	value  float64
}

type expectedHistogramBucket struct {
	Lower float64
	Upper float64
	Count uint64
}

func compareFloat64(a, b float64) bool {
	opt := cmp.Comparer(func(a, b float64) bool {
		return math.Abs(a-b) <= tolerance
	})
	return cmp.Diff(a, b, opt) != ""
}

func TestCounter(t *testing.T) {
	testCases := map[string]struct {
		targetMetrics map[string][]string
		countFunc     func(m *Client)
		want          []expectedRecord
	}{
		"success: no_label": {
			targetMetrics: map[string][]string{
				"test_counter": {},
			},
			countFunc: func(m *Client) {
				m.Count("test_counter", 1.1)
				m.Count("test_counter", 2.2)
				m.Count("test_counter", 3)
			},
			want: []expectedRecord{
				{
					name:   "test_counter",
					labels: []string{},
					value:  6.3,
				},
			},
		},
		"success: single_label": {
			targetMetrics: map[string][]string{
				"test2_counter": {"test_label"},
			},
			countFunc: func(m *Client) {
				m.Count("test2_counter", 1.1, "test_label")
				m.Count("test2_counter", 2.2, "test_label")
				m.Count("test2_counter", 3, "test_label")
			},
			want: []expectedRecord{
				{
					name:   "test2_counter",
					labels: []string{"test_label"},
					value:  6.3,
				},
			},
		},
		"success: single_name_multiple_labels": {
			targetMetrics: map[string][]string{
				"test3_counter": {"test_label1", "test_label2"},
			},
			countFunc: func(m *Client) {
				m.Count("test3_counter", 1.1, "1", "2")
				m.Count("test3_counter", 2.2, "10", "20")
				m.Count("test3_counter", 3, "1", "2")
				m.Count("test3_counter", 4, "10", "20")
			},
			want: []expectedRecord{
				{
					name:   "test3_counter",
					labels: []string{"1", "2"},
					value:  4.1, // 1.1 + 3
				},
				{
					name:   "test3_counter",
					labels: []string{"10", "20"},
					value:  6.2, // 2.2 + 4
				},
			},
		},
		"success: multiple_name_multiple_labels": {
			targetMetrics: map[string][]string{
				"apple_counter":      {"Kogyoku", "Jonagold"}, // {"紅玉", "ジョナゴールド"},
				"strawberry_counter": {"Amao", "Tochiotome"},  // {"あまおう", "とちおとめ"},
			},
			countFunc: func(m *Client) {
				const delicious = "delicious"
				const bad = "bad"
				m.Count("apple_counter", 1, delicious, delicious)
				m.Count("apple_counter", 2, delicious, bad)
				m.Count("apple_counter", 3, delicious, delicious)
				m.Count("apple_counter", 4, bad, bad)
				m.Count("strawberry_counter", 5, delicious, delicious)
				m.Count("strawberry_counter", 6, bad, delicious)
				m.Count("strawberry_counter", 7, bad, delicious)
			},
			want: []expectedRecord{
				{
					name:   "apple_counter",
					labels: []string{"delicious", "delicious"},
					value:  4, // 1 + 3
				},
				{
					name:   "apple_counter",
					labels: []string{"delicious", "bad"},
					value:  2, // 2
				},
				{
					name:   "apple_counter",
					labels: []string{"bad", "bad"},
					value:  4, // 4
				},
				{
					name:   "strawberry_counter",
					labels: []string{"delicious", "delicious"},
					value:  5, // 5
				},
				{
					name:   "strawberry_counter",
					labels: []string{"bad", "delicious"},
					value:  13, // 6 + 7
				},
			},
		},
	}

	for tc, tt := range testCases {
		tt := tt
		t.Run(tc, func(t *testing.T) {
			t.Parallel()

			m := NewClient()
			for name, labels := range tt.targetMetrics {
				fmt.Println(name, labels)
				m.RegisterCounter(name, "dummy", labels...)
			}
			tt.countFunc(m)

			for _, w := range tt.want {
				metric := &dto.Metric{}
				if err := m.counterVecMap[w.name].WithLabelValues(w.labels...).Write(metric); err != nil {
					t.Errorf("failed to get metric: %v", err)
				}
				if compareFloat64(w.value, metric.Counter.GetValue()) {
					t.Errorf("want %v, but got %v", w.value, metric.Counter.GetValue())
				}
			}
		})
	}
}

func TestGauge(t *testing.T) {
	testCases := map[string]struct {
		targetMetrics map[string][]string
		setFunc       func(m *Client)
		want          []expectedRecord
	}{
		"success: no_label": {
			targetMetrics: map[string][]string{
				"test_gauge": {},
			},
			setFunc: func(m *Client) {
				m.SetGauge("test_gauge", 1.1)
				m.SetGauge("test_gauge", 3)
				m.SetGauge("test_gauge", 2.2)
			},
			want: []expectedRecord{
				{
					name:   "test_gauge",
					labels: []string{},
					value:  2.2, // 最後にセットした値
				},
			},
		},
		"success: single_label": {
			targetMetrics: map[string][]string{
				"test2_gauge": {"test_label"},
			},
			setFunc: func(m *Client) {
				m.SetGauge("test2_gauge", 1.1, "test_label")
				m.SetGauge("test2_gauge", 2.2, "test_label")
				m.SetGauge("test2_gauge", 3, "test_label")
			},
			want: []expectedRecord{
				{
					name:   "test2_gauge",
					labels: []string{"test_label"},
					value:  3, // value of last set
				},
			},
		},
		"success: single_name_multiple_labels": {
			targetMetrics: map[string][]string{
				"test3_gauge": {"test_label1", "test_label2"},
			},
			setFunc: func(m *Client) {
				m.SetGauge("test3_gauge", 1.1, "1", "2")
				m.SetGauge("test3_gauge", 2.2, "10", "20")
				m.SetGauge("test3_gauge", 3, "1", "2")
				m.SetGauge("test3_gauge", 4, "10", "20")
			},
			want: []expectedRecord{
				{
					name:   "test3_gauge",
					labels: []string{"1", "2"},
					value:  3, // value of last set in the same label
				},
				{
					name:   "test3_gauge",
					labels: []string{"10", "20"},
					value:  4, // value of last set in the same label
				},
			},
		},
		"success: multiple_name_multiple_labels": {
			targetMetrics: map[string][]string{
				"apple_gauge":      {"Kogyoku", "Jonagold"}, // Apple Kinds
				"strawberry_gauge": {"Amao", "Tochiotome"},  // Strawberry Kinds
			},
			setFunc: func(m *Client) {
				const delicious = "delicious"
				const bad = "bad"
				m.SetGauge("apple_gauge", 1, delicious, delicious)
				m.SetGauge("apple_gauge", 2, delicious, bad)
				m.SetGauge("apple_gauge", 3, delicious, delicious)
				m.SetGauge("apple_gauge", 4, bad, bad)
				m.SetGauge("strawberry_gauge", 5, delicious, delicious)
				m.SetGauge("strawberry_gauge", 6, bad, delicious)
				m.SetGauge("strawberry_gauge", 7, bad, delicious)
			},
			want: []expectedRecord{
				{
					name:   "apple_gauge",
					labels: []string{"delicious", "delicious"},
					value:  3,
				},
				{
					name:   "apple_gauge",
					labels: []string{"delicious", "bad"},
					value:  2,
				},
				{
					name:   "apple_gauge",
					labels: []string{"bad", "bad"},
					value:  4,
				},
				{
					name:   "strawberry_gauge",
					labels: []string{"delicious", "delicious"},
					value:  5,
				},
				{
					name:   "strawberry_gauge",
					labels: []string{"bad", "delicious"},
					value:  7,
				},
			},
		},
	}

	for tc, tt := range testCases {
		tt := tt
		t.Run(tc, func(t *testing.T) {
			t.Parallel()

			m := NewClient()
			for name, labels := range tt.targetMetrics {
				fmt.Println(name, labels)
				m.RegisterGauge(name, "dummy", labels...)
			}
			tt.setFunc(m)

			for _, w := range tt.want {
				metric := &dto.Metric{}
				if err := m.gaugeVecMap[w.name].WithLabelValues(w.labels...).Write(metric); err != nil {
					t.Errorf("failed to get metric: %v", err)
				}
				if compareFloat64(w.value, metric.Gauge.GetValue()) {
					t.Errorf("want %v, but got %v", w.value, metric.Gauge.GetValue())
				}
			}
		})
	}
}

func TestHistogram(t *testing.T) {
	type registeredHistogram struct {
		labels  []string
		buckets []float64
	}
	type expectedHistogram struct {
		name    string
		labels  []string
		buckets []expectedHistogramBucket
	}

	testcases := map[string]struct {
		targetMetrics map[string]registeredHistogram
		observeFunc   func(m *Client)
		want          []expectedHistogram
	}{
		"success: no_label": {
			targetMetrics: map[string]registeredHistogram{
				"test_histogram": {
					labels:  []string{},
					buckets: []float64{1.0, 2.0, 3.0, 4.0},
				},
			},
			observeFunc: func(m *Client) {
				m.ObserveHistogram("test_histogram", 0.9)
				m.ObserveHistogram("test_histogram", 1.5)
				m.ObserveHistogram("test_histogram", 2.2)
				m.ObserveHistogram("test_histogram", 2.7)
				m.ObserveHistogram("test_histogram", 3)
				m.ObserveHistogram("test_histogram", 3.4)
				m.ObserveHistogram("test_histogram", 4.00001)
			},
			want: []expectedHistogram{
				{
					name:   "test_histogram",
					labels: []string{},
					buckets: []expectedHistogramBucket{
						{Lower: math.Inf(-1), Upper: 1.0, Count: 1}, // 0.9
						{Lower: 1.0, Upper: 2.0, Count: 2},          // 1.5, 2.2
						{Lower: 2.0, Upper: 3.0, Count: 5},          // 0.9, 2.2, 2.2, 2.7, 3
						{Lower: 3.0, Upper: 4.0, Count: 6},          // 0.9, 2.2, 2.2, 2.7, 3, 3.4
					},
				},
			},
		},
		"success: single_label": {
			targetMetrics: map[string]registeredHistogram{
				"test_histogram2": {
					labels:  []string{"test_label"},
					buckets: []float64{1.0, 2.0, 3.0, 4.0},
				},
			},
			observeFunc: func(m *Client) {
				m.ObserveHistogram("test_histogram2", 0.9, "test_label")
				m.ObserveHistogram("test_histogram2", 1.5, "test_label")
				m.ObserveHistogram("test_histogram2", 2.2, "test_label")
				m.ObserveHistogram("test_histogram2", 2.7, "test_label")
				m.ObserveHistogram("test_histogram2", 3, "test_label")
				m.ObserveHistogram("test_histogram2", 3.4, "test_label")
				m.ObserveHistogram("test_histogram2", 4.00001, "test_label")
			},
			want: []expectedHistogram{
				{
					name:   "test_histogram2",
					labels: []string{"test_label"},
					buckets: []expectedHistogramBucket{
						{Lower: math.Inf(-1), Upper: 1.0, Count: 1}, // 0.9
						{Lower: 1.0, Upper: 2.0, Count: 2},          // 1.5, 2.2
						{Lower: 2.0, Upper: 3.0, Count: 5},          // 0.9, 2.2, 2.2, 2.7, 3
						{Lower: 3.0, Upper: 4.0, Count: 6},          // 0.9, 2.2, 2.2, 2.7, 3, 3.4
					},
				},
			},
		},
	}

	for tc, tt := range testcases {
		tt := tt
		t.Run(tc, func(t *testing.T) {
			t.Parallel()

			m := NewClient()
			for name, hist := range tt.targetMetrics {
				m.RegisterHistogram(name, "dummy", hist.buckets, hist.labels...)
			}
			tt.observeFunc(m)

			for _, w := range tt.want {
				m, err := m.histogramVecMap[w.name].MetricVec.GetMetricWithLabelValues(w.labels...)
				if err != nil {
					t.Errorf("failed to get metric: %v", err)
				}
				hist := &dto.Metric{}
				if err := m.Write(hist); err != nil {
					t.Errorf("failed to get metric: %v", err)
				}

				lower := math.Inf(-1)
				gotBuckets := hist.GetHistogram().GetBucket()
				for i := range gotBuckets {
					upper := gotBuckets[i].GetUpperBound()
					got := expectedHistogramBucket{
						Lower: lower,
						Upper: upper,
						Count: gotBuckets[i].GetCumulativeCount(),
					}
					want := w.buckets[i]

					if diff := cmp.Diff(got, want); diff != "" {
						t.Errorf("unexpected histogram bucket: %v", diff)
					}
					lower = upper
				}
			}
		})
	}
}
