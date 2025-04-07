package agent

import (
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

	agent := NewAgent(
		10*time.Millisecond,
		20*time.Millisecond,
		mockSender,
	)

	// Запускаем агент на короткое время
	go agent.Run()
	time.Sleep(50 * time.Millisecond)

	mockSender.AssertCalled(t, "Send", mock.Anything)
	mockSender.AssertNumberOfCalls(t, "Send", 2) // Должен успеть отправить 2 раза
}

func TestNewAgent(t *testing.T) {
	mockSender := new(MockSender)
	pollInterval := 1 * time.Second
	reportInterval := 2 * time.Second

	agent := NewAgent(pollInterval, reportInterval, mockSender)

	assert.Equal(t, pollInterval, agent.pollInterval)
	assert.Equal(t, reportInterval, agent.reportInterval)
	assert.Equal(t, mockSender, agent.sender)
	assert.NotNil(t, agent.metrics)
}
