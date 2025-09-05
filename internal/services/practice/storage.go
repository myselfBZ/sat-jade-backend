package practice

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/myselfBZ/sat-jade/internal/queries/qpractices"
)

type Practice struct {
	CreatedAt time.Time `json:"created_at"`
	ID        int32     `json:"id"`
	Title     string    `json:"title"`
	Modules   []Module  `json:"sections"`
}

type Module struct {
	Name      string      `json:"name"`
	Questions []*Question `json:"questions"`
}

type Question struct {
	Number        int32          `json:"number"`
	PracticeId    int32          `json:"practice_id"`
	Module        string         `json:"module"`
	Domain        string         `json:"domain"`
	Correct       string         `json:"correct"`
	Difficulty    string         `json:"difficulty"`
	Paragraph     string         `json:"paragraph"`
	Prompt        string         `json:"prompt"`
	Explanation   string         `json:"explanation"`
	AnswerChoices []AnswerChoice `json:"choices"`
}

type AnswerChoice struct {
	Label      string `json:"label"`
	Text       string `json:"text"`
	QuestionID int32  `json:"question_id"`
}

type ResultPreview struct {
	CorrectAnswers int64     `json:"correct_answers"`
	EnglishScore   int32     `json:"english_score"`
	MathScore      int32     `json:"math_score"`
	Total          int32     `json:"total_score"`
	PracticeTitle  string    `json:"practice_title"`
	ID             int32     `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
}

type TestSession struct {
	ID           int32
	UserId       string
	PracticeId   int32
	CreatedAt    time.Time
	TotalScore   int32
	MathScore    int32
	EnglishScore int32
	Answers      []*TestSessionAnswers
}

type TestSessionAnswers struct {
	ID            int32  `json:"id"`
	UserAnswer    string `json:"yourAnswer"`
	SessionID     int32  `json:"session_id"`
	CorrectAnswer string `json:"correctAnswer"`
	Module        string `json:"module"`
	Status        string `json:"status"`

	// qustion part
	Number   int32  `json:"number"`
	Passage  string `json:"passage"`
	Question string `json:"question"`
	ChoiceA  string `json:"choiceA"`
	ChoiceB  string `json:"choiceB"`
	ChoiceC  string `json:"choiceC"`
	ChoiceD  string `json:"choiceD"`
}

type Storage interface {
	GetById(ctx context.Context, id int32) (*Practice, error)
	Delete(ctx context.Context, id int32) error
	GetAllPreview(ctx context.Context) ([]*Practice, error)
	Create(ctx context.Context, title string) (int32, error)
	AddQuestion(ctx context.Context, sectionID int32, q *Question) error
	GetModuleId(ctx context.Context, practiceID int32, name string) (int32, error)
	CreateSession(ctx context.Context, session *TestSession) error
	GetResultPreviews(ctx context.Context, userId string) ([]*ResultPreview, error)
	GetSessionById(ctx context.Context, sessionId int32) (*TestSession, error)
	GetSessionAnswers(ctx context.Context, sessionId int32) ([]*TestSessionAnswers, error)
	DeleteSessionById(ctx context.Context, userId uuid.UUID, sessionId int32) error
	GetLastSession(ctx context.Context, userId uuid.UUID) (*ResultPreview, error)
}

type PostgresStorage struct {
	queries *qpractices.Queries
}

func (s *PostgresStorage) GetById(ctx context.Context, id int32) (*Practice, error) {
	practiceRow, err := s.queries.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	practice := &Practice{
		ID:        practiceRow.ID,
		Title:     practiceRow.Title,
		CreatedAt: practiceRow.CreatedAt.Time,
		Modules:   []Module{},
	}

	// Get modules for this practice
	moduleRows, err := s.queries.GetModulesByPracticeId(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, moduleRow := range moduleRows {
		module := Module{

			Name: moduleRow.Name,

			Questions: []*Question{},
		}

		// Get questions for this module
		questionRows, err := s.queries.GetQuestionsByModuleId(ctx, moduleRow.ID)
		if err != nil {
			return nil, err
		}

		for _, questionRow := range questionRows {
			question := Question{
				Number:        questionRow.Number.Int32,
				Domain:        questionRow.Domain,
				Difficulty:    questionRow.Difficulty,
				Paragraph:     questionRow.Paragraph,
				Prompt:        questionRow.Prompt,
				Explanation:   questionRow.Explanation,
				Correct:       questionRow.Correct,
				AnswerChoices: []AnswerChoice{},
			}

			// Get answer choices for this question
			choiceRows, err := s.queries.GetAnswerChoicesByQuestionId(ctx, questionRow.ID)
			if err != nil {
				return nil, err
			}

			for _, choiceRow := range choiceRows {
				choice := AnswerChoice{
					Label:      choiceRow.Label,
					Text:       choiceRow.Text,
					QuestionID: choiceRow.QuestionID,
				}
				question.AnswerChoices = append(question.AnswerChoices, choice)
			}

			module.Questions = append(module.Questions, &question)
		}

		practice.Modules = append(practice.Modules, module)
	}

	return practice, nil
}

func (s *PostgresStorage) GetAllPreview(ctx context.Context) ([]*Practice, error) {
	practicesDB, err := s.queries.GetPracticePreviews(ctx)
	if err != nil {
		return nil, err
	}
	var practices []*Practice
	for _, p := range practicesDB {
		practices = append(practices, &Practice{
			Title:     p.Title,
			CreatedAt: p.CreatedAt.Time,
			ID:        p.ID,
		})
	}
	return practices, nil
}

func (s *PostgresStorage) Create(ctx context.Context, title string) (int32, error) {
	p, err := s.queries.Create(ctx, title)
	return p.ID, err
}

func (s *PostgresStorage) AddQuestion(ctx context.Context, sectionID int32, q *Question) error {
	params := qpractices.AddQuestionParams{
		Domain:      q.Domain,
		Number:      pgtype.Int4{Int32: q.Number, Valid: true},
		SectionID:   sectionID,
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

	_, err := s.queries.AddQuestion(ctx, params)
	return err
}

func (s *PostgresStorage) GetModuleId(ctx context.Context, practiceID int32, name string) (int32, error) {
	id, err := s.queries.GetModuleID(ctx, qpractices.GetModuleIDParams{
		PracticeID: practiceID,
		Name:       name,
	})
	return id, err

}

func (s *PostgresStorage) Delete(ctx context.Context, id int32) error {
	_, err := s.queries.Delete(ctx, id)
	return err
}

func (s *PostgresStorage) CreateSession(ctx context.Context, session *TestSession) error {
	userId, _ := uuid.Parse(session.UserId)
	testSession, err := s.queries.CreateTestSession(ctx, qpractices.CreateTestSessionParams{
		PracticeID:   session.PracticeId,
		UserID:       pgtype.UUID{Bytes: userId, Valid: true},
		EnglishScore: pgtype.Int4{Int32: session.EnglishScore, Valid: true},
		TotalScore:   pgtype.Int4{Int32: session.TotalScore, Valid: true},
		MathScore:    pgtype.Int4{Int32: session.MathScore, Valid: true},
	})

	session.ID = testSession.ID

	if err != nil {
		return err
	}

	params := qpractices.CreateTestSessionAnswersParams{
		UserAnswer:    make([]string, len(session.Answers)),
		SessionID:     make([]int32, len(session.Answers)),
		CorrectAnswer: make([]string, len(session.Answers)),
		Module:        make([]string, len(session.Answers)),
		Status:        make([]string, len(session.Answers)),
	}

	for i, answer := range session.Answers {
		params.UserAnswer[i] = answer.UserAnswer
		params.SessionID[i] = testSession.ID
		params.CorrectAnswer[i] = answer.CorrectAnswer
		params.Module[i] = answer.Module
		params.Status[i] = answer.Status
	}

	_, err = s.queries.CreateTestSessionAnswers(ctx, params)
	return err
}

func (s *PostgresStorage) GetResultPreviews(ctx context.Context, userId string) ([]*ResultPreview, error) {
	id, _ := uuid.Parse(userId)
	resultsDB, err := s.queries.GetExamResultsByUserID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	var resultPreviews []*ResultPreview
	for _, r := range resultsDB {
		result := &ResultPreview{
			CorrectAnswers: r.CorrectAnswers,
			PracticeTitle:  r.PracticeTitle,
			MathScore:      r.MathScore.Int32,
			EnglishScore:   r.EnglishScore.Int32,
			Total:          r.TotalScore.Int32,
			CreatedAt:      r.CreatedAt.Time,
			ID:             r.ID,
		}
		resultPreviews = append(resultPreviews, result)
	}
	return resultPreviews, nil
}

func (s *PostgresStorage) GetSessionAnswers(ctx context.Context, sessionId int32) ([]*TestSessionAnswers, error) {
	answers, err := s.queries.GetSessionAnswers(ctx, sessionId)
	if err != nil {
		return nil, err
	}
	var sessionAnswers []*TestSessionAnswers

	for _, a := range answers {
		sessionAnswers = append(sessionAnswers, &TestSessionAnswers{
			ID:            a.ID,
			UserAnswer:    a.UserAnswer.String,
			SessionID:     a.SessionID,
			CorrectAnswer: a.CorrectAnswer,
			Module:        a.Module,
			Status:        a.Status,
		})
	}
	return sessionAnswers, nil
}

func (s *PostgresStorage) GetSessionById(ctx context.Context, sessionId int32) (*TestSession, error) {

	session, err := s.queries.GetSessionById(ctx, sessionId)
	if err != nil {
		return nil, err
	}
	return &TestSession{
		ID:         session.ID,
		UserId:     session.UserID.String(),
		PracticeId: session.PracticeID,
		CreatedAt:  session.CreatedAt.Time,
	}, nil
}

func (s *PostgresStorage) DeleteSessionById(ctx context.Context, userId uuid.UUID, sessionId int32) error {
	_, err := s.queries.DeleteSessionById(ctx, qpractices.DeleteSessionByIdParams{
		UserID: pgtype.UUID{Bytes: userId, Valid: true},
		ID:     sessionId,
	})
	return err
}

func (s *PostgresStorage) GetLastSession(ctx context.Context, userId uuid.UUID) (*ResultPreview, error) {
	result, err := s.queries.GetLastSession(ctx, pgtype.UUID{Bytes: userId, Valid: true})
	if err != nil {
		return nil, err
	}
	return &ResultPreview{
		EnglishScore: result.EnglishScore.Int32,
		MathScore:    result.MathScore.Int32,
		Total:        result.TotalScore.Int32,
		ID:           result.ID,
		CreatedAt:    result.CreatedAt.Time,
	}, nil
}
