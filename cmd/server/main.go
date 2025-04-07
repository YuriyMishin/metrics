package main

import (
	"YuriyMishin/metrics/internal/server"
	"flag"
	"fmt"
	"log"
	"os"
)

func parseServerFlags() (string, error) {
	defaultAddr := "localhost:8080"
	envAddr := os.Getenv("ADDRESS")

	var flagAddr string
	flag.StringVar(&flagAddr, "a", defaultAddr, "HTTP server endpoint address")

	flag.Parse()

	if flag.NArg() > 0 {
		return "", fmt.Errorf("unknown flags: %v", flag.Args())
	}
	addr := defaultAddr
	if flagAddr != defaultAddr {
		addr = flagAddr
	}
	if envAddr != "" {
		addr = envAddr
	}

	return addr, nil
}

func main() {
	addr, err := parseServerFlags()
	if err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}
	srv := server.NewServer()
	if err := srv.Start(addr); err != nil {
		panic(err)
	}
}
