package src

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func SaveResults(results []RequestResult, filename string) error {
	dir := filepath.Dir(filename)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	resultsJSON, err := json.Marshal(results)
	if err != nil {
		return fmt.Errorf("error serializing results: %w", err)
	}

	err = os.WriteFile(filename, resultsJSON, 0644)
	if err != nil {
		return fmt.Errorf("error saving results to file: %w", err)
	}

	fmt.Printf("Results saved to %s\n", filename)
	return nil
}

func ReadResults(filename string) ([]RequestResult, error) {
	var results []RequestResult

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	if err := json.Unmarshal(data, &results); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return results, nil
}
