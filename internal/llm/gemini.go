package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

func NewGemini(apiKey string) (LLM, error) {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		return nil, err
	}

	return &GeminiLLM{
		apikey: apiKey,
		client: client,
	}, nil
}

type GeminiLLM struct {
	apikey string
	client *genai.Client
}

func (l *GeminiLLM) GeneratePracticeOverview(params *PracticeOverviewParams) (*PracticeOverview, error) {
	mistakesByDomain := ""
	for domain, mistakeCount := range params.Mistakes {
		mistakesByDomain += fmt.Sprintf(domain+": %d\n", mistakeCount)
	}

	prompt := fmt.Sprintf(THE_PRACTICE_PROMPT, 98-params.CorrectAnswers, params.CorrectAnswers, mistakesByDomain)

	result, err := l.client.Models.GenerateContent(
		context.TODO(),
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return nil, err
	}

	var payload PracticeOverview

	text := result.Text()
	clean := strings.TrimPrefix(text, "```json")
	clean = strings.TrimSuffix(clean, "```")

	if err := json.Unmarshal([]byte(clean), &payload); err != nil {
		return nil, err
	}

	return &payload, nil
}
