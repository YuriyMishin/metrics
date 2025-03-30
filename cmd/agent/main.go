package main

import (
	"YuriyMishin/metrics/agent"
	"flag"
	"fmt"
	"log"
	"time"
)

func parseAgentFlags() (string, time.Duration, time.Duration, error) {
	var (
		addr           string
		pollInterval   int
		reportInterval int
	)

	flag.StringVar(&addr, "a", "localhost:8080", "HTTP server endpoint address")
	flag.IntVar(&pollInterval, "p", 2, "Poll interval in seconds")
	flag.IntVar(&reportInterval, "r", 10, "Report interval in seconds")

	flag.Parse()

	if flag.NArg() > 0 {
		return "", 0, 0, fmt.Errorf("unknown flags: %v", flag.Args())
	}

	if pollInterval <= 0 {
		return "", 0, 0, fmt.Errorf("poll интервал должен быть положительным")
	}
	if reportInterval <= 0 {
		return "", 0, 0, fmt.Errorf("report интервал должен быть положительным")
	}

	return addr, time.Duration(pollInterval) * time.Second,
		time.Duration(reportInterval) * time.Second, nil
}

func main() {
	addr, pollInterval, reportInterval, err := parseAgentFlags()
	if err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	// Создаем отправитель с RESTy клиентом
	sender := agent.NewRestySender("http://" + addr)

	// Настраиваем агент
	agent := agent.NewAgent(pollInterval, reportInterval, sender)

	log.Println("Starting agent...")
	agent.Run()
}
