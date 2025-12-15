package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/myselfBZ/sat-jade/internal/queries/qb_answers"
)

type QBStore struct {
	queries *qb_answers.Queries
}

func NewQBStore(db *pgxpool.Pool) *QBStore {
	queries := qb_answers.New(db)
	return &QBStore{
		queries: queries,
	}
}


func (s *QBStore) Create(ctx context.Context, a *qb_answers.CreateParams) error {
	_, err := s.queries.Delete(ctx, qb_answers.DeleteParams{
		UserID: a.UserID,
		QuestionID: a.QuestionID,
	})

	if err != nil && err != pgx.ErrNoRows {
		return err
	}

	_, err = s.queries.Create(ctx, *a)
	return err
}

func (s *QBStore) GetByUser(ctx context.Context, userId string) ([]qb_answers.QuestionBankAnswer,error) {
	validId, err :=  uuid.Parse(userId)
	if err != nil {
		return nil, err
	}
	rows, err := s.queries.GetByUser(ctx, pgtype.UUID{Bytes: validId, Valid: true})
	if err != nil {
		if err == pgx.ErrNoRows {
			return []qb_answers.QuestionBankAnswer{}, nil
		}
		return nil, err
	}

	if rows == nil {
		return []qb_answers.QuestionBankAnswer{}, nil
	}

	return rows, nil
}
