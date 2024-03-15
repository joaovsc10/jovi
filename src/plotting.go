package src

import (
	"fmt"
	"sort"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func GenerateResponseTimePlot(results []RequestResult) {
	responseTimeBySecond := make(map[int64][]float64)
	for _, req := range results {
		roundedTimestamp := req.Timestamp.Round(time.Second)
		second := roundedTimestamp.Unix()
		responseTimeBySecond[second] = append(responseTimeBySecond[second], float64(req.Duration))
	}

	var seconds []int64
	for second := range responseTimeBySecond {
		seconds = append(seconds, second)
	}
	sort.Slice(seconds, func(i, j int) bool {
		return seconds[i] < seconds[j]
	})

	var avgResponseTimePoints plotter.XYs
	for _, second := range seconds {
		responseTimes := responseTimeBySecond[second]
		sum := 0.0
		for _, rt := range responseTimes {
			sum += rt
		}
		avgResponseTime := sum / float64(len(responseTimes))
		avgResponseTimePoints = append(avgResponseTimePoints, plotter.XY{
			X: float64(second),
			Y: avgResponseTime,
		})
	}

	p := plot.New()

	p.Title.Text = "Average Response Time Evolution Over Time"
	p.X.Label.Text = "Time (seconds)"
	p.Y.Label.Text = "Average Response Time (ms)"
	p.X.Tick.Marker = plot.TimeTicks{Format: "15:04:05"}

	avgResponseTimeLine, err := plotter.NewLine(avgResponseTimePoints)
	if err != nil {
		fmt.Println("Error creating line plot:", err)
		return
	}
	avgResponseTimeLine.LineStyle.Width = vg.Points(1)
	avgResponseTimeLine.LineStyle.Color = plotutil.Color(0)
	p.Add(avgResponseTimeLine)

	if err := p.Save(6*vg.Inch, 4*vg.Inch, "report/response_time_evolution.png"); err != nil {
		fmt.Println("Error saving plot:", err)
		return
	}
}

func GenerateStatusCodePlot(results []RequestResult) {
	statusCodeCounts := make(map[int]map[time.Time]int)

	for _, req := range results {
		statusCode := req.ResponseCode
		timestamp := req.Timestamp

		if _, ok := statusCodeCounts[statusCode]; !ok {
			statusCodeCounts[statusCode] = make(map[time.Time]int)
		}
		statusCodeCounts[statusCode][timestamp]++
	}

	p := plot.New()

	p.Title.Text = "Status Code Distribution Over Time"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Count"

	var series []interface{}
	for code, counts := range statusCodeCounts {
		var pts plotter.XYs
		for timestamp, count := range counts {
			pts = append(pts, plotter.XY{X: float64(timestamp.Unix()), Y: float64(count)})
		}
		sort.Slice(pts, func(i, j int) bool { return pts[i].X < pts[j].X })
		series = append(series, fmt.Sprintf("%d", code), pts)
	}

	p.Legend.Top = true
	p.Legend.Left = false
	p.Legend.YOffs = -10
	p.Legend.XOffs = -10
	p.Legend.Padding = 5
	p.Legend.TextStyle.Font.Size = 8

	p.X.Tick.Marker = plot.TimeTicks{Format: "15:04:05"}

	if err := plotutil.AddLines(p, series...); err != nil {
		fmt.Println("Error adding series to plot:", err)
		return
	}

	if err := p.Save(6*vg.Inch, 4*vg.Inch, "report/status_code_distribution.png"); err != nil {
		fmt.Println("Error saving plot:", err)
		return
	}

}

func GenerateRequestsPerSecondPlot(results []RequestResult, testDuration float64) {

	requestsPerSecond := make(map[int]int)
	for _, result := range results {
		second := int(result.Timestamp.Sub(results[0].Timestamp).Seconds())

		requestsPerSecond[second]++
	}

	var pts plotter.XYs
	for second := 0; second < int(testDuration); second++ {
		rps := float64(requestsPerSecond[second])

		pts = append(pts, plotter.XY{X: float64(second), Y: rps})
	}

	p := plot.New()

	p.Title.Text = "Requests Per Second (RPS) Over Time"
	p.X.Label.Text = "Time (seconds)"
	p.Y.Label.Text = "Requests Per Second (RPS)"

	line, err := plotter.NewLine(pts)
	if err != nil {
		fmt.Println("Error creating line plot:", err)
		return
	}

	p.Add(line)

	if err := p.Save(6*vg.Inch, 4*vg.Inch, "report/requests_per_second.png"); err != nil {
		fmt.Println("Error saving plot:", err)
		return
	}
}
