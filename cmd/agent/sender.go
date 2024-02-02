package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPClient interface {
	Post(string, string, io.Reader) (*http.Response, error)
}

func NewSender(duration int, metrics *Metrics, done <-chan bool, client HTTPClient) {
	var interval = time.Duration(duration) * time.Second

loop:
	for {

		select {
		case <-time.After(interval):
		case <-done:
			break loop
		}

		mutex.RLock()
		metricsCopy := *metrics
		mutex.RUnlock()

		fmt.Println("Sending metrics")

		for _, gauge := range metricsCopy.gauges {
			fmt.Printf("%s: %f\n", gauge.name, gauge.value)
			resp, err := client.Post(fmt.Sprintf("http://127.0.0.1:8989/update/gauge/%s/%f", gauge.name, gauge.value), "plain/text", nil)
			p := []byte{}
			_, _ = resp.Body.Read(p)
			resp.Body.Close() // wtf?

			if err != nil {
				fmt.Println(err)
			}
		}

		for _, counter := range metricsCopy.counters {
			fmt.Printf("%s: %d\n", counter.name, counter.value)
			resp, err := client.Post(fmt.Sprintf("http://127.0.0.1:8989/update/counter/%s/%d", counter.name, counter.value), "plain/text", nil)
			p := []byte{}
			_, _ = resp.Body.Read(p)
			resp.Body.Close() // wtf?

			if err != nil {
				fmt.Println(err)
			}
		}

	}
}
