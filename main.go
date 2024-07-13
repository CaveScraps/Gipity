package main;

import (
    "bytes"
    "errors"
    "fmt"
    "os"
    "path"
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

type Config struct {
	APIKey string
	Prompt string
}

func setup() (Config, error) {
	exePath, _ := os.Executable()
	pathOfExe := path.Dir(exePath)

    err := godotenv.Load(pathOfExe + "/.gipityenv")
    if err != nil {
		return Config{}, errors.New("Error loading .env file: " + err.Error())
    }

    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
		return Config{}, errors.New("OPENAI_API_KEY is not set")
    }

	if(len(os.Args) < 2){
		return Config{}, errors.New("Prompt is required")
	}
    prompt := os.Args[1]

	return Config{APIKey: apiKey, Prompt: prompt}, nil
}

func getResponseFromOpenAI(prompt string, apiKey string) ([]byte, error){
	requestBody, err := json.Marshal(OpenAIRequest{
		Model:     "gpt-3.5-turbo-instruct",
		Prompt:    prompt,
		MaxTokens: 500,
	})
	if err != nil {
		return nil, errors.New("Error marshalling request body: " + err.Error())
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, errors.New("Error creating request: " + err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Error making request: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New("Received non-200 response status: " + resp.Status + "\n" + string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Error reading response body: " + err.Error())
	}
	return body, nil
}

func main() {
	Config, err := setup()
	if err != nil {
		fmt.Println("Error setting up: " + err.Error())
		os.Exit(1)
	}

	body, err := getResponseFromOpenAI(Config.Prompt, Config.APIKey)
	if err != nil {
		fmt.Println("Error getting response from OpenAI: " + err.Error())
		os.Exit(1)
	}

	var openAIResponse OpenAIResponse
	err = json.Unmarshal(body, &openAIResponse)
	if err != nil {
		fmt.Printf("Error unmarshalling response body: %v\n", err)
		os.Exit(1)
	}

	if len(openAIResponse.Choices) > 0 {
		responseText := openAIResponse.Choices[0].Text
		fmt.Println(strings.TrimSpace(responseText))
	} else {
		fmt.Println("No response received")
	}
}
