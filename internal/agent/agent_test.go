package agent

import (
	"YuriyMishin/metrics/internal/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSender struct {
	mock.Mock
}

func (m *MockSender) Send(metrics *Metrics) error {
	args := m.Called(metrics)
	return args.Error(0)
}

func TestNewAgent(t *testing.T) {
	mockSender := new(MockSender)

	agentConfig, _ := config.NewAgentConfig()
	agentConfig.PollInterval = 1 * time.Second
	agentConfig.ReportInterval = 2 * time.Second

	agent := NewAgent(agentConfig)
	agent.sender = mockSender
	assert.Equal(t, agentConfig.PollInterval, agent.config.PollInterval)
	assert.Equal(t, agentConfig.ReportInterval, agent.config.ReportInterval)
	assert.Equal(t, mockSender, agent.sender)
	assert.NotNil(t, agent.metrics)
}
