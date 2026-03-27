package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/myselfBZ/sat-jade/internal/queries/feedbacks"
)

type FeedbackStore struct {
	queries feedbacks.Queries
}

func NewFeedBackStore(db *pgxpool.Pool) *FeedbackStore {
	queries := feedbacks.New(db)
	return &FeedbackStore{
		queries: *queries,
	}
}

func (s *FeedbackStore) Create(ctx context.Context, userID uuid.UUID, resultID int32, feedback string) error {
	_, err := s.queries.Create(ctx, feedbacks.CreateParams{
		UserID:  pgtype.UUID{Bytes: userID, Valid: true},
		ResultID: resultID,
		Content: feedback,
	})

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23503" {
				return ErrForeignConstraint
			} else {
				return err
			}
		}
	}

	return err
}
