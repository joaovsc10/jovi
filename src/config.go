package src

import (
	"encoding/json"
	"os"
	"time"
)

type RequestConfig struct {
	Method            string            `json:"method"`
	URL               string            `json:"url"`
	TimeoutSec        int               `json:"timeoutSec"`
	Phases            []Phase           `json:"phases"`
	ExpectedStatus    int               `json:"expectedStatus"`
	PassingConditions PassingConditions `json:"passingConditions"`
}

type PassingConditions struct {
	SuccessRate    string `json:"successRate"`
	ResponseTime95 string `json:"responseTime95"`
}

type Phase struct {
	Duration          int `json:"duration"`
	RequestsPerSecond int `json:"requestsPerSecond"`
}

type RequestResult struct {
	Duration     time.Duration
	ResponseCode int
	DataSent     int64 `json:"data_sent"`
	DataReceived int64 `json:"data_received"`
	Timestamp    time.Time
}

func ReadJSONConfig(filePath string) (RequestConfig, error) {
	var config RequestConfig
	file, err := os.Open(filePath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
