package helpers

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func RunPerformanceTest(config RequestConfig) {
	var wg sync.WaitGroup
	var currentVelocity int
	totalStart := time.Now()
	totalDuration := getTotalDuration(config.Phases)

	var results []RequestResult

	for _, phase := range config.Phases {
		increment := float64(phase.RequestsPerSecond-currentVelocity) / float64(phase.Duration)

		for i := 0; i < phase.Duration; i++ {
			reqsPerSecond := currentVelocity + int(float64(i)*increment)

			if reqsPerSecond < 0 {
				reqsPerSecond = 0
			}

			percentage := (float64(i) + getPhaseOffset(phase, config.Phases)) / totalDuration * 100

			fmt.Printf("Current Velocity: %d requests/second, %.2f%% of test duration passed\n", reqsPerSecond, percentage)

			for j := 0; j < reqsPerSecond; j++ {
				wg.Add(1)
				go func(requestNum int) {
					defer wg.Done()

					start := time.Now()
					client := &http.Client{Timeout: time.Duration(config.TimeoutSec) * time.Second}
					req, err := http.NewRequest(config.Method, config.URL, nil)
					if err != nil {
						fmt.Printf("Request %d failed to create: %s\n", requestNum, err)
						return
					}

					resp, err := client.Do(req)
					if err != nil {
						fmt.Printf("Request %d failed: %s\n", requestNum, err)
						return
					}
					resp.Body.Close()
					elapsed := time.Since(start)
					fmt.Printf("Request took %s\n", elapsed)

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
	}

	wg.Wait()
	totalElapsed := time.Since(totalStart)
	fmt.Printf("Total test duration: %.2fs\n", totalElapsed.Seconds())

	filename := "output.json"
	err := SaveResults(results, filename)
	if err != nil {
		fmt.Printf("Error saving results: %v\n", err)
	}
}

// computes statistics for response codes over time
func ComputeResponseCodeStats(results []RequestResult) map[int]int {
	responseCodeStats := make(map[int]int)

	for _, result := range results {
		responseCodeStats[result.ResponseCode]++
	}

	return responseCodeStats
}

// calculates the total duration of all phases
func getTotalDuration(phases []Phase) float64 {
	var totalDuration float64
	for _, phase := range phases {
		totalDuration += float64(phase.Duration)
	}
	return totalDuration
}

// calculates the offset of the current phase within the total duration
func getPhaseOffset(currentPhase Phase, phases []Phase) float64 {
	var offset float64
	for _, phase := range phases {
		if phase == currentPhase {
			break
		}
		offset += float64(phase.Duration)
	}
	return offset
}
