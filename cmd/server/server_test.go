package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMemStorage(t *testing.T) {
	storage := &MemStorage{}

	t.Run("Gauge operations", func(t *testing.T) {
		// Test SetGauge and GetGauge
		storage.SetGauge("temperature", 23.5)
		value, exists := storage.GetGauge("temperature")
		if !exists {
			t.Error("Expected temperature to exist")
		}
		if value != 23.5 {
			t.Errorf("Expected temperature 23.5, got %v", value)
		}

		// Test non-existent gauge
		_, exists = storage.GetGauge("nonexistent")
		if exists {
			t.Error("Expected nonexistent gauge to not exist")
		}
	})

	t.Run("Counter operations", func(t *testing.T) {
		// Test AddCounter and GetCounter
		storage.AddCounter("requests", 10)
		storage.AddCounter("requests", 5)
		value, exists := storage.GetCounter("requests")
		if !exists {
			t.Error("Expected requests to exist")
		}
		if value != 15 {
			t.Errorf("Expected requests counter 15, got %v", value)
		}

		// Test non-existent counter
		_, exists = storage.GetCounter("nonexistent")
		if exists {
			t.Error("Expected nonexistent counter to not exist")
		}
	})

	t.Run("Concurrent access", func(t *testing.T) {
		// This would be more thorough with actual concurrency tests
		// but this at least verifies the maps are initialized properly
		storage.SetGauge("new_metric", 1.0)
		storage.AddCounter("new_counter", 1)
	})
}
func TestUpdateHandler(t *testing.T) {
	storage := &MemStorage{}
	handler := updateHandler(storage)

	tests := []struct {
		name          string
		path          string
		expectedCode  int
		expectedGauge float64
		expectedCount int64
	}{
		{
			name:          "Valid gauge update",
			path:          "/update/gauge/temp/42.5",
			expectedCode:  http.StatusOK,
			expectedGauge: 42.5,
		},
		{
			name:          "Valid counter update",
			path:          "/update/counter/requests/10",
			expectedCode:  http.StatusOK,
			expectedCount: 10,
		},
		{
			name:         "Invalid gauge value",
			path:         "/update/gauge/temp/invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid counter value",
			path:         "/update/counter/requests/invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid metric type",
			path:         "/update/invalid/name/10",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid path format",
			path:         "/update/gauge/name",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedCode)
			}

			// Check storage if request was successful
			if tt.expectedCode == http.StatusOK {
				if strings.Contains(tt.path, "gauge") {
					metricName := strings.Split(tt.path, "/")[3]
					value, exists := storage.GetGauge(metricName)
					if !exists {
						t.Errorf("Expected gauge %s to be set", metricName)
					}
					if value != tt.expectedGauge {
						t.Errorf("Expected gauge value %v, got %v", tt.expectedGauge, value)
					}
				} else if strings.Contains(tt.path, "counter") {
					metricName := strings.Split(tt.path, "/")[3]
					value, exists := storage.GetCounter(metricName)
					if !exists {
						t.Errorf("Expected counter %s to be set", metricName)
					}
					if value != tt.expectedCount {
						t.Errorf("Expected counter value %v, got %v", tt.expectedCount, value)
					}
				}
			}
		})
	}
}
