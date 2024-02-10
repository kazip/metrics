package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var mutex sync.RWMutex = sync.RWMutex{}

func main() {
	var wg sync.WaitGroup

	pollInterval := flag.Int("p", 2, "poll interval")
	reportInterval := flag.Int("r", 10, "report interval")
	serverAddress := flag.String("a", "localhost:8080", "server address")
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	doneMonitor := make(chan bool, 1)
	doneSender := make(chan bool, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		sig := <-sigs
		fmt.Println()
		fmt.Printf("Received signal: %s\n", sig)
		fmt.Println("Shutdown initiated...")
		doneMonitor <- true
		doneSender <- true
	}()

	metrics := &Metrics{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		NewMonitor(*pollInterval, metrics, doneMonitor, &MemStats{})
		fmt.Println("Shutdown monitor")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		NewSender(*reportInterval, *serverAddress, metrics, doneSender, &http.Client{Timeout: time.Duration(1) * time.Second})
		fmt.Println("Shutdown sender")
	}()

	wg.Wait()
	fmt.Println("Shutdown success")
}
