package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Struct to hold the relevant fields
type message struct {
	Content string `json:"content"`
}

type choice struct {
	Message message `json:"message"`
}

type response struct {
	Choices []choice `json:"choices"`
}

// Struct to hold the error response
type errorDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type errorResponse struct {
	Error errorDetails `json:"error"`
}



func aiRequest(conf ConfigProvider, data string, keys string) (string, error) {

	requestBody := map[string]interface{}{
		"messages": []map[string]string{
			{
				"role": "system",
				"content": `You are an intelligent, detail-oriented data extraction assistant trained to identify and organize information from unstructured text. Your task is to extract data fields with high accuracy and provide results in a structured format. Always ensure you:
- Analyze the entire context of the text.
- Prioritize accuracy and completeness when extracting requested fields.
- Avoid skipping any requested fields.
- If data is presented in tables or lists, extract all possible values systematically.

Output only the requested fields in the exact format: <field>: <value>; separate them with a semicolon (;). If a field is missing or unclear, return 'N/A' as the value.`,
			},
			{
				"role": "user",
				"content": fmt.Sprintf(`Extract the following fields from the provided text while maintaining accuracy and capturing all relevant data. For fields presented in tables, extract all possible values and include them appropriately. Do not skip any requested fields. 
Fields to extract: { %s }. 
Text: { %s }`, keys, data),
			},
		},
		"model":       conf.GetModel(),
		"temperature": 0.3,
		"max_tokens":  500,
		"top_p":       1,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("Error marshaling JSON: %v", err)
	}

	r, err := http.NewRequest("POST",conf.GetEndpoint(), bytes.NewBuffer([]byte(body)))

	if err != nil {
		return "", err
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add(conf.GetAPIHeader(), conf.GetAPIKey())

	client := &http.Client{}
	res, err := client.Do(r)

	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	// Read and print the response body as JSON
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body: %v\n", err)
	}

	if res.StatusCode != http.StatusOK {
		var errorResponse errorResponse
		err := json.Unmarshal(responseBody, &errorResponse)
		if err != nil {
			return "", fmt.Errorf("Failed to parse error response.")
		}

		return "", fmt.Errorf("%s", errorResponse.Error.Message)
	}

	var successResponse response
	if err := json.Unmarshal(responseBody, &successResponse); err != nil || len(successResponse.Choices) == 0 {
		return "", fmt.Errorf("Failed to parse error response.")
	}
	return successResponse.Choices[0].Message.Content, nil
}
