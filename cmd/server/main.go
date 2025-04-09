package main

import (
	"YuriyMishin/metrics/internal/config"
	"YuriyMishin/metrics/internal/server"
	"log"
)

func main() {
	config, err := config.NewServerConfig()
	if err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}
	srv := server.NewServer()
	if err := srv.Start(config.Addr); err != nil {
		panic(err)
	}
}
