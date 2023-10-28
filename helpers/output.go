package helpers

import (
	"encoding/json"
	"fmt"
	"os"
)

// saves the performance test results to a JSON file.
func SaveResults(results []RequestResult, filename string) error {
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
