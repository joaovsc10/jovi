package main

import (
	"flag"
	"fmt"
	"jovi/helpers"
)

func main() {
	jsonFile := flag.String("config", "", "Path to the JSON configuration file")
	flag.Parse()

	if *jsonFile == "" {
		fmt.Println("Please provide a JSON configuration file using the -config flag.")
		return
	}

	config, err := helpers.ReadJSONConfig(*jsonFile)
	if err != nil {
		fmt.Printf("Error reading JSON configuration: %s\n", err)
		return
	}

	helpers.RunPerformanceTest(config)
}
