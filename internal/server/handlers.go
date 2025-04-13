package server

import (
	"YuriyMishin/metrics/internal/logger"
	"YuriyMishin/metrics/internal/storage"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type MetricHandlers struct {
	storage storage.Repositories
}

func NewMetricHandlers(s storage.Repositories) *MetricHandlers {
	return &MetricHandlers{storage: s}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (w *loggingResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := &loggingResponseWriter{ResponseWriter: w}

		next.ServeHTTP(lw, r)

		duration := time.Since(start)
		log := logger.Get().Sugar()

		log.Infow("request completed",
			"uri", r.RequestURI,
			"method", r.Method,
			"status", lw.status,
			"size", lw.size,
			"duration", duration,
		)
	})
}

func (h *MetricHandlers) RootHandler(w http.ResponseWriter, r *http.Request) {
	gauges, counters := h.storage.GetAllMetrics()

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "All Metrics:")
	fmt.Fprintln(w, "\nGauges:")
	for name, value := range gauges {
		fmt.Fprintf(w, "%s: %g\n", name, value)
	}

	fmt.Fprintln(w, "\nCounters:")
	for name, value := range counters {
		fmt.Fprintf(w, "%s: %d\n", name, value)
	}
}

const (
	MetricTypeCounter string = "counter"
	MetricTypeGauge   string = "gauge"
)

func (h *MetricHandlers) UpdateHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r) // Получаем параметры из URL
	metricType := vars["metricType"]
	metricName := vars["metricName"]
	metricValue := vars["metricValue"]

	switch metricType {
	case MetricTypeGauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Invalid gauge value", http.StatusBadRequest)
			return
		}
		h.storage.SetGauge(metricName, value)
	case MetricTypeCounter:
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
	case MetricTypeGauge:
		if value, exists := h.storage.GetGauge(metricName); exists {
			fmt.Fprintf(w, "%g", value)
			return
		}
	case MetricTypeCounter:
		if value, exists := h.storage.GetCounter(metricName); exists {
			fmt.Fprintf(w, "%d", value)
			return
		}
	default:
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		return
	}

	http.Error(w, "Metric not found", http.StatusNotFound)
}
