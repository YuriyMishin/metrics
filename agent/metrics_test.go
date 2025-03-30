package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetrics_UpdateRuntimeMetrics(t *testing.T) {
	m := NewMetrics()
	m.UpdateRuntimeMetrics()

	assert.NotZero(t, m.gauges["Alloc"])
	assert.NotZero(t, m.gauges["HeapAlloc"])
	// Добавьте проверки для других runtime метрик
}

func TestMetrics_UpdateCustomMetrics(t *testing.T) {
	m := NewMetrics()
	initialCount := m.pollCount

	m.UpdateCustomMetrics()

	assert.Equal(t, initialCount+1, m.pollCount)
	assert.NotZero(t, m.gauges["RandomValue"])
	assert.Equal(t, int64(1), m.counters["PollCount"])
}

func TestMetrics_UpdateAll(t *testing.T) {
	m := NewMetrics()
	m.UpdateAll()

	assert.NotEmpty(t, m.gauges)
	assert.NotEmpty(t, m.counters)
}
