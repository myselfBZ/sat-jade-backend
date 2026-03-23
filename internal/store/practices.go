package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/myselfBZ/sat-jade/internal/queries/qpractices"
)

type Practice struct {
	CreatedAt time.Time `json:"created_at"`
	ID        int32     `json:"id"`
	Title     string    `json:"title"`
	Modules   []*Module `json:"sections"`
}

type PracticePreview struct {
	CreatedAt time.Time `json:"created_at"`
	ID        int32     `json:"id"`
	Title     string    `json:"title"`
}

type PracticeStore struct {
	queries *qpractices.Queries
}

func NewPracticeStore(db *pgxpool.Pool) *PracticeStore {
	queries := qpractices.New(db)
	return &PracticeStore{
		queries: queries,
	}
}


func (s *PracticeStore) GetAllPreview(ctx context.Context) ([]*PracticePreview, error) {
	practicesDB, err := s.queries.GetPracticePreviews(ctx)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			// empty slice instead of error
			return []*PracticePreview{}, nil
		default:
			return nil, err
		}
	}
	var practices []*PracticePreview
	for _, p := range practicesDB {
		practices = append(practices, &PracticePreview{
			Title:     p.Title,
			CreatedAt: p.CreatedAt.Time,
			ID:        p.ID,
		})
	}
	return practices, nil
}

func (s *PracticeStore) Create(ctx context.Context, title string) (int32, error) {
	p, err := s.queries.Create(ctx, title)
	return p.ID, err
}

func (s *PracticeStore) Delete(ctx context.Context, id int32) error {
	_, err := s.queries.Delete(ctx, id)
	return err
}

func (s *PracticeStore) GetFullTest(ctx context.Context, id int32) (*Practice, error) {
	row, err := s.queries.GetFullPracticeTest(ctx, id)

	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	var practice Practice

	practice.CreatedAt = row.CreatedAt.Time
	practice.Title = row.Title
	practice.ID = row.ID

	if err := json.Unmarshal(row.Sections, &practice.Modules); err != nil {
		return nil, err
	}

	return &practice, nil
}
