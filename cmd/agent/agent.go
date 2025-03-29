package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

// Agent - структура агента
type Agent struct {
	pollInterval   time.Duration // Интервал сбора метрик
	reportInterval time.Duration // Интервал отправки метрик
	sender         Sender        // Отправитель метрик
}

// NewAgent - конструктор для Agent
func NewAgent(pollInterval, reportInterval time.Duration, sender Sender) *Agent {
	return &Agent{
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		sender:         sender,
	}
}

// Run - запуск агента
func (a *Agent) Run() {
	metrics := NewMetrics()
	pollTicker := time.NewTicker(a.pollInterval)
	reportTicker := time.NewTicker(a.reportInterval)

	for {
		select {
		case <-pollTicker.C:
			metrics.UpdateRuntimeMetrics()
			metrics.UpdateCustomMetrics()
		case <-reportTicker.C:
			if err := a.sender.Send(metrics); err != nil {
				log.Printf("Failed to send metrics: %v", err)
			}
		}
	}
}

// Metrics - структура для хранения метрик
type Metrics struct {
	gauges    map[string]float64
	counters  map[string]int64
	pollCount int64
}

// NewMetrics - конструктор для Metrics
func NewMetrics() *Metrics {
	return &Metrics{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

// UpdateRuntimeMetrics - обновление метрик из пакета runtime
func (m *Metrics) UpdateRuntimeMetrics() {

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	m.gauges["Alloc"] = float64(memStats.Alloc)
	m.gauges["BuckHashSys"] = float64(memStats.BuckHashSys)
	m.gauges["Frees"] = float64(memStats.Frees)
	m.gauges["GCCPUFraction"] = memStats.GCCPUFraction
	m.gauges["GCSys"] = float64(memStats.GCSys)
	m.gauges["HeapAlloc"] = float64(memStats.HeapAlloc)
	m.gauges["HeapIdle"] = float64(memStats.HeapIdle)
	m.gauges["HeapInuse"] = float64(memStats.HeapInuse)
	m.gauges["HeapObjects"] = float64(memStats.HeapObjects)
	m.gauges["HeapReleased"] = float64(memStats.HeapReleased)
	m.gauges["HeapSys"] = float64(memStats.HeapSys)
	m.gauges["LastGC"] = float64(memStats.LastGC)
	m.gauges["Lookups"] = float64(memStats.Lookups)
	m.gauges["MCacheInuse"] = float64(memStats.MCacheInuse)
	m.gauges["MCacheSys"] = float64(memStats.MCacheSys)
	m.gauges["MSpanInuse"] = float64(memStats.MSpanInuse)
	m.gauges["MSpanSys"] = float64(memStats.MSpanSys)
	m.gauges["Mallocs"] = float64(memStats.Mallocs)
	m.gauges["NextGC"] = float64(memStats.NextGC)
	m.gauges["NumForcedGC"] = float64(memStats.NumForcedGC)
	m.gauges["NumGC"] = float64(memStats.NumGC)
	m.gauges["OtherSys"] = float64(memStats.OtherSys)
	m.gauges["PauseTotalNs"] = float64(memStats.PauseTotalNs)
	m.gauges["StackInuse"] = float64(memStats.StackInuse)
	m.gauges["StackSys"] = float64(memStats.StackSys)
	m.gauges["Sys"] = float64(memStats.Sys)
	m.gauges["TotalAlloc"] = float64(memStats.TotalAlloc)
}

// UpdateCustomMetrics - обновление кастомных метрик
func (m *Metrics) UpdateCustomMetrics() {

	m.pollCount++
	m.counters["PollCount"] = m.pollCount
	m.gauges["RandomValue"] = rand.Float64()
}

// Sender - интерфейс для отправки метрик
type Sender interface {
	Send(metrics *Metrics) error
}

// HTTPSender - реализация Sender для отправки метрик по HTTP
type HTTPSender struct {
	serverURL string
}

// NewHTTPSender - конструктор для HTTPSender
func NewHTTPSender(serverURL string) *HTTPSender {
	return &HTTPSender{serverURL: serverURL}
}

// Send - отправка метрик на сервер
func (s *HTTPSender) Send(metrics *Metrics) error {
	for name, value := range metrics.gauges {
		url := fmt.Sprintf("%s/update/gauge/%s/%f", s.serverURL, name, value)
		if err := sendRequest(url); err != nil {
			return err
		}
	}

	for name, value := range metrics.counters {
		url := fmt.Sprintf("%s/update/counter/%s/%d", s.serverURL, name, value)
		if err := sendRequest(url); err != nil {
			return err
		}
	}

	return nil
}

// sendRequest - отправка HTTP-запроса
func sendRequest(url string) error {
	resp, err := http.Post(url, "text/plain", bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func main() {
	sender := NewHTTPSender("http://localhost:8080")
	agent := NewAgent(2000, 10000, sender)
	agent.Run()
}
