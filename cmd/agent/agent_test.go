package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMetrics(t *testing.T) {
	m := NewMetrics()
	if m.gauges == nil {
		t.Error("Expected gauges map to be initialized")
	}
	if m.counters == nil {
		t.Error("Expected counters map to be initialized")
	}
	if m.pollCount != 0 {
		t.Error("Expected pollCount to be 0")
	}
}

func TestUpdateRuntimeMetrics(t *testing.T) {
	m := NewMetrics()
	m.UpdateRuntimeMetrics()

	// Проверяем наличие основных метрик
	requiredGauges := []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
		"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects",
	}

	for _, metric := range requiredGauges {
		if _, exists := m.gauges[metric]; !exists {
			t.Errorf("Expected metric %s to be present", metric)
		}
	}

	// Проверяем, что значения метрик не нулевые (хотя бы некоторые)
	if m.gauges["Alloc"] <= 0 {
		t.Error("Expected Alloc to have positive value")
	}
}

func TestUpdateCustomMetrics(t *testing.T) {
	m := NewMetrics()
	initialPollCount := m.pollCount

	m.UpdateCustomMetrics()

	if m.pollCount != initialPollCount+1 {
		t.Error("Expected pollCount to increment by 1")
	}

	if _, exists := m.counters["PollCount"]; !exists {
		t.Error("Expected PollCount counter to be set")
	}

	if _, exists := m.gauges["RandomValue"]; !exists {
		t.Error("Expected RandomValue gauge to be set")
	}

	// Проверяем, что RandomValue в допустимом диапазоне
	if m.gauges["RandomValue"] < 0 || m.gauges["RandomValue"] >= 1 {
		t.Error("Expected RandomValue to be between 0 and 1")
	}
}

func TestHTTPSender_Send(t *testing.T) {
	// Создаем тестовый сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	sender := NewHTTPSender(ts.URL)
	metrics := NewMetrics()
	metrics.UpdateRuntimeMetrics()
	metrics.UpdateCustomMetrics()

	err := sender.Send(metrics)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestHTTPSender_Send_Error(t *testing.T) {
	// Сервер, который возвращает ошибку
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	sender := NewHTTPSender(ts.URL)
	metrics := NewMetrics()
	metrics.UpdateCustomMetrics()

	err := sender.Send(metrics)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestSendRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	err := sendRequest(ts.URL)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSendRequest_Error(t *testing.T) {
	// Несуществующий URL для теста ошибки соединения
	err := sendRequest("http://nonexistent-server")
	if err == nil {
		t.Error("Expected connection error, got nil")
	}
}
