package practice

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/myselfBZ/sat-jade/internal/queries/qpractices"
)

func New(db *pgxpool.Pool) *PracticeService {
	queries := qpractices.New(db)
	return &PracticeService{
		storage: &PostgresStorage{
			queries: queries,
		},
	}
}

type PracticeService struct {
	storage Storage
}
