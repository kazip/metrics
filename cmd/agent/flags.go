package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

type Config struct {
	PollInterval   int
	ReportInterval int
	ServerAddress  string
}

func parseFlags(config *Config) {
	pollInterval := flag.Int("p", 2, "poll interval")
	reportInterval := flag.Int("r", 10, "report interval")
	serverAddress := flag.String("a", "localhost:8080", "server address")
	flag.Parse()

	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		*serverAddress = envServerAddress
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		interval, err := strconv.Atoi(envReportInterval)
		if err != nil {
			log.Fatal(err)
		}
		*reportInterval = interval
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		interval, err := strconv.Atoi(envPollInterval)
		if err != nil {
			log.Fatal(err)
		}
		*pollInterval = interval
	}

	config.PollInterval = *pollInterval
	config.ReportInterval = *reportInterval
	config.ServerAddress = *serverAddress
}
