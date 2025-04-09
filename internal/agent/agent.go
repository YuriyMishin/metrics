package agent

import (
	"YuriyMishin/metrics/internal/config"
	"log"
	"time"
)

type Agent struct {
	config  *config.AgentConfig
	sender  Sender
	metrics *Metrics
}

func NewAgent(agentConfig *config.AgentConfig) *Agent {
	return &Agent{
		config:  agentConfig,
		sender:  NewRestySender("http://" + agentConfig.Addr),
		metrics: NewMetrics(),
	}
}

func (a *Agent) Run() {
	pollTicker := time.NewTicker(a.config.PollInterval)
	reportTicker := time.NewTicker(a.config.ReportInterval)
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
