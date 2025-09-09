package llm

// ChatCompletionResponse represents the full API response
type groqPayload struct {
	ID                string           `json:"id"`
	Object            string           `json:"object"`
	Created           int64            `json:"created"`
	Model             string           `json:"model"`
	Choices           []ResponseChoice `json:"choices"`
	Usage             Usage            `json:"usage"`
	UsageBreakdown    *UsageBreakdown  `json:"usage_breakdown"`
	SystemFingerprint string           `json:"system_fingerprint"`
	XGroq             *XGroq           `json:"x_groq,omitempty"`
	ServiceTier       string           `json:"service_tier"`
}

// Choice represents individual response choices
type ResponseChoice struct {
	Index        int             `json:"index"`
	Message      ResponseMessage `json:"message"`
	LogProbs     *LogProb        `json:"logprobs"`
	FinishReason string          `json:"finish_reason"`
}

// Message represents the assistant's message
type ResponseMessage struct {
	Role      string  `json:"role"`
	Content   string  `json:"content"`
	Reasoning *string `json:"reasoning,omitempty"`
}

// LogProb represents log probability information (currently null in the payload)
type LogProb struct {
	// This would contain log probability data if present
	// Structure depends on the specific implementation
}

// Usage represents token usage statistics
type Usage struct {
	QueueTime        float64 `json:"queue_time"`
	PromptTokens     int     `json:"prompt_tokens"`
	PromptTime       float64 `json:"prompt_time"`
	CompletionTokens int     `json:"completion_tokens"`
	CompletionTime   float64 `json:"completion_time"`
	TotalTokens      int     `json:"total_tokens"`
	TotalTime        float64 `json:"total_time"`
}

// UsageBreakdown represents detailed usage breakdown (currently null in the payload)
type UsageBreakdown struct {
	// This would contain breakdown data if present
	// Structure depends on the specific implementation
}

// XGroq represents Groq-specific metadata
type XGroq struct {
	ID string `json:"id"`
}
