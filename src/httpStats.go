package src

import (
	"io"
	"net/http"
	"sort"
)

func calculateDataSent(req *http.Request) int64 {
	if req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err == nil {
			return int64(len(body))
		}
	}
	return 0
}

func calculateDataReceived(resp *http.Response) int64 {
	if resp.Body != nil {
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			return int64(len(body))
		}
	}
	return 0
}

func ComputeResponseCodeStats(results []RequestResult) map[int]int {
	responseCodeStats := make(map[int]int)

	for _, result := range results {
		responseCodeStats[result.ResponseCode]++
	}

	return responseCodeStats
}

func ComputeResponseTimeStats(results []RequestResult) (mean, p95, p99, min, max float64) {
	var responseTimes []float64
	for _, result := range results {
		responseTimes = append(responseTimes, float64(result.Duration))
	}

	sort.Float64s(responseTimes)

	total := len(responseTimes)
	if total == 0 {
		return 0, 0, 0, 0, 0
	}

	var sum float64
	for _, rt := range responseTimes {
		sum += rt
	}
	mean = sum / float64(total)

	p95Index := int(float64(total) * 0.95)
	if p95Index < total {
		p95 = responseTimes[p95Index]
	} else {
		p95 = responseTimes[total-1]
	}

	p99Index := int(float64(total) * 0.99)
	if p99Index < total {
		p99 = responseTimes[p99Index]
	} else {
		p99 = responseTimes[total-1]
	}

	min = responseTimes[0]
	max = responseTimes[total-1]

	return mean, p95, p99, min, max
}
