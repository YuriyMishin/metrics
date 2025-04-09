package main

import (
	"YuriyMishin/metrics/internal/agent"
	"YuriyMishin/metrics/internal/config"
	"log"
)

func main() {
	config, err := config.NewAgentConfig()

	if err != nil {
		panic(err)
	}

	agent := agent.NewAgent(config)

	log.Println("Starting agent...")
	agent.Run()
}
