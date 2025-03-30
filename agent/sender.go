package agent

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type Sender interface {
	Send(metrics *Metrics) error
}

type RestySender struct {
	client    *resty.Client
	serverURL string
}

func NewRestySender(serverURL string) *RestySender {
	return &RestySender{
		client:    resty.New(),
		serverURL: serverURL,
	}
}

func (s *RestySender) Send(metrics *Metrics) error {
	// Отправляем gauge метрики
	for name, value := range metrics.gauges {
		url := fmt.Sprintf("%s/update/gauge/%s/%f", s.serverURL, name, value)
		if err := s.sendMetric(url); err != nil {
			return err
		}
	}

	// Отправляем counter метрики
	for name, value := range metrics.counters {
		url := fmt.Sprintf("%s/update/counter/%s/%d", s.serverURL, name, value)
		if err := s.sendMetric(url); err != nil {
			return err
		}
	}

	return nil
}

func (s *RestySender) sendMetric(url string) error {
	resp, err := s.client.R().
		SetHeader("Content-Type", "text/plain").
		Post(url)

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return nil
}
