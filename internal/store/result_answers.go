package store

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/myselfBZ/sat-jade/internal/queries/results_answers"
)

type ResultAnswer struct {
	ID            int32  `json:"id"`
	UserAnswer    string `json:"yourAnswer"`
	ResultID      int32  `json:"result_id"`
	CorrectAnswer string `json:"correctAnswer"`
	Module        string `json:"module"`
	Status        string `json:"status"`

	// qustion part
	Number      int32  `json:"number"`
	Passage     string `json:"passage"`
	Question    string `json:"question"`
	ChoiceA     string `json:"choiceA"`
	ChoiceB     string `json:"choiceB"`
	ChoiceC     string `json:"choiceC"`
	ChoiceD     string `json:"choiceD"`
	Explanation string `json:"explanation"`
}

type ResultAnswersStore struct {
	queries *results_answers.Queries
}

func NewResultAnswersStore(db *pgxpool.Pool) *ResultAnswersStore {
	queries := results_answers.New(db)
	return &ResultAnswersStore{
		queries: queries,
	}
}

func (s *ResultAnswersStore) GetByResultID(ctx context.Context, resultID int32) ([]*ResultAnswer, error) {
	answersRow, err := s.queries.GetByResultID(ctx, resultID)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return []*ResultAnswer{}, nil
		default:
			return nil, err
		}
	}
	var answers []*ResultAnswer
	for _, a := range answersRow {
		answers = append(answers, &ResultAnswer{
			ID:            a.ID,
			UserAnswer:    a.UserAnswer.String,
			ResultID:      a.SessionID,
			CorrectAnswer: a.CorrectAnswer,
			Module:        a.Module,
			Status:        a.Status,
		})
	}
	return answers, nil
}

func (s *ResultAnswersStore) CreateMany(ctx context.Context, resultID int32, answers []*ResultAnswer) error {
	params := results_answers.CreateManyParams{
		UserAnswer:    make([]string, len(answers)),
		SessionID:     make([]int32, len(answers)),
		CorrectAnswer: make([]string, len(answers)),
		Module:        make([]string, len(answers)),
		Status:        make([]string, len(answers)),
	}

	for i, answer := range answers {
		params.UserAnswer[i] = answer.UserAnswer
		params.SessionID[i] = resultID
		params.CorrectAnswer[i] = answer.CorrectAnswer
		params.Module[i] = answer.Module
		params.Status[i] = answer.Status
	}

	_, err := s.queries.CreateMany(ctx, params)

	return err
}
