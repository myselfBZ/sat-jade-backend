package store

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/myselfBZ/sat-jade/internal/queries/qb_answers"
)

var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrDuplicateEmail    = errors.New("this email is already taken")
	ErrForeignConstraint = errors.New("foreign key constraint violated")
	ErrConflict 		 = errors.New("conflict request")
	ErrInvalidUUID = errors.New("invalid uuid")
)

type Storage struct {
	Users         UserRepository
	Practices     PracticeRepository
	Modules       ModuleRepository
	Questions     QuestionRepository
	AnswerChoices AnswerChoiceRepository
	Results       ResultRepository
	ResultAnswers ResultAnswerRepository
	Feedback      FeedbackRepository
	QuestionBank  QuestionBankRepository
	QBAnswers     QBAnswerRepository
}

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetMany(ctx context.Context) ([]User, error)
	Delete(ctx context.Context, id string) error
}

type PracticeRepository interface {
	Delete(ctx context.Context, id int32) error
	GetAllPreview(ctx context.Context) ([]PracticePreview, error)
	Create(ctx context.Context, title string) (int32, error)
	GetFullTest(ctx context.Context, id int32) (*Practice, error)
	GetCorrectAnswersWithAnswerChoices(ctx context.Context, id int32) ([]CorrectAnswerWithAnswerChoices, error) 
}

type ModuleRepository interface {
	GetAllByPracticeID(ctx context.Context, practiceID int32) ([]Module, error)
	GetByID(ctx context.Context, id int32) (*Module, error)
	GetByNameAndPracticeID(ctx context.Context, name string, practiceID int32) (*Module, error)
}

type QuestionRepository interface {
	CreateWithAnswerChoices(ctx context.Context, moduleID int32, q *Question) error
}

type AnswerChoiceRepository interface {
	GetByQuestionID(ctx context.Context, questID int32) ([]AnswerChoice, error)
	UpdateAnswerChoice(ctx context.Context)
}

type ResultRepository interface {
	Create(ctx context.Context, result *Result) error
	GetByUserID(ctx context.Context, userId string) ([]ResultPreview, error)
	GetById(ctx context.Context, sessionId int32) (*Result, error)
	Delete(ctx context.Context, userID string, id int32) error
	GetAll(ctx context.Context) ([]ResultPreview, error)
}

type ResultAnswerRepository interface {
	CreateMany(ctx context.Context, resultID int32, answers []ResultAnswer) error
	GetByResultID(ctx context.Context, resultID int32) ([]ResultAnswer, error)
}

type FeedbackRepository interface {
	Create(ctx context.Context, userID uuid.UUID, resultID int32, feedback string) error
}

type QuestionBankRepository interface {
	GetIdBySkill(ctx context.Context, skill string) ([]int32, error)
	Create(ctx context.Context, q *SQBQuestion) error
	GetCollectionDetail(ctx context.Context) ([]CollectionDetail, error)
	GetById(ctx context.Context, id int, userId string) (*SQBQuestion, error)
}

type QBAnswerRepository interface {
	Create(ctx context.Context, a *qb_answers.CreateParams) error
	GetByUser(ctx context.Context, userId string) ([]qb_answers.QuestionBankAnswer, error)
}

func New(db *pgxpool.Pool) *Storage {
	return &Storage{
		Users:               NewUserStore(db),
		Practices:           NewPracticeStore(db),
		Modules:             NewModuleStore(db),
		Questions:           NewQuestionStore(db),
		AnswerChoices: NewAnswerChoiceStore(db),
		Results:             NewResultStore(db),
		ResultAnswers:       NewResultAnswersStore(db),
		Feedback:            NewFeedBackStore(db),
		QuestionBank:        NewQuestionBank(db),
		QBAnswers:           NewQBStore(db),
	}
}
