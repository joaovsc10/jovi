package src

import (
	"fmt"
	"strconv"
	"strings"
)

func GetTotalDuration(phases []Phase) float64 {
	var totalDuration float64
	for _, phase := range phases {
		totalDuration += float64(phase.Duration)
	}
	return totalDuration
}

func GenerateStatistics(config RequestConfig) {
	results, err := ReadResults("report/output.json")
	if err != nil {
		fmt.Printf("Error reading results: %v\n", err)
		return
	}

	responseCodeStats := ComputeResponseCodeStats(results)
	totalDuration := GetTotalDuration(config.Phases)

	passed := checkSuccessRate(config, responseCodeStats, results) && checkResponseTime(config, results)

	passFailInfo := "Test Passed"
	passFailColor := "\033[32m" // Green color
	if !passed {
		passFailInfo = "Test Failed"
		passFailColor = "\033[31m" // Red color
	}

	GenerateHtml(results, totalDuration, config, computeSuccessRate(responseCodeStats, results, config), passed)

	fmt.Printf("\n%s%s\033[0m\n", passFailColor, passFailInfo)

	printResponseCodeStatistics(responseCodeStats)
	fmt.Println("\nTotal Duration:", totalDuration)
}

func checkSuccessRate(config RequestConfig, responseCodeStats map[int]int, results []RequestResult) bool {
	actualStatusCodeCount, ok := responseCodeStats[config.ExpectedStatus]
	if !ok {
		return false
	}

	successRate := float64(actualStatusCodeCount) / float64(len(results)) * 100

	expectedSuccessRate, err := parseThreshold(config.PassingConditions.SuccessRate)
	if err != nil {
		return false
	}

	return checkThreshold(config.PassingConditions.SuccessRate, successRate, expectedSuccessRate)
}

func checkResponseTime(config RequestConfig, results []RequestResult) bool {
	_, p95ResponseTime, _, _, _ := ComputeResponseTimeStats(results)

	expectedResponseTime95, err := parseThreshold(config.PassingConditions.ResponseTime95)
	if err != nil {
		return false
	}

	return checkThreshold(config.PassingConditions.ResponseTime95, p95ResponseTime, expectedResponseTime95)
}

func parseThreshold(threshold string) (float64, error) {
	thresholdValue, err := strconv.ParseFloat(strings.TrimPrefix(threshold, "<>"), 64)
	if err != nil {
		return 0, err
	}

	return thresholdValue, nil
}

func checkThreshold(condition string, actual, expected float64) bool {
	switch {
	case strings.HasPrefix(condition, ">"):
		return actual > expected
	case strings.HasPrefix(condition, "<"):
		return actual < expected
	default:
		return false
	}
}

func printResponseCodeStatistics(responseCodeStats map[int]int) {
	totalRequests := 0
	for _, count := range responseCodeStats {
		totalRequests += count
	}

	fmt.Println("Response Code Statistics:")
	fmt.Println("+------------+-------+------------+")
	fmt.Println("| Status Code| Count | Percentage |")
	fmt.Println("+------------+-------+------------+")
	for statusCode, count := range responseCodeStats {
		percentage := float64(count) / float64(totalRequests) * 100
		fmt.Printf("| %-11d| %-6d| %-10.2f%%|\n", statusCode, count, percentage)
	}
	fmt.Println("+------------+-------+------------+")
}

func computeSuccessRate(responseCodeStats map[int]int, results []RequestResult, config RequestConfig) float64 {
	actualStatusCodeCount := responseCodeStats[config.ExpectedStatus]
	return float64(actualStatusCodeCount) / float64(len(results)) * 100
}
