package main

import (
	"YuriyMishin/metrics/server"
	"flag"
	"fmt"
	"log"
)

func parseServerFlags() (string, error) {
	var addr string
	flag.StringVar(&addr, "a", "localhost:8080", "HTTP server endpoint address")

	flag.Parse()

	if flag.NArg() > 0 {
		return "", fmt.Errorf("unknown flags: %v", flag.Args())
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
