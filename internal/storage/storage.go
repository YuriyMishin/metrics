package storage

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

type Repositories interface {
	SetGauge(name string, value float64)
	AddCounter(name string, value int64)
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
	GetAllMetrics() (map[string]float64, map[string]int64)
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (s *MemStorage) SetGauge(name string, value float64) {
	s.gauges[name] = value
}

func (s *MemStorage) AddCounter(name string, value int64) {
	s.counters[name] += value
}

func (s *MemStorage) GetGauge(name string) (float64, bool) {
	value, exists := s.gauges[name]
	return value, exists
}

func (s *MemStorage) GetCounter(name string) (int64, bool) {
	value, exists := s.counters[name]
	return value, exists
}

func (s *MemStorage) GetAllMetrics() (map[string]float64, map[string]int64) {
	return s.gauges, s.counters
}
