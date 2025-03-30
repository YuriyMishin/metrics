package main

import (
	"YuriyMishin/metrics/agent"
	"log"
	"time"
)

func main() {
	// Создаем отправитель с RESTy клиентом
	sender := agent.NewRestySender("http://localhost:8080")

	// Настраиваем агент
	agent := agent.NewAgent(
		2*time.Second,  // Интервал сбора метрик
		10*time.Second, // Интервал отправки
		sender,
	)

	log.Println("Starting agent...")
	agent.Run()
}
