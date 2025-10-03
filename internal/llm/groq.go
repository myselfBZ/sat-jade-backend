package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Groq struct {
	apikey string
}

// Request payload structs
type GroqRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Response payload struct (you already have this)
type GroqPayload struct {
	Choices []Choice `json:"choices"`
	// ... other fields from your existing struct
}

type Choice struct {
	Message Message `json:"message"`
	// ... other fields
}

func (g *Groq) GeneratePracticeOverview(params *PracticeOverviewParams) (*PracticeOverview, error) {
	mistakesByDomain := ""
	for domain, mistakeCount := range params.Mistakes {
		mistakesByDomain += fmt.Sprintf(domain+": %d ", mistakeCount)
	}

	data, err := g.requestLLM(fmt.Sprintf(THE_PRACTICE_PROMPT, 98-params.CorrectAnswers, params.CorrectAnswers, mistakesByDomain))
	if err != nil {
		return nil, err
	}

	var overview PracticeOverview
	if err := json.Unmarshal(data, &overview); err != nil {
		log.Println("errro unmarshaling the data:", string(data))
		return nil, err
	}

	return &overview, nil
}

func (g *Groq) requestLLM(prompt string) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: false,
		},
	}

	// Create request payload using proper structs
	requestPayload := GroqRequest{
		Model: "openai/gpt-oss-20b",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// Marshal to JSON properly
	jsonData, err := json.Marshal(requestPayload)
	if err != nil {
		log.Println("error marshaling request payload: ", err)
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.groq.com/openai/v1/chat/completions",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Println("error creating a request for LLM: ", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.apikey)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("error requesting llm: ", err)
		return nil, err
	}
	defer resp.Body.Close() // Don't forget to close the response body

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var payload GroqPayload
	if err := json.Unmarshal(respBody, &payload); err != nil {
		log.Println("Couldn't decode response: ", err)
		return nil, err
	}

	if len(payload.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned from API")
	}

	text := payload.Choices[0].Message.Content

	// Clean up markdown code blocks if present
	clean := strings.TrimPrefix(text, "```json")
	clean = strings.TrimSuffix(clean, "```")
	clean = strings.TrimSpace(clean)

	return []byte(clean), nil
}
