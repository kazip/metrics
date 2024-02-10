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

func NewSender(reportInterval int, serverAddress string, metrics *Metrics, done <-chan bool, client HTTPClient) {
	var interval = time.Duration(reportInterval) * time.Second

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
			resp, err := client.Post(fmt.Sprintf("http://%s/update/gauge/%s/%f", serverAddress, gauge.name, gauge.value), "plain/text", nil)
			if err != nil {
				fmt.Println(err)
			} else {
				p := []byte{}
				_, _ = resp.Body.Read(p)
				resp.Body.Close() // wtf?
			}
		}

		for _, counter := range metricsCopy.counters {
			fmt.Printf("%s: %d\n", counter.name, counter.value)
			resp, err := client.Post(fmt.Sprintf("http://%s/update/counter/%s/%d", serverAddress, counter.name, counter.value), "plain/text", nil)

			if err != nil {
				fmt.Println(err)
			} else {
				p := []byte{}
				_, _ = resp.Body.Read(p)
				resp.Body.Close() // wtf?
			}

		}

	}
}
