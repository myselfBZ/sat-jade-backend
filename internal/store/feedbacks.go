package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/myselfBZ/sat-jade/internal/queries/feedbacks"
)

type FeedbackStore struct {
	queries feedbacks.Queries
}

type Feedback struct {
	UserId    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Header    string    `json:"header"`
	Body      string    `json:"body"`
	Footer    string    `json:"footer"`
}

func NewFeedBackStore(db *pgxpool.Pool) *FeedbackStore {
	queries := feedbacks.New(db)
	return &FeedbackStore{
		queries: *queries,
	}
}

func (s *FeedbackStore) Get(ctx context.Context, resultId int32) (*Feedback, error) {
	f, err := s.queries.Get(ctx, resultId)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	feedback := &Feedback{
		Header:    f.Header,
		Body:      f.Body,
		Footer:    f.Footer,
		CreatedAt: f.CreatedAt.Time,
		UserId:    f.UserID.String(),
	}

	return feedback, nil
}


func (s *FeedbackStore) Create(ctx context.Context, params feedbacks.CreateParams) (*Feedback, error) {
	f, err := s.queries.Create(ctx, params)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23503" {
				return nil, ErrForeignConstraint
			} else {
				return nil, err
			}
		}
	}

	feedback := &Feedback{
		Header:    f.Header,
		Body:      f.Body,
		Footer:    f.Footer,
		CreatedAt: f.CreatedAt.Time,
		UserId:    f.UserID.String(),
	}

	return feedback, err
}
