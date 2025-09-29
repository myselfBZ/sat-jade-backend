package store

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/myselfBZ/sat-jade/internal/queries/questions"
)

type Question struct {
	ModuleID      int32           `json:"module_id"`
	ID            int32           `json:"id"`
	Number        int32           `json:"number"`
	PracticeId    int32           `json:"practice_id"`
	Module        string          `json:"module"`
	Domain        string          `json:"domain"`
	Correct       string          `json:"correct"`
	Difficulty    string          `json:"difficulty"`
	Paragraph     string          `json:"paragraph"`
	Prompt        string          `json:"prompt"`
	Explanation   string          `json:"explanation"`
	AnswerChoices []*AnswerChoice `json:"choices"`
}

func NewQuestionStore(db *pgxpool.Pool) *QuestionStore {
	queries := questions.New(db)
	return &QuestionStore{
		queries: queries,
	}
}

type QuestionStore struct {
	queries *questions.Queries
}

func (s *QuestionStore) GetByModuleID(ctx context.Context, moduleID int32) ([]*Question, error) {
	questionRows, err := s.queries.GetByModuleId(ctx, moduleID)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	var quests []*Question
	for _, questionRow := range questionRows {
		quests = append(quests, &Question{
			ID:            questionRow.ID,
			Number:        questionRow.Number.Int32,
			Domain:        questionRow.Domain,
			Difficulty:    questionRow.Difficulty,
			Paragraph:     questionRow.Paragraph,
			Prompt:        questionRow.Prompt,
			Explanation:   questionRow.Explanation,
			Correct:       questionRow.Correct,
			AnswerChoices: []*AnswerChoice{},
		})
	}
	return quests, nil
}

func (s *QuestionStore) CreateWithAnswerChoices(ctx context.Context, moduleID int32, q *Question) error {
	params := questions.CreateWithAnswerChoicesParams{
		Domain:      q.Domain,
		Number:      pgtype.Int4{Int32: q.Number, Valid: true},
		SectionID:   moduleID,
		Paragraph:   q.Paragraph,
		Correct:     q.Correct,
		Prompt:      q.Prompt,
		Explanation: q.Explanation,
		Difficulty:  q.Difficulty,
	}

	if len(q.AnswerChoices) >= 4 {
		params.Label = q.AnswerChoices[0].Label
		params.Text = q.AnswerChoices[0].Text

		params.Label_2 = q.AnswerChoices[1].Label
		params.Text_2 = q.AnswerChoices[1].Text

		params.Label_3 = q.AnswerChoices[2].Label
		params.Text_3 = q.AnswerChoices[2].Text

		params.Label_4 = q.AnswerChoices[3].Label
		params.Text_4 = q.AnswerChoices[3].Text
	}
	dbQuestions, err := s.queries.CreateWithAnswerChoices(ctx, params)
	if err != nil {
		return err
	}
	q.ID = dbQuestions.ID
	return nil
}
