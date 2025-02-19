package core

import (
	"fmt"
	"strings"
	"sync"
)

type ConfigProvider interface {
	GetEndpoint() string
	GetModel() string
	GetAPIHeader() string
	GetAPIKey() string
	GetMaxConcurrentRequests() int
	GetMaxRetries() int
	GetMaxChunkSize() int
}

func DataProcessing(conf ConfigProvider, data string, schema map[string]interface{}) (map[string]interface{}, error) {
	chunks := chunkText(conf.GetMaxChunkSize(), data)
	if len(chunks) > 10 {
		return nil, fmt.Errorf("The content exceeds the allowed size of characters per request")
	}

	type chunkResult struct {
		extracted string
		err       error
	}

	results := make(chan chunkResult, len(chunks))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, conf.GetMaxConcurrentRequests()) // Limit concurrency

	flattened := FlattenObject(schema, "")

	// Prepare keys for extraction
	keys := make([]string, 0, len(flattened))
	for k := range flattened {
		keys = append(keys, k)
	}

	// Extract fields from chunks
	for _, chunk := range chunks {
		wg.Add(1)
		go func(chunk string) {
			defer wg.Done()

			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			extracted, err := aiRequestWithRetries(conf, chunk, strings.Join(keys, ","))
			results <- chunkResult{extracted: extracted, err: err}
		}(chunk)
	}

	// Collect extracted results
	go func() {
		wg.Wait()
		close(results)
	}()

	var extractedData []string

	for result := range results {
		if result.err != nil {
			return nil, fmt.Errorf("errors encountered: %v", result.err)
		}
		extractedData = append(extractedData, result.extracted)
	}

	// Format extracted data
	mergedData, missedKeys, err := formatSchema(extractedData, flattened)
	if err != nil {
		return nil, err
	}

	// Retry extraction for missing keys with goroutines
	for attempt := 0; attempt < 3 && len(missedKeys) > 0; attempt++ {
		missingResults := make(chan chunkResult, len(chunks))

		// Run missing key extraction concurrently
		for _, chunk := range chunks {
			wg.Add(1)
			go func(chunk string) {
				defer wg.Done()

				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				extractedMissingData, err := aiRequestWithRetries(conf, chunk, strings.Join(missedKeys, ","))
				missingResults <- chunkResult{extracted: extractedMissingData, err: err}
			}(chunk)
		}

		go func() {
			wg.Wait()
			close(missingResults)
		}()

		// Collect missing data
		var newExtractedData []string
		for result := range missingResults {
			if result.err != nil {
				return nil, fmt.Errorf("errors encountered: %v", result.err)
			}
			newExtractedData = append(newExtractedData, result.extracted)
		}

		// Call formatSchema again with the newly extracted data
		extractedData = append(extractedData, newExtractedData...)
		mergedData, missedKeys, err = formatSchema(extractedData, flattened)
		if err != nil {
			return nil, err
		}

		// Break if all keys are found
		if len(missedKeys) == 0 {
			break
		}
	}

	// Assign "N/A" for any remaining missing keys
	for _, key := range missedKeys {
		mergedData[key] = "N/A"
	}

	return unflattenObject(mergedData, schema), nil
}
