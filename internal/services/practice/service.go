package practice

import (
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/myselfBZ/sat-jade/internal/llm"
	"github.com/myselfBZ/sat-jade/internal/queries/qpractices"
)

func New(db *pgxpool.Pool) *PracticeService {
	queries := qpractices.New(db)
	//groq := llm.NewGroq("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")
	gemini, err := llm.NewGemini("AIzaSyDsOClmYzPr5ydj_LpKg7SNvkTzNoJqRIY")
	if err != nil {
		log.Fatal("couldn't create an LLM client: ", err)
	}
	return &PracticeService{
		storage: &PostgresStorage{
			queries: queries,
		},
		LLM: gemini,
	}
}

type PracticeService struct {
	storage Storage
	LLM     llm.LLM
}
