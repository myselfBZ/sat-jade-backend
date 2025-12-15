package store

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/myselfBZ/sat-jade/internal/queries/qb_answers"
	"github.com/myselfBZ/sat-jade/internal/queries/question_bank"
)

type SQBQuestion struct {
	ID          int32  `json:"id"`
	Domain      string `json:"domain"`
	Paragraph   string `json:"paragraph"`
	Skill       string `json:"skill"`
	QuestionID  string `json:"questionID"`
	Correct     string `json:"correct"`
	Prompt      string `json:"prompt"`
	Explanation string `json:"explanation"`
	Difficulty  string `json:"difficulty"`
	AnswerType  string `json:"answerType"`
	Active      bool   `json:"active"`
	ChoiceA     string `json:"choiceA"`
	ChoiceB     string `json:"choiceB"`
	ChoiceC     string `json:"choiceC"`
	ChoiceD     string `json:"choiceD"`

	Response *qb_answers.QuestionBankAnswer `json:"response,omitempty"`
}

type CollectionDetail struct {
	Quantity int      `json:"quantity"`
	Domain   string   `json:"domain"`
	Skills   []*Skill `json:"skills"`
}

type Skill struct {
	Quantity int    `json:"quantity"`
	Text     string `json:"text"`
}

func NewQuestionBank(db *pgxpool.Pool) *QuestionBank {
	queries := question_bank.New(db)
	return &QuestionBank{queries}
}

type QuestionBank struct {
	queries *question_bank.Queries
}

func (s *QuestionBank) GetById(ctx context.Context, id int, userId string) (*SQBQuestion, error) {
	validId, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("Invalid user id")
	}
	row, err := s.queries.GetById(ctx, question_bank.GetByIdParams{
		ID:     int32(id),
		UserID: pgtype.UUID{Bytes: validId, Valid: true},
	})

	if err != nil {
		return nil, err
	}
	sqb := &SQBQuestion{
		ID:          row.ID,
		Active:      true,
		Domain:      row.Domain,
		Skill:       row.Skill.String,
		QuestionID:  row.QuestionID.String,
		Difficulty:  row.Difficulty,
		Explanation: row.Explanation,
		Prompt:      row.Prompt,
		Paragraph:   row.Paragraph,
		Correct:     row.Correct,
		AnswerType:  row.AnswerType.String,
		ChoiceA:     row.ChoiceA.String,
		ChoiceB:     row.ChoiceB.String,
		ChoiceC:     row.ChoiceC.String,
		ChoiceD:     row.ChoiceD.String,
	}

	if row.Status.String != "" {
		sqb.Response = &qb_answers.QuestionBankAnswer{
			UserID:           pgtype.UUID{Bytes: validId, Valid: true},
			Answer:           row.Answer.String,
			QuestionID:       row.QuestionID_2.Int32,
			Status:           row.Status.String,
			ResponseDuration: row.ResponseDuration,
			CreatedAt:        row.CreatedAt,
		}
	}

	return sqb, nil
}

func (s *QuestionBank) GetIdBySkill(ctx context.Context, skill string) ([]int32, error) {
	rows, err := s.queries.GetIdBySkill(ctx, pgtype.Text{String: skill, Valid: true})
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *QuestionBank) Create(ctx context.Context, q *SQBQuestion) error {
	_, err := s.queries.Create(ctx, question_bank.CreateParams{
		Domain:      q.Domain,
		Paragraph:   q.Paragraph,
		Prompt:      q.Prompt,
		Skill:       pgtype.Text{String: q.Skill, Valid: true},
		ChoiceA:     pgtype.Text{String: q.ChoiceA, Valid: true},
		ChoiceB:     pgtype.Text{String: q.ChoiceB, Valid: true},
		ChoiceC:     pgtype.Text{String: q.ChoiceC, Valid: true},
		ChoiceD:     pgtype.Text{String: q.ChoiceD, Valid: true},
		Correct:     strings.TrimSpace(q.Correct),
		AnswerType:  pgtype.Text{String: q.AnswerType, Valid: true},
		QuestionID:  pgtype.Text{String: q.QuestionID, Valid: true},
		Active:      pgtype.Bool{Bool: q.Active, Valid: true},
		Difficulty:  q.Difficulty,
		Explanation: q.Explanation,
	})
	return err
}

func (s *QuestionBank) GetCollectionDetail(ctx context.Context) ([]*CollectionDetail, error) {
	rows, err := s.queries.GetCollectionDetails(ctx)

	if err != nil {
		return nil, err
	}

	collectionMap := make(map[string]*CollectionDetail)
	for _, row := range rows {
		details, ok := collectionMap[row.Domain]
		if !ok {
			collectionMap[row.Domain] = &CollectionDetail{
				Quantity: int(row.Count),
				Domain:   row.Domain,
				Skills: []*Skill{
					{
						Text:     row.Skill.String,
						Quantity: int(row.Count),
					},
				},
			}
			continue
		}
		details.Quantity += int(row.Count)
		details.Skills = append(details.Skills, &Skill{Text: row.Skill.String, Quantity: int(row.Count)})
	}

	result := make([]*CollectionDetail, len(collectionMap))

	i := 0

	for _, v := range collectionMap {
		result[i] = v
		i++
	}

	return result, nil
}
