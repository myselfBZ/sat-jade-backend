package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/myselfBZ/sat-jade/internal/queries/daily_questions"
)

type DailyQuestionStore struct {
	queries *daily_questions.Queries
}

type DailyQuestion struct {
	ID          int32     `json:"id"`
	Domain      string    `json:"domain"`
	Correct     string    `json:"correct"`
	Difficulty  string    `json:"difficulty"`
	Paragraph   string    `json:"paragraph"`
	Prompt      string    `json:"prompt"`
	Explanation string    `json:"explanation"`
	CreatedAt   time.Time `json:"created_at"`
	ChoiceA     string    `json:"choice_a"`
	ChoiceB     string    `json:"choice_b"`
	ChoiceC     string    `json:"choice_c"`
	ChoiceD     string    `json:"choice_d"`
}

func NewDailyQuestionStore(db *pgxpool.Pool) *DailyQuestionStore {
	queries := daily_questions.New(db)
	return &DailyQuestionStore{
		queries,
	}
}

func (s *DailyQuestionStore) GetAll(ctx context.Context) ([]*DailyQuestion, error) {
	questionRows, err := s.queries.GetAll(ctx)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return []*DailyQuestion{}, nil
		default:
			return nil, err
		}
	}

	var questions []*DailyQuestion

	for _, q := range questionRows {
		questions = append(questions, &DailyQuestion{
			ID:          q.ID,
			Domain:      q.Domain,
			Paragraph:   q.Paragraph,
			Correct:     q.Correct,
			Prompt:      q.Prompt,
			Explanation: q.Explanation,
			Difficulty:  q.Difficulty,

			ChoiceA: q.ChoiceA.String,
			ChoiceB: q.ChoiceB.String,
			ChoiceC: q.ChoiceC.String,
			ChoiceD: q.ChoiceD.String,
		})
	}
	return questions, nil
}

func (s *DailyQuestionStore) Create(ctx context.Context, q *DailyQuestion) error {
	qRow, err := s.queries.Create(ctx, daily_questions.CreateParams{
		Domain:      q.Domain,
		Paragraph:   q.Paragraph,
		Correct:     q.Correct,
		Prompt:      q.Prompt,
		Explanation: q.Explanation,
		Difficulty:  q.Difficulty,
		CreatedAt:   pgtype.Timestamp{Time: q.CreatedAt, Valid: true},
		ChoiceA:     pgtype.Text{String: q.ChoiceA, Valid: true},
		ChoiceB:     pgtype.Text{String: q.ChoiceB, Valid: true},
		ChoiceC:     pgtype.Text{String: q.ChoiceC, Valid: true},
		ChoiceD:     pgtype.Text{String: q.ChoiceD, Valid: true},
	})

	if err != nil {
		return err
	}

	q.ID = qRow.ID
	return nil
}

func (s *DailyQuestionStore) GetLatest(ctx context.Context) ([]*DailyQuestion, error) {
	questionsRow, err := s.queries.GetLatest(ctx)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return []*DailyQuestion{}, nil
		default:
			return nil, err
		}
	}

	var questions []*DailyQuestion

	for _, q := range questionsRow {
		questions = append(questions, &DailyQuestion{
			ID:          q.ID,
			Domain:      q.Domain,
			Paragraph:   q.Paragraph,
			Correct:     q.Correct,
			Prompt:      q.Prompt,
			Explanation: q.Explanation,
			Difficulty:  q.Difficulty,

			ChoiceA: q.ChoiceA.String,
			ChoiceB: q.ChoiceB.String,
			ChoiceC: q.ChoiceC.String,
			ChoiceD: q.ChoiceD.String,
		})
	}

	return questions, nil
}
