package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/myselfBZ/sat-jade/internal/queries/result_answers"
)

type ResultAnswer struct {
	ID           int32  `json:"id"`
	UserAnswerId *int32 `json:"user_answer_id,omitempty"`
	ResultID     int32  `json:"result_id"`
	QuestionId   int32  `json:"question_id"`
	Status       string `json:"status"`

	UserAnswer   string `json:"yourAnswer"`
	// qustion part
	CorrectAnswer string `json:"correctAnswer"`
	Number        int32  `json:"number"`
	Passage       string `json:"passage"`
	Module        string `json:"module"`
	Question      string `json:"question"`
	ChoiceA       string `json:"choiceA"`
	ChoiceB       string `json:"choiceB"`
	ChoiceC       string `json:"choiceC"`
	ChoiceD       string `json:"choiceD"`
	Explanation   string `json:"explanation"`
}

type ResultAnswersStore struct {
	queries *result_answers.Queries
}

func NewResultAnswersStore(db *pgxpool.Pool) *ResultAnswersStore {
	queries := result_answers.New(db)
	return &ResultAnswersStore{
		queries: queries,
	}
}

func (s *ResultAnswersStore) GetByResultID(ctx context.Context, resultID int32) ([]ResultAnswer, error) {
	rows, err := s.queries.GetByResultID(ctx, resultID)

	if err != nil {
		// TODO map it to domain error
		return nil, err
	}

	data := make([]ResultAnswer, len(rows))

	for i := 0; i < len(rows); i++ {
		data[i].CorrectAnswer = rows[i].CorrectAnswer 
		data[i].UserAnswer = rows[i].UserAnswer
		data[i].ID = rows[i].ID
		data[i].ChoiceA = rows[i].ChoiceA.(string)
		data[i].ChoiceB = rows[i].ChoiceB.(string)
		data[i].ChoiceC = rows[i].ChoiceC.(string)
		data[i].ChoiceD = rows[i].ChoiceD.(string)
		data[i].QuestionId = rows[i].QuestionID
		data[i].Status = rows[i].Status
		data[i].Passage = rows[i].Passage
		data[i].Question = rows[i].Question
		data[i].Number = rows[i].Number.Int32
		data[i].Explanation = rows[i].Explanation
		data[i].Module = rows[i].Module
	}

	return data, nil
}

func (s *ResultAnswersStore) CreateMany(ctx context.Context, resultID int32, answers []ResultAnswer) error {
	params := result_answers.CreateManyParams{
		ResultID:   make([]int32, len(answers)),
		QuestionID: make([]int32, len(answers)),
		AnswerID:   make([]int32, len(answers)),
		Status:     make([]string, len(answers)),
	}

	for i := 0; i < len(answers); i++ {
		params.ResultID[i] = resultID
		params.QuestionID[i] = answers[i].QuestionId

		if answers[i].UserAnswerId != nil {
			params.AnswerID[i] = *answers[i].UserAnswerId
		}

		params.Status[i] = answers[i].Status
	}

	_, err := s.queries.CreateMany(ctx, params)

	return err
}
