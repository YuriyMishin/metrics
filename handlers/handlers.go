package handlers

import (
	"YuriyMishin/metrics/storage"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type MetricHandlers struct {
	storage storage.Storage
}

func NewMetricHandlers(s storage.Storage) *MetricHandlers {
	return &MetricHandlers{storage: s}
}

func (h *MetricHandlers) RootHandler(w http.ResponseWriter, r *http.Request) {
	gauges, counters := h.storage.GetAllMetrics()

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "All Metrics:")
	fmt.Fprintln(w, "\nGauges:")
	for name, value := range gauges {
		fmt.Fprintf(w, "%s: %f\n", name, value)
	}

	fmt.Fprintln(w, "\nCounters:")
	for name, value := range counters {
		fmt.Fprintf(w, "%s: %d\n", name, value)
	}
}

func (h *MetricHandlers) UpdateHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r) // Получаем параметры из URL
	metricType := vars["metricType"]
	metricName := vars["metricName"]
	metricValue := vars["metricValue"]

	switch metricType {
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Invalid gauge value", http.StatusBadRequest)
			return
		}
		h.storage.SetGauge(metricName, value)
	case "counter":
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Invalid counter value", http.StatusBadRequest)
			return
		}
		h.storage.AddCounter(metricName, value)
	default:
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MetricHandlers) ValueHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Получаем параметры из URL
	metricType := vars["metricType"]
	metricName := vars["metricName"]

	switch metricType {
	case "gauge":
		if value, exists := h.storage.GetGauge(metricName); exists {
			fmt.Fprintf(w, "Gauge %s: %f", metricName, value)
			return
		}
	case "counter":
		if value, exists := h.storage.GetCounter(metricName); exists {
			fmt.Fprintf(w, "Counter %s: %d", metricName, value)
			return
		}
	default:
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		return
	}

	http.Error(w, "Metric not found", http.StatusNotFound)
}
