package llm

import (
	"log"
	"testing"
)

func Test_Prompting(t *testing.T) {
	stats := PracticeOverviewParams{
		CorrectAnswers: 97,
		Mistakes: MistakeCountByDomain{
			"Advanced Math": 1,
		},
	}

	gemini := &GeminiLLM{
		apikey: "AIzaSyDsOClmYzPr5ydj_LpKg7SNvkTzNoJqRIY",
	}

	overview, err := gemini.GeneratePracticeOverview(&stats)
	if err != nil {
		t.Fatalf("error isn't nil: %v", err)
	}
	log.Println(overview.Overview)
	log.Println(overview.Motivation)
}
