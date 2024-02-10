package main

import (
	"flag"
	"os"
)

type Config struct {
	ListenURI string
}

func parseFlags(config *Config) {

	listenAddr := flag.String("a", "localhost:8080", "listen address")
	flag.Parse()

	if envListenAddr := os.Getenv("ADDRESS"); envListenAddr != "" {
		*listenAddr = envListenAddr
	}

	config.ListenURI = *listenAddr
}
