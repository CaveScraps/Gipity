package main;

import (
    "bytes"
    "fmt"
    "os"
    "encoding/json"
    "net/http"
    "io"
    "strings"
    "github.com/joho/godotenv"
)

const apiURL = "https://api.openai.com/v1/completions"

type OpenAIRequest struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

type OpenAIResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func main() {
    err := godotenv.Load()
    if err != nil {
        fmt.Println("Error loading .env file")
        os.Exit(1)
    }

    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        fmt.Println("OPENAI_API_KEY is not set")
        os.Exit(1)
    }

	if(len(os.Args) < 2){
		fmt.Println("Prompt is required")
		os.Exit(1)
	}

    prompt := os.Args[1]

	requestBody, err := json.Marshal(OpenAIRequest{
		Model:     "gpt-3.5-turbo-instruct",
		Prompt:    prompt,
		MaxTokens: 500,
	})
	if err != nil {
		fmt.Printf("Error marshalling request body: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received non-200 response status: %d\n", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Response body:", string(body))
		return
	}


	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	var openAIResponse OpenAIResponse
	err = json.Unmarshal(body, &openAIResponse)
	if err != nil {
		fmt.Printf("Error unmarshalling response body: %v\n", err)
		return
	}

	if len(openAIResponse.Choices) > 0 {
		responseText := openAIResponse.Choices[0].Text
		fmt.Println(strings.TrimSpace(responseText))
	} else {
		fmt.Println("No response received")
	}
}
