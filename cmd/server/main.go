package main

import (
	"net/http"
)

type Config struct {
	ListenURI string
}

func main() {

	storage := NewMemStorage()
	config := Config{ListenURI: ":8080"}
	err := http.ListenAndServe(config.ListenURI, Router(storage))
	if err != nil {
		panic(err)
	}
}
