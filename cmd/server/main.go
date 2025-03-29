package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// MemStorage - структура для хранения метрик
type MemStorage struct {
	gauges   map[string]float64 // Хранение метрик типа gauge
	counters map[string]int64   // Хранение метрик типа counter
}

// Storage - интерфейс для взаимодействия с хранилищем
type Storage interface {
	SetGauge(name string, value float64)  // Установить значение gauge
	AddCounter(name string, value int64)  // Добавить значение counter
	GetGauge(name string) (float64, bool) // Получить значение gauge
	GetCounter(name string) (int64, bool) // Получить значение counter
}

// SetGauge - устанавливает значение метрики типа gauge
func (s *MemStorage) SetGauge(name string, value float64) {
	if s.gauges == nil {
		s.gauges = make(map[string]float64)
	}
	s.gauges[name] = value
}

// AddCounter - добавляет значение метрики типа counter
func (s *MemStorage) AddCounter(name string, value int64) {
	if s.counters == nil {
		s.counters = make(map[string]int64)
	}
	s.counters[name] += value
}

// GetGauge - возвращает значение метрики типа gauge
func (s *MemStorage) GetGauge(name string) (float64, bool) {
	value, exists := s.gauges[name]
	return value, exists
}

// GetCounter - возвращает значение метрики типа counter
func (s *MemStorage) GetCounter(name string) (int64, bool) {
	value, exists := s.counters[name]
	return value, exists
}

// Обработчик для пути /update
func updateHandler(storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Разбиваем путь на части
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 5 {
			http.Error(w, "Invalid request path", http.StatusNotFound)
			return
		}

		// Извлекаем тип метрики, имя и значение
		metricType := parts[2]
		metricName := parts[3]
		metricValue := parts[4]

		// Обрабатываем метрику в зависимости от типа
		switch metricType {
		case "gauge":
			value, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				http.Error(w, "Invalid gauge value", http.StatusBadRequest)
				return
			}
			storage.SetGauge(metricName, value)
		case "counter":
			value, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				http.Error(w, "Invalid counter value", http.StatusBadRequest)
				return
			}
			storage.AddCounter(metricName, value)
		default:
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		// Возвращаем успешный статус
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	// Создаем хранилище
	storage := &MemStorage{}

	// Регистрируем обработчик
	http.HandleFunc("/update/", updateHandler(storage))

	// Запускаем сервер
	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
