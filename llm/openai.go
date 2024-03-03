package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Config stores the API host and key.
type Config struct {
	Host     string
	APIKey   string
	Endpoint url.URL
}

// inference makes a POST request to the OpenAI API with the given request data.
func (c Config) inference(r *Request) (*http.Response, error) {
	requestBytes, err := json.Marshal(r)
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

func (c Config) Infer(request *Request) (Response, error) {
	resp, err := c.inference(request)
	if err != nil {
		return Response{}, fmt.Errorf("failed to make inference request: %w", err)
	}

	var response Response
	if request.Stream {
		lines, err := handleStreamedResponse(resp)
		if err != nil {
			return Response{}, fmt.Errorf("failed to handle streamed response: %w", err)
		}
		response, err = UnmarshalResponse([]byte(strings.Join(lines, "")))
		if err != nil {
			return Response{}, fmt.Errorf("failed to unmarshal streamed response: %w", err)
		}
	} else {
		response, err = handleResponse(resp)
		if err != nil {
			return Response{}, fmt.Errorf("failed to handle response: %w", err)
		}
	}

	return response, nil
}

// handleResponse parses the HTTP response from the inference API call using io.ReadAll.
func handleResponse(response *http.Response) (Response, error) {
	defer response.Body.Close()

	// Use io.ReadAll to read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Response{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var messages Response
	if err := json.Unmarshal(body, &messages); err != nil {
		return Response{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return messages, nil
}

// handleStreamedResponse processes the HTTP response from a streamed API call.
func handleStreamedResponse(response *http.Response) ([]string, error) {
	defer response.Body.Close()

	// Create a new buffered reader to read the response body line by line
	reader := bufio.NewReader(response.Body)

	var lines []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			// If we reach the end of the stream, break out of the loop
			if err.Error() == "EOF" {
				break
			}
			// For any other error, return it
			return nil, fmt.Errorf("error reading streamed response: %w", err)
		}
		lines = append(lines, line)
	}

	return lines, nil
}
