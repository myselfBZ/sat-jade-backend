package store

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/myselfBZ/sat-jade/internal/queries/results"
)

type ResultPreview struct {
	CorrectAnswers int64     `json:"correct_answers"`
	EnglishScore   int32     `json:"english_score"`
	MathScore      int32     `json:"math_score"`
	Total          int32     `json:"total_score"`
	PracticeTitle  string    `json:"practice_title"`
	ID             int32     `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
}

type Result struct {
	ID           int32
	UserId       string
	PracticeId   int32
	CreatedAt    time.Time
	TotalScore   int32
	MathScore    int32
	EnglishScore int32
	Feedback     []byte
	Answers      []*ResultAnswer
}

type ResultStore struct {
	queries *results.Queries
}

func NewResultStore(db *pgxpool.Pool) *ResultStore {
	queries := results.New(db)
	return &ResultStore{
		queries: queries,
	}
}

func (s *ResultStore) Create(ctx context.Context, result *Result) error {
	userId, _ := uuid.Parse(result.UserId)
	resultRow, err := s.queries.Create(ctx, results.CreateParams{
		PracticeID:   result.PracticeId,
		UserID:       pgtype.UUID{Bytes: userId, Valid: true},
		EnglishScore: pgtype.Int4{Int32: result.EnglishScore, Valid: true},
		TotalScore:   pgtype.Int4{Int32: result.TotalScore, Valid: true},
		MathScore:    pgtype.Int4{Int32: result.MathScore, Valid: true},
	})

	if err != nil {
		return err
	}

	result.ID = resultRow.ID
	return nil
}

func (s *ResultStore) GetByUserID(ctx context.Context, userId string) ([]*ResultPreview, error) {
	validUserID, _ := uuid.Parse(userId)

	resultRows, err := s.queries.GetByUserID(ctx, pgtype.UUID{Bytes: validUserID, Valid: true})

	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return []*ResultPreview{}, nil
		default:
			return nil, err
		}
	}

	var results []*ResultPreview
	for _, r := range resultRows {
		results = append(results, &ResultPreview{
			EnglishScore: r.EnglishScore.Int32,
			MathScore:    r.MathScore.Int32,
			Total:        r.TotalScore.Int32,
			ID:           r.ID,
			CreatedAt:    r.CreatedAt.Time,
		})
	}
	return results, nil

}

func (s *ResultStore) GetById(ctx context.Context, id int32) (*Result, error) {
	resultRow, err := s.queries.GetByID(ctx, id)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &Result{
		ID:         resultRow.ID,
		UserId:     resultRow.UserID.String(),
		PracticeId: resultRow.PracticeID,
		CreatedAt:  resultRow.CreatedAt.Time,
		Feedback:   resultRow.AiFeedback,
	}, nil
}

func (s *ResultStore) Delete(ctx context.Context, userID string, id int32) error {
	validUserID, _ := uuid.Parse(userID)

	_, err := s.queries.DeleteById(ctx, results.DeleteByIdParams{
		UserID: pgtype.UUID{Bytes: validUserID, Valid: true},
		ID:     id,
	})

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23503" {
				return ErrRecordNotFound
			}
		}

		switch err {
		case pgx.ErrNoRows:
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *ResultStore) GetAll(ctx context.Context) ([]*ResultPreview, error) {
	resultRows, err := s.queries.GetAll(ctx)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return []*ResultPreview{}, nil
		default:
			return nil, err
		}
	}
	var results []*ResultPreview
	for _, r := range resultRows {
		results = append(results, &ResultPreview{
			EnglishScore: r.EnglishScore.Int32,
			MathScore:    r.MathScore.Int32,
			Total:        r.TotalScore.Int32,
			ID:           r.ID,
			CreatedAt:    r.CreatedAt.Time,
		})
	}
	return results, nil
}
