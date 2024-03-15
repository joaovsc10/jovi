package src

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func RunPerformanceTest(config RequestConfig) {
	var wg sync.WaitGroup
	var currentVelocity int
	var previousPhaseDuration int

	totalStart := time.Now()
	totalDuration := GetTotalDuration(config.Phases)

	client := &http.Client{Timeout: time.Duration(config.TimeoutSec) * time.Second}

	var results []RequestResult
	var failedRequests int

	for _, phase := range config.Phases {
		increment := float64(phase.RequestsPerSecond-currentVelocity) / float64(phase.Duration)

		for i := 0; i < phase.Duration; i++ {
			reqsPerSecond := currentVelocity + int(float64(i)*increment)

			if reqsPerSecond < 0 {
				reqsPerSecond = 0
			}

			percentage := (float64(i) + float64(previousPhaseDuration)) / totalDuration * 100

			fmt.Printf("Current Velocity: %d requests/second, %.2f%% of test duration passed\n", reqsPerSecond, percentage)

			sem := make(chan struct{}, 1000)

			for j := 0; j < reqsPerSecond; j++ {
				sem <- struct{}{}
				wg.Add(1)
				go func(requestNum int) {
					defer func() {
						<-sem
						wg.Done()
					}()

					req, err := http.NewRequest(config.Method, config.URL, nil)
					if err != nil {
						fmt.Printf("Request %d failed to create: %s\n", requestNum, err)
						failedRequests++
						return
					}

					start := time.Now()
					resp, err := client.Do(req)
					elapsed := time.Since(start)

					if err != nil {
						// fmt.Printf("Request %d failed: %s\n", requestNum, err)
						failedRequests++
						return
					}
					resp.Body.Close()
					//fmt.Printf("Request took %s\n", elapsed)

					result := RequestResult{
						Duration:     time.Duration(elapsed.Milliseconds()),
						ResponseCode: resp.StatusCode,
						DataSent:     calculateDataSent(req),
						DataReceived: calculateDataReceived(resp),
						Timestamp:    time.Now(),
					}
					results = append(results, result)
				}(j)
			}

			time.Sleep(time.Second)
		}
		currentVelocity = phase.RequestsPerSecond
		previousPhaseDuration += phase.Duration
	}

	wg.Wait()
	totalElapsed := time.Since(totalStart)
	fmt.Printf("\nTotal test duration: %.2fs\n", totalElapsed.Seconds())

	filename := "report/output.json"
	err := SaveResults(results, filename)
	if err != nil {
		fmt.Printf("Error saving results: %v\n", err)
	}
}
