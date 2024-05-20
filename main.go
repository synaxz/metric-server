package metricmemory

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	gauges     = make(map[string]prometheus.Gauge)
	counters   = make(map[string]prometheus.Counter)
	histograms = make(map[string]prometheus.Histogram)
	summaries  = make(map[string]prometheus.Summary)
)

var (
	ErrInvalidAction     = errors.New("invalid action")
	ErrInvalidMetricType = errors.New("invalid metric type")
)

type Metric struct {
	Key    string            `json:"key"`
	Value  float64           `json:"value"`
	Labels map[string]string `json:"labels"`
	Help   string            `json:"help"`
	Action string            `json:"action"`
}

func main() {

	port := os.Args[1]
	if port == "" {
		port = "8080"
	}

	registry := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = registry

	http.HandleFunc("/store/gauge", gaugeHandler)
	http.HandleFunc("/store/counter", counterHandler)
	http.HandleFunc("/store/histogram", histogramHandler)
	http.HandleFunc("/store/summary", summaryHandler)
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func gaugeHandler(w http.ResponseWriter, r *http.Request) {
	storeMetric(w, r, "gauge")
}

func counterHandler(w http.ResponseWriter, r *http.Request) {
	storeMetric(w, r, "counter")
}

func histogramHandler(w http.ResponseWriter, r *http.Request) {
	storeMetric(w, r, "histogram")
}

func summaryHandler(w http.ResponseWriter, r *http.Request) {
	storeMetric(w, r, "summary")
}

func storeMetric(w http.ResponseWriter, r *http.Request, metricType string) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var metric Metric
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&metric); err != nil {
		http.Error(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	switch metricType {
	case "gauge":
		err := storeGaugeMetric(metric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case "counter":
		err := storeCounterMetric(metric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case "histogram":
		err := storeHistogramMetric(metric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case "summary":
		err := storeSummaryMetric(metric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, ErrInvalidMetricType.Error(), http.StatusBadRequest)

	}

	w.Write([]byte("metric stored\n"))
}

func storeGaugeMetric(m Metric) error {
	gauge, exists := gauges[m.Key]
	if !exists {
		gauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        m.Key,
			Help:        m.Help,
			ConstLabels: m.Labels,
		})
		prometheus.MustRegister(gauge)
		gauges[m.Key] = gauge
	}

	switch m.Action {
	case "set":
		gauge.Set(m.Value)
	case "inc":
		gauge.Inc()
	case "dec":
		gauge.Dec()
	case "add":
		gauge.Add(m.Value)
	case "sub":
		gauge.Sub(m.Value)
	default:
		return ErrInvalidAction
	}

	return nil

}
func storeCounterMetric(m Metric) error {
	counter, exists := counters[m.Key]
	if !exists {
		counter = prometheus.NewCounter(prometheus.CounterOpts{
			Name:        m.Key,
			Help:        m.Help,
			ConstLabels: m.Labels,
		})
		prometheus.MustRegister(counter)
		counters[m.Key] = counter
	}

	switch m.Action {
	case "inc":
		counter.Inc()
	case "add":
		counter.Add(m.Value)
	default:
		return ErrInvalidAction
	}

	return nil

}
func storeHistogramMetric(m Metric) error {
	histogram, exists := histograms[m.Key]
	if !exists {
		histogram = prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:        m.Key,
			Help:        m.Help,
			ConstLabels: m.Labels,
			Buckets:     prometheus.DefBuckets,
		})
		prometheus.MustRegister(histogram)
		histograms[m.Key] = histogram
	}

	switch m.Action {
	case "observe":
		histogram.Observe(m.Value)
	default:
		return ErrInvalidAction
	}

	return nil
}
func storeSummaryMetric(m Metric) error {
	summary, exists := summaries[m.Key]
	if !exists {
		summary = prometheus.NewSummary(prometheus.SummaryOpts{
			Name:        m.Key,
			Help:        m.Help,
			ConstLabels: m.Labels,
			Objectives:  map[float64]float64{},
		})
		prometheus.MustRegister(summary)
		summaries[m.Key] = summary
	}

	switch m.Action {
	case "observe":
		summary.Observe(m.Value)
	default:
		return ErrInvalidAction
	}

	return nil
}
