package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type AgentConfig struct {
	Addr           string
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func NewAgentConfig() (*AgentConfig, error) {
	return parseAgentFlags()
}

func parseAgentFlags() (*AgentConfig, error) {
	var (
		flagAddr   string
		flagPoll   int
		flagReport int
	)
	defaultAddr := "localhost:8080"
	defaultPoll := 2 * time.Second
	defaultReport := 10 * time.Second

	envAddr := os.Getenv("ADDRESS")
	envPoll := os.Getenv("POLL_INTERVAL")
	envReport := os.Getenv("REPORT_INTERVAL")

	flag.StringVar(&flagAddr, "a", defaultAddr, "HTTP server endpoint address")
	flag.IntVar(&flagPoll, "p", int(defaultPoll.Seconds()), "Poll interval in seconds")
	flag.IntVar(&flagReport, "r", int(defaultReport.Seconds()), "Report interval in seconds")

	flag.Parse()

	if flag.NArg() > 0 {
		return nil, fmt.Errorf("unknown flags: %v", flag.Args())
	}

	config := &AgentConfig{
		Addr:           defaultAddr,
		PollInterval:   defaultPoll,
		ReportInterval: defaultReport,
	}

	if flagAddr != defaultAddr {
		config.Addr = flagAddr
	}
	if flagPoll != int(defaultPoll.Seconds()) {
		config.PollInterval = time.Duration(flagPoll) * time.Second
	}
	if flagReport != int(defaultReport.Seconds()) {
		config.ReportInterval = time.Duration(flagReport) * time.Second
	}

	if envAddr != "" {
		config.Addr = envAddr
	}
	if envPoll != "" {
		poll, err := strconv.Atoi(envPoll)
		if err != nil || poll <= 0 {
			return nil, fmt.Errorf("invalid POLL_INTERVAL value: %s", envPoll)
		}
		config.PollInterval = time.Duration(poll) * time.Second
	}
	if envReport != "" {
		report, err := strconv.Atoi(envReport)
		if err != nil || report <= 0 {
			return nil, fmt.Errorf("invalid REPORT_INTERVAL value: %s", envReport)
		}
		config.ReportInterval = time.Duration(report) * time.Second
	}

	if config.PollInterval <= 0 {
		return nil, fmt.Errorf("poll interval must be positive")
	}
	if config.ReportInterval <= 0 {
		return nil, fmt.Errorf("report interval must be positive")
	}

	return config, nil
}

type ServerConfig struct {
	Addr string
}

func NewServerConfig() (*ServerConfig, error) {
	return parseServerFlags()
}

func parseServerFlags() (*ServerConfig, error) {
	defaultAddr := "localhost:8080"
	envAddr := os.Getenv("ADDRESS")

	var flagAddr string
	flag.StringVar(&flagAddr, "a", defaultAddr, "HTTP server endpoint address")

	flag.Parse()

	config := &ServerConfig{
		Addr: defaultAddr,
	}

	if flag.NArg() > 0 {
		return config, fmt.Errorf("unknown flags: %v", flag.Args())
	}

	if flagAddr != defaultAddr {
		config.Addr = flagAddr
	}
	if envAddr != "" {
		config.Addr = envAddr
	}

	return config, nil
}
