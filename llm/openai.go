package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Config stores the API host and key.
type Config struct {
	Host     string
	APIKey   string
	Endpoint url.URL
}

func localhost() Config {
	return Config{
		Host:   "localhost:7869",
		APIKey: "api-key",
		Endpoint: url.URL{
			Scheme: "http",
			Host:   "localhost:7869",
			Path:   "/v1/chat/completions",
		},
	}
}

func defaultRequest() *Request {
	return &Request{
		Messages: []Message{
			{
				Role:    SystemRole,
				Content: "The following is a conversation with an AI assistant. The assistant is helpful, creative, clever, and very friendly.",
			},
			{
				Role:    UserRole,
				Content: "How can I help you today?",
			},
		},
		Temperature: 0.7,
		MaxTokens:   2048,
		Stream:      false,
	}
}

// inference makes a POST request to the OpenAI API with the given request data.
func (c Config) inference(requestData *Request) (*http.Response, error) {
	requestBytes, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", c.Endpoint.String(), bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+c.APIKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	return resp, nil
}

// handleResponse parses the HTTP response from the inference API call using io.ReadAll.
func handleResponse(response *http.Response) (*[]Message, error) {
	defer response.Body.Close()

	// Use io.ReadAll to read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var messages []Message
	if err := json.Unmarshal(body, &messages); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &messages, nil
}

// handleStreamedResponse processes the HTTP response from a streamed API call.
func handleStreamedResponse(response *http.Response) error {
	defer response.Body.Close()

	// Create a new buffered reader to read the response body line by line
	reader := bufio.NewReader(response.Body)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			// If we reach the end of the stream, break out of the loop
			if err.Error() == "EOF" {
				break
			}
			// For any other error, return it
			return fmt.Errorf("error reading streamed response: %w", err)
		}
		// Process the line here (for example, print it)
		fmt.Print(line)
	}

	return nil
}
