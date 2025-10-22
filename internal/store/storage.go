package store

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrDuplicateEmail    = errors.New("this email is already taken")
	ErrForeignConstraint = errors.New("foreign key constraint violated")
)

type Storage struct {
	Users interface {
		Create(ctx context.Context, u *User) error
		GetByEmail(ctx context.Context, email string) (*User, error)
		GetByID(ctx context.Context, id string) (*User, error)
		GetMany(ctx context.Context) ([]*User, error)
	}

	Practices interface {
		GetById(ctx context.Context, id int32) (*Practice, error)
		Delete(ctx context.Context, id int32) error
		// Fetches the preview of practice tests
		// If there is none, empty slice will be returned
		// []*PracticePrview

		GetAllPreview(ctx context.Context) ([]*PracticePreview, error)
		Create(ctx context.Context, title string) (int32, error)
	}

	Modules interface {
		GetAllByPracticeID(ctx context.Context, practiceID int32) ([]*Module, error)
		GetByID(ctx context.Context, id int32) (*Module, error)
		GetByNameAndPracticeID(ctx context.Context, name string, practiceID int32) (*Module, error)
	}

	Questions interface {
		CreateWithAnswerChoices(ctx context.Context, moduleID int32, q *Question) error
		GetByModuleID(ctx context.Context, moduleID int32) ([]*Question, error)
		GetByModuleWithChoices(ctx context.Context, moduleID int32) ([]*Question, error) 
	}

	AnswerChoiceStorage interface {
		GetByQuestionID(ctx context.Context, questID int32) ([]*AnswerChoice, error)
		// Implementation TBD
		UpdateAnswerChoice(ctx context.Context)
	}

	Results interface {
		Create(ctx context.Context, result *Result) error
		GetByUserID(ctx context.Context, userId string) ([]*ResultPreview, error)
		GetById(ctx context.Context, sessionId int32) (*Result, error)
		Delete(ctx context.Context, userID string, id int32) error
		GetAll(ctx context.Context) ([]*ResultPreview, error)
	}

	ResultAnswers interface {
		CreateMany(ctx context.Context, resultID int32, answers []*ResultAnswer) error
		GetByResultID(ctx context.Context, resultID int32) ([]*ResultAnswer, error)
	}

	Feedback interface {
		Create(ctx context.Context, userID uuid.UUID, resultID int32, feedback []byte) error
	}

	DailyQuestions interface {
		Create(ctx context.Context, q *DailyQuestion) error
		GetLatest(ctx context.Context) ([]*DailyQuestion, error)
		GetAll(ctx context.Context) ([]*DailyQuestion, error)
	}
}

func New(db *pgxpool.Pool) *Storage {
	return &Storage{
		Users:               NewUserStore(db),
		Practices:           NewPracticeStore(db),
		Modules:             NewModuleStore(db),
		Questions:           NewQuestionStore(db),
		AnswerChoiceStorage: NewAnswerChoiceStore(db),
		Results:             NewResultStore(db),
		ResultAnswers:       NewResultAnswersStore(db),
		Feedback:            NewFeedBackStore(db),
		DailyQuestions:      NewDailyQuestionStore(db),
	}
}
