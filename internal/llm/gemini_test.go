package llm

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func Test_Prompting(t *testing.T) {
	if err := godotenv.Load(".env"); err != nil {
		t.Fatal(".env file not found")
	}

	GEMINI_API_KEY := os.Getenv("GEMINI_API_KEY")
	if GEMINI_API_KEY == "" {
		t.Fatal("GEMINI_API_KEY env var not found")
	}

	stats := PracticeOverviewParams{
		CorrectAnswers: 97,
		Mistakes: MistakeCountByDomain{
			"Advanced Math": 1,
		},
	}

	gemini := &GeminiLLM{
		apikey: GEMINI_API_KEY,
	}

	_, err := gemini.GeneratePracticeOverview(&stats)
	if err != nil {
		t.Fatalf("error isn't nil: %v", err)
	}
}
