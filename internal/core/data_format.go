package core

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// formatSchema formats extracted data into the provided schema (dummy implementation).
func formatSchema(extractedData []string, flattenedSchema map[string]interface{}) (map[string]interface{}, []string, error) {
	missedKeys := []string{}
	mergedData := make(map[string]interface{})

	// Process each extracted data item
	for _, data := range extractedData {
		parsed := parseToMap(data)
		for key, value := range parsed {
			if expectedType, exists := flattenedSchema[key]; exists {
				convertedValue, err := validateAndConvert(value, expectedType)
				if err != nil {
					return nil, nil, fmt.Errorf("error validating key '%s' with value '%v': %v", key, value, err)
				}

				// Update only if:
				// - Key does not exist in mergedData
				// - Current value in mergedData is "N/A"
				// - Current value in mergedData is empty
				if currentValue, exists := mergedData[key]; !exists || currentValue == "N/A" || currentValue == "" {
					mergedData[key] = convertedValue
				}
			} else {
				// Track keys that aren't in the schema
				missedKeys = append(missedKeys, key)
			}
		}
	}

	// Unflatten the merged data before returning
	return mergedData, missedKeys, nil
}

func parseToMap(response string) map[string]string {
	result := make(map[string]string)

	lines := strings.Split(response, ";")

	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result[key] = value
		}
	}

	return result
}

// Validate and convert value based on the expected type
func validateAndConvert(value string, expectedType interface{}) (interface{}, error) {
	if value == "N/A" {
		return "N/A", nil
	}
	switch expectedType.(type) {
	case string:
		if expectedType == "string" {
			return value, nil
		}
	case bool:
		if expectedType == "boolean" {
			convertedValue, err := strconv.ParseBool(value)
			if err != nil {
				return nil, fmt.Errorf("cannot be converted to boolean")
			}
			return convertedValue, nil
		}
	case map[string]interface{}:
		var jsonObject map[string]interface{}
		err := json.Unmarshal([]byte(value), &jsonObject)
		if err != nil {
			return nil, fmt.Errorf("cannot be parsed as JSON object")
		}
		return jsonObject, nil
	case []interface{}:
		var jsonArray []interface{}
		parts := strings.Split(value, ",")
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			jsonArray = append(jsonArray, trimmed)
		}

		return jsonArray, nil
	}

	if expectedType == "number" {
		convertedValue, err := strconv.ParseFloat(strings.ReplaceAll(value, ",", "."), 64)
		if err != nil {
			return nil, fmt.Errorf("cannot be converted to a number")
		}
		return convertedValue, nil
	}
	if expectedType == "float" {
		convertedValue, err := strconv.ParseFloat(strings.ReplaceAll(value, ",", "."), 64)
		if err != nil || !strings.Contains(value, ".") {
			return nil, fmt.Errorf("value '%s' must be a valid floating-point number", value)
		}
		return convertedValue, nil
	}

	return nil, fmt.Errorf("is of unsupported type: %v", expectedType)
}
