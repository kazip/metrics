package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockHTTPClient struct {
	Urls []string
}

func (m *MockHTTPClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	m.Urls = append(m.Urls, url)
	fmt.Println(url)
	r := io.NopCloser(bytes.NewReader([]byte("ok!")))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil
}

func TestNewSender(t *testing.T) {
	type args struct {
		duration int
		metrics  *Metrics
		done     chan bool
		client   MockHTTPClient
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{
				duration: 1,
				metrics: &Metrics{
					gauges: []GaugeMetric{
						{
							name:  "Alloc",
							value: 1,
						},
					},
				},
				done: make(chan bool),
				client: MockHTTPClient{
					Urls: []string{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fmt.Println(tt.args.metrics.gauges)

			go NewSender(tt.args.duration, tt.args.metrics, tt.args.done, &tt.args.client)
			time.Sleep(time.Duration(1500) * time.Millisecond)
			tt.args.done <- true
			assert.Equal(t, 1, len(tt.args.client.Urls))
		})
	}
}
