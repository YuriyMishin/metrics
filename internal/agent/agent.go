package agent

import (
	"log"
	"time"
)

type Agent struct {
	pollInterval   time.Duration
	reportInterval time.Duration
	sender         Sender
	metrics        *Metrics
}

func NewAgent(pollInterval, reportInterval time.Duration, sender Sender) *Agent {
	return &Agent{
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		sender:         sender,
		metrics:        NewMetrics(),
	}
}

func (a *Agent) Run() {
	pollTicker := time.NewTicker(a.pollInterval)
	reportTicker := time.NewTicker(a.reportInterval)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			a.metrics.UpdateMetrics()
			log.Println("Metrics updated")
		case <-reportTicker.C:
			if err := a.sender.Send(a.metrics); err != nil {
				log.Printf("Failed to send metrics: %v", err)
			} else {
				log.Println("Metrics sent successfully")
			}
		}
	}
}
