package src

import (
	"fmt"
	"os"
)

func GenerateHtml(results []RequestResult, totalDuration float64, config RequestConfig, successRate float64, passed bool) {
	GenerateStatusCodePlot(results)
	GenerateResponseTimePlot(results)
	GenerateRequestsPerSecondPlot(results, totalDuration)

	var statusCodeTable string
	statusCodeTable += "<h2>Status Code Statistics</h2>"
	statusCodeTable += "<table class='statistics-table'>"
	statusCodeTable += "<tr><th>Status Code</th><th>Count</th><th>Percentage</th></tr>"

	responseCodeStats := ComputeResponseCodeStats(results)
	totalRequests := len(results)

	for statusCode, count := range responseCodeStats {
		percentage := float64(count) / float64(totalRequests) * 100
		statusCodeTable += fmt.Sprintf("<tr><td>%d</td><td>%d</td><td>%.2f%%</td></tr>", statusCode, count, percentage)
	}

	statusCodeTable += "</table>"

	var responseTimeTable string
	responseTimeTable += "<h2>Response Time Statistics</h2>"
	responseTimeTable += "<table class='statistics-table'>"
	responseTimeTable += "<tr><th>Statistic</th><th>Value (ms)</th></tr>"

	meanResponseTime, p95ResponseTime, p99ResponseTime, minResponseTime, maxResponseTime := ComputeResponseTimeStats(results)

	responseTimeTable += fmt.Sprintf("<tr><td>Mean</td><td>%.0f</td></tr>", meanResponseTime)
	responseTimeTable += fmt.Sprintf("<tr><td>p(95)</td><td>%.0f</td></tr>", p95ResponseTime)
	responseTimeTable += fmt.Sprintf("<tr><td>p(99)</td><td>%.0f</td></tr>", p99ResponseTime)
	responseTimeTable += fmt.Sprintf("<tr><td>Min</td><td>%.0f</td></tr>", minResponseTime)
	responseTimeTable += fmt.Sprintf("<tr><td>Max</td><td>%.0f</td></tr>", maxResponseTime)

	responseTimeTable += "</table>"

	testDescription := fmt.Sprintf("<h2>Test Case Description</h2><p>Method: %s</p><p>URL: %s</p>", config.Method, config.URL)

	passFailInfo := "Test Passed"
	passFailColor := "green"
	if !passed {
		passFailInfo = "Test Failed"
		passFailColor = "red"
	}

	passFailTable := fmt.Sprintf(`
		<h2>Test Pass/Fail</h2>
		<table class='pass-fail-table'>
			<tr><th>Test Result</th><td style='color: %s;'>%s</td></tr>
		</table>
	`, passFailColor, passFailInfo)

	passingConditionsTable := fmt.Sprintf(`
		<h2>Passing Conditions</h2>
		<table class='passing-conditions-table'>
			<tr><th>Condition</th><th>Expected</th><th>Actual</th></tr>
			<tr><td>Success Rate</td><td>%s</td><td>%.2f</td></tr>
			<tr><td>Response Time 95th Percentile</td><td>%s</td><td>%.0f ms</td></tr>
		</table>
	`, config.PassingConditions.SuccessRate, successRate, config.PassingConditions.ResponseTime95, p95ResponseTime)

	htmlReport := fmt.Sprintf(`
		<html>
		<head>
		<title>Performance Test Report</title>
		<style>
		body {
			font-family: Arial, sans-serif;
			margin: 0;
			padding: 20px;
		}
		.container {
			max-width: 800px;
			margin: auto;
		}
		h1, h2 {
			text-align: center;
			margin-top: 50px;
		}
		.statistics-table {
			width: 100%%;
			border-collapse: collapse;
			margin-top: 40px;
			margin-bottom: 40px;
		}
		.statistics-table th, .statistics-table td {
			border: 1px solid #ddd;
			padding: 8px;
			text-align: left;
		}
		.statistics-table th {
			background-color: #f2f2f2;
		}
		.plot-image {
			display: block;
			margin: auto;
			margin-top: 40px;
			margin-bottom: 40px;
		}
		.pass-fail-table {
			width: 100%%;
			margin-top: 20px;
			margin-bottom: 20px;
			border-collapse: collapse;
		}
		.pass-fail-table th, .pass-fail-table td {
			border: 1px solid #ddd;
			padding: 8px;
			text-align: left;
		}
		.passing-conditions-table {
			width: 100%%;
			margin-top: 20px;
			margin-bottom: 20px;
			border-collapse: collapse;
		}
		.passing-conditions-table th, .passing-conditions-table td {
			border: 1px solid #ddd;
			padding: 8px;
			text-align: left;
		}
		.passing-conditions-table th {
			background-color: #f2f2f2;
		}
		.passing-conditions-table td:first-child {
			font-weight: bold;
		}
		</style>
		</head>
		<body>
		<div class="container">
			<h1>Performance Test Report</h1>
			%s
			%s
			%s
			%s
			%s
			<h2>Statistics Distribution Over Time</h2>
			<img class="plot-image" src="status_code_distribution.png" alt="Status Code Distribution">
			<img class="plot-image" src="response_time_evolution.png" alt="Response Time Distribution">
			<img class="plot-image" src="requests_per_second.png" alt="Requests Per Second (RPS) Over Time">
		</div>
		</body>
		</html>
	`, testDescription, passFailTable, passingConditionsTable, statusCodeTable, responseTimeTable)

	if err := os.WriteFile("report/report.html", []byte(htmlReport), 0644); err != nil {
		fmt.Println("Error writing HTML report:", err)
		return
	}

	fmt.Println("\nHTML report generated successfully.")
}
