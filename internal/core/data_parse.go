package core

import (
	"fmt"
	"strings"
	"time"
)

const objectDelimiter = "=>"

func  aiRequestWithRetries(conf ConfigProvider,data string, keys string) (string, error) {
	baseDelay := time.Second

	for i := 0; i < conf.GetMaxRetries(); i++ {

		response, err := aiRequest(conf, data, keys)

		if err != nil {
			// Check for rate-limit error
			if strings.Contains(err.Error(), "Rate limit") {
				waitTime := baseDelay * time.Duration(1<<i)
				time.Sleep(waitTime)
				continue
			}

			// For other errors, return immediately
			return "", err
		}

		return response, nil
	}

	return "", fmt.Errorf("maximum retries exceeded for OpenAI request")
}

func FlattenObject(input map[string]interface{}, prefix string) map[string]interface{} {
	flat := make(map[string]interface{})
	for key, value := range input {
		if nestedMap, ok := value.(map[string]interface{}); ok {
			nestedFlat := FlattenObject(nestedMap, prefix+key+objectDelimiter)
			for nestedKey, nestedValue := range nestedFlat {
				flat[nestedKey] = nestedValue
			}
		} else {
			flat[prefix+key] = value
		}
	}
	return flat
}

func unflattenObject(flat map[string]interface{}, schema map[string]interface{}) map[string]interface{} {
	nested := make(map[string]interface{})

	for key, value := range flat {
		parts := strings.Split(key, objectDelimiter)

		if !isInSchema(key, schema) {
			nested[key] = value
			continue
		}

		// Handle nested keys
		currentMap := nested
		for _, part := range parts[:len(parts)-1] {
			// Traverse or create intermediate maps
			if _, exists := currentMap[part]; !exists {
				currentMap[part] = make(map[string]interface{})
			}
			currentMap = currentMap[part].(map[string]interface{})
		}
		// Add the value at the final level
		currentMap[parts[len(parts)-1]] = value
	}

	return nested
}

func isInSchema(key string, schema map[string]interface{}) bool {
	parts := strings.Split(key, "=>")
	current := schema
	for _, part := range parts {
		// Traverse schema for nested keys
		if nestedSchema, exists := current[part]; exists {
			if nested, ok := nestedSchema.(map[string]interface{}); ok {
				current = nested
			} else {
				return true // Reached a leaf
			}
		} else {
			return false // Key not in schema
		}
	}
	return true
}
