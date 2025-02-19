package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Dcup-dev/Dcup-lib/internal/core"
)

func clean(buf bytes.Buffer) (string, error) {
	var jsonData map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &jsonData); err != nil {
		return "", fmt.Errorf("failed to parse JSON content: %v", err)
	}

	flatten_json := core.FlattenObject(jsonData, "")
	return MapToKeyValueString(flatten_json), nil

}

func MapToKeyValueString(data map[string]interface{}) string {
	var result strings.Builder

	for key, value := range data {
		valueStr := fmt.Sprintf("%v", value)
		result.WriteString(fmt.Sprintf("%s:%s; ", key, valueStr))
	}

	return strings.TrimSpace(result.String())
}
