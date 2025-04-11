package server_test

import (
	"YuriyMishin/metrics/internal/server"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) SetGauge(name string, value float64) {
	m.Called(name, value)
}

func (m *MockStorage) AddCounter(name string, value int64) {
	m.Called(name, value)
}

func (m *MockStorage) GetGauge(name string) (float64, bool) {
	args := m.Called(name)
	return args.Get(0).(float64), args.Bool(1)
}

func (m *MockStorage) GetCounter(name string) (int64, bool) {
	args := m.Called(name)
	return args.Get(0).(int64), args.Bool(1)
}

func (m *MockStorage) GetAllMetrics() (map[string]float64, map[string]int64) {
	args := m.Called()
	return args.Get(0).(map[string]float64), args.Get(1).(map[string]int64)
}

func TestRootHandler(t *testing.T) {
	tests := []struct {
		name           string
		gauges         map[string]float64
		counters       map[string]int64
		expectedOutput string
	}{
		{
			name:           "empty metrics",
			gauges:         map[string]float64{},
			counters:       map[string]int64{},
			expectedOutput: "All Metrics:\n\nGauges:\n\nCounters:\n",
		},
		{
			name:           "with metrics",
			gauges:         map[string]float64{"Alloc": 123.45},
			counters:       map[string]int64{"PollCount": 1},
			expectedOutput: "All Metrics:\n\nGauges:\nAlloc: 123.45\n\nCounters:\nPollCount: 1\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockStorage)
			mockStorage.On("GetAllMetrics").Return(tt.gauges, tt.counters)

			handler := server.NewMetricHandlers(mockStorage)
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			handler.RootHandler(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
			assert.Equal(t, tt.expectedOutput, w.Body.String())
			mockStorage.AssertExpectations(t)
		})
	}
}

func TestUpdateHandler(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		metricType    string
		metricName    string
		metricValue   string
		expectedCalls func(*MockStorage)
		expectedCode  int
		expectedBody  string
	}{
		{
			name:        "valid gauge",
			url:         "/update/gauge/Alloc/123.45",
			metricType:  "gauge",
			metricName:  "Alloc",
			metricValue: "123.45",
			expectedCalls: func(m *MockStorage) {
				m.On("SetGauge", "Alloc", 123.45).Once()
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "valid counter",
			url:         "/update/counter/PollCount/1",
			metricType:  "counter",
			metricName:  "PollCount",
			metricValue: "1",
			expectedCalls: func(m *MockStorage) {
				m.On("AddCounter", "PollCount", int64(1)).Once()
			},
			expectedCode: http.StatusOK,
		},
		{
			name:          "invalid gauge value",
			url:           "/update/gauge/Alloc/invalid",
			metricType:    "gauge",
			metricName:    "Alloc",
			metricValue:   "invalid",
			expectedCalls: func(m *MockStorage) {},
			expectedCode:  http.StatusBadRequest,
			expectedBody:  "Invalid gauge value\n",
		},
		{
			name:          "invalid metric type",
			url:           "/update/invalid/name/123",
			metricType:    "invalid",
			metricName:    "name",
			metricValue:   "123",
			expectedCalls: func(m *MockStorage) {},
			expectedCode:  http.StatusBadRequest,
			expectedBody:  "Invalid metric type\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockStorage)
			tt.expectedCalls(mockStorage)

			handler := server.NewMetricHandlers(mockStorage)

			req := httptest.NewRequest("POST", tt.url, nil)
			w := httptest.NewRecorder()

			// Устанавливаем параметры маршрута вручную
			vars := map[string]string{
				"metricType":  tt.metricType,
				"metricName":  tt.metricName,
				"metricValue": tt.metricValue,
			}
			req = mux.SetURLVars(req, vars)

			handler.UpdateHandler(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, w.Body.String())
			}
			mockStorage.AssertExpectations(t)
		})
	}
}

func TestValueHandler(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		metricType    string
		metricName    string
		expectedCalls func(*MockStorage)
		expectedCode  int
		expectedBody  string
	}{
		{
			name:       "existing gauge",
			url:        "/value/gauge/Alloc",
			metricType: "gauge",
			metricName: "Alloc",
			expectedCalls: func(m *MockStorage) {
				m.On("GetGauge", "Alloc").Return(123.45, true).Once()
			},
			expectedCode: http.StatusOK,
			expectedBody: "123.45",
		},
		{
			name:       "non-existing gauge",
			url:        "/value/gauge/Nonexistent",
			metricType: "gauge",
			metricName: "Nonexistent",
			expectedCalls: func(m *MockStorage) {
				m.On("GetGauge", "Nonexistent").Return(0.0, false).Once()
			},
			expectedCode: http.StatusNotFound,
			expectedBody: "Metric not found\n",
		},
		{
			name:       "existing counter",
			url:        "/value/counter/PollCount",
			metricType: "counter",
			metricName: "PollCount",
			expectedCalls: func(m *MockStorage) {
				m.On("GetCounter", "PollCount").Return(int64(1), true).Once()
			},
			expectedCode: http.StatusOK,
			expectedBody: "1",
		},
		{
			name:          "invalid metric type",
			url:           "/value/invalid/name",
			metricType:    "invalid",
			metricName:    "name",
			expectedCalls: func(m *MockStorage) {},
			expectedCode:  http.StatusBadRequest,
			expectedBody:  "Invalid metric type\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockStorage)
			tt.expectedCalls(mockStorage)

			handler := server.NewMetricHandlers(mockStorage)

			req := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()

			// Устанавливаем параметры маршрута вручную
			vars := map[string]string{
				"metricType": tt.metricType,
				"metricName": tt.metricName,
			}
			req = mux.SetURLVars(req, vars)

			handler.ValueHandler(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
			mockStorage.AssertExpectations(t)
		})
	}
}
