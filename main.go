package main

import (
	"flag"
	"fmt"
	"jovi/src"
)

func main() {
	jsonFile := flag.String("config", "", "Path to the JSON configuration file")
	flag.Parse()

	if *jsonFile == "" {
		fmt.Println("Please provide a JSON configuration file using the -config flag.")
		return
	}

	config, err := src.ReadJSONConfig(*jsonFile)
	if err != nil {
		fmt.Printf("Error reading JSON configuration: %s\n", err)
		return
	}

	src.RunPerformanceTest(config)

	src.GenerateStatistics(config)
}
