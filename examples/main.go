package main

import (
	"fmt"

	"github.com/Dcup-dev/Dcup-lib/cmd/dcup"
)

func main() {
	config := dcup.Config{
		Endpoint:  "https://api.example.com",
		Model:     "gpt-4",
		APIHeader: "Authorization",
		APIKey:    "sk-xxx",
	}

	if _, err := dcup.Init(config); err != nil {
		fmt.Println("Error:", err)
		return
	}
}
