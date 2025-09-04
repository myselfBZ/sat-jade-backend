package practice

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/myselfBZ/sat-jade/internal/llm"
	"github.com/myselfBZ/sat-jade/internal/queries/qpractices"
)

func New(db *pgxpool.Pool) *PracticeService {
	queries := qpractices.New(db)
	gemini := llm.NewGemini("AIzaSyDsOClmYzPr5ydj_LpKg7SNvkTzNoJqRIY")
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
