package main

import (
	"net/http"
)

type Config struct {
	ListenURI string
}

type HandleError struct {
	msg    string
	status int
}

func (e HandleError) Error() string {
	return e.msg
}

func (e HandleError) Status() int {
	return e.status
}

func main() {
	
	storage := NewMemStorage()
	config := Config{ListenURI: ":8080"}

	mux := http.NewServeMux()
	mux.HandleFunc("/update/counter/", handleCounterFunc(storage))
	mux.HandleFunc("/update/gauge/", handleGaugeFunc(storage))
	mux.HandleFunc("/update/", handleUnknownType)
	err := http.ListenAndServe(config.ListenURI, mux)
	if err != nil {
		panic(err)
	}
}
