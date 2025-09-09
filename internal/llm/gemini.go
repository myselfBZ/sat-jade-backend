package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

func NewGemini(apiKey string) (*GeminiLLM, error) {
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

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
			Role string `json:"role"`
		} `json:"content"`
		FinishReason string  `json:"finishReason"`
		AvgLogprobs  float64 `json:"avgLogprobs"`
	} `json:"candidates"`
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

/*func (l *GeminiLLM) GeneratePracticeOverview(params *PracticeOverviewParams) (*PracticeOverview, error) {
	mistakesByDomain := ""

	for domain, mistakeCount := range params.Mistakes {
		mistakesByDomain += fmt.Sprintf(domain+": %d\n", mistakeCount)
	}

	data, err := l.requestLLM(fmt.Sprintf(THE_PRACTICE_PROMPT, 98-params.CorrectAnswers, params.CorrectAnswers, mistakesByDomain))
	log.Println("I have the data!!!!", err)
	if err != nil {
		return nil, err
	}

	var overview PracticeOverview
	if err := json.Unmarshal(data, &overview); err != nil {
		return nil, err
	}
	return &overview, nil
}

func (l *GeminiLLM) requestLLM(prompt string) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: false,
		},
	}

	data := fmt.Sprintf(`{
    "contents": [
      {
        "parts": [
          {
            "text": "%s"
          }
        ]
      }
    ]
  }`, prompt)

	req, err := http.NewRequest(
		"POST",
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent",
		bytes.NewBuffer([]byte(data)),
	)

	if err != nil {
		log.Println("error creating a request for LLM: ", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-goog-api-key", l.apikey)

	resp, err := client.Do(req)

	if err != nil {
		log.Println("error requesting llm: ", err)
		return nil, err
	}

	var payload geminiResponse

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		log.Println("Couldn't decode: ", err)
		return nil, err
	}

	if !(len(payload.Candidates) > 0) {
		return nil, errors.New("No Ai response")
	}

	if !(len(payload.Candidates[0].Content.Parts) > 0) {
		return nil, errors.New("No Ai response")
	}

	text := payload.Candidates[0].Content.Parts[0].Text
	clean := strings.TrimPrefix(text, "```json")
	clean = strings.TrimSuffix(clean, "```")
	log.Println("text is here: ", text)
	return []byte(clean), nil
}
*/
