package store

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/myselfBZ/sat-jade/internal/queries/answer_choices"
)

type AnswerChoice struct {
	ID         int32  `json:"id"`
	Label      string `json:"label"`
	Text       string `json:"text"`
	QuestionID int32  `json:"question_id"`
}

func NewAnswerChoiceStore(db *pgxpool.Pool) *AnswerChoiceStore {
	queries := answer_choices.New(db)
	return &AnswerChoiceStore{
		queries: queries,
	}
}

type AnswerChoiceStore struct {
	queries *answer_choices.Queries
}

func (s *AnswerChoiceStore) GetByQuestionID(ctx context.Context, questID int32) ([]AnswerChoice, error) {
	answerChoiceRows, err := s.queries.GetByQuestionId(ctx, questID)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	var answerChoices []AnswerChoice
	for _, answerChoiceRow := range answerChoiceRows {
		answerChoices = append(answerChoices, AnswerChoice{
			QuestionID: answerChoiceRow.ID,
			ID:         answerChoiceRow.ID,
			Text:       answerChoiceRow.Text,
			Label:      answerChoiceRow.Label,
		})
	}
	return answerChoices, nil
}
func (s *AnswerChoiceStore) UpdateAnswerChoice(ctx context.Context) {

}
