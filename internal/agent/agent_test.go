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

func TestAgent_Run(t *testing.T) {
	mockSender := new(MockSender)
	mockSender.On("Send", mock.Anything).Return(nil)

	agent_config, _ := config.NewAgentConfig()
	agent_config.PollInterval = 10 * time.Millisecond
	agent_config.ReportInterval = 20 * time.Millisecond
	agent := NewAgent(agent_config)
	agent.sender = mockSender
	// Запускаем агент на короткое время
	go agent.Run()
	time.Sleep(50 * time.Millisecond)

	mockSender.AssertCalled(t, "Send", mock.Anything)
	mockSender.AssertNumberOfCalls(t, "Send", 2) // Должен успеть отправить 2 раза
}

func TestNewAgent(t *testing.T) {
	mockSender := new(MockSender)

	agent_config, _ := config.NewAgentConfig()
	agent_config.PollInterval = 1 * time.Second
	agent_config.ReportInterval = 2 * time.Second

	agent := NewAgent(agent_config)
	agent.sender = mockSender
	assert.Equal(t, agent_config.PollInterval, agent.config.PollInterval)
	assert.Equal(t, agent_config.ReportInterval, agent.config.ReportInterval)
	assert.Equal(t, mockSender, agent.sender)
	assert.NotNil(t, agent.metrics)
}
