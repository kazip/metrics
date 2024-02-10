package main

import (
	"net/http"
)

func main() {

	config := Config{}
	parseFlags(&config)

	storage := NewMemStorage()

	err := http.ListenAndServe(config.ListenURI, Router(storage))
	if err != nil {
		panic(err)
	}
}
