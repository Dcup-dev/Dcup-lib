package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Dcup-dev/Dcup-lib/cmd/dcup"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := dcup.Config{
		Endpoint:  os.Getenv("OPENAI_URL"),
		Model:     "gpt-4o-mini",
		APIHeader: "Authorization",
		APIKey:    fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_KEY")),
	}

	client, err := dcup.Init(config)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	schema := map[string]interface{}{
		"title":  "string",
		"auther": "string",
	}

	res, err := client.Docx.CleanUrl("https://calibre-ebook.com/downloads/demos/demo.docx", schema)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(res)
}
