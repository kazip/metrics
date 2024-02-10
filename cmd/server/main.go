package main

import (
	"flag"
	"net/http"
	"os"
)

type Config struct {
	ListenURI string
}

func main() {

	listenAddr := flag.String("a", "localhost:8080", "listen address")
	flag.Parse()

	if envListenAddr := os.Getenv("ADDRESS"); envListenAddr != "" {
		*listenAddr = envListenAddr
	}

	storage := NewMemStorage()
	config := Config{ListenURI: *listenAddr}
	err := http.ListenAndServe(config.ListenURI, Router(storage))
	if err != nil {
		panic(err)
	}
}
