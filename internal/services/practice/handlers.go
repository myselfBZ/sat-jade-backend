package practice

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/myselfBZ/sat-jade/internal/llm"
	"github.com/myselfBZ/sat-jade/internal/services/users"
)

type practicePayload struct {
	Title string `json:"title"`
}

type testSessionPayload struct {
	ExamID  int32    `json:"exam_id"`
	Answers []string `json:"answers"`
}

type questionPayload struct {
	Number        int32          `json:"number"`
	PracticeId    int32          `json:"practice_id"`
	Difficulty    string         `json:"difficulty"`
	Module        string         `json:"section"`
	Domain        string         `json:"domain"`
	Correct       string         `json:"correct"`
	Paragraph     string         `json:"paragraph"`
	Prompt        string         `json:"prompt"`
	Explanation   string         `json:"explanation"`
	AnswerChoices []AnswerChoice `json:"choices"`
}

func (s *PracticeService) GetById(c echo.Context) error {
	id := c.Param("id")
	validId, err := strconv.Atoi(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	practice, err := s.storage.GetById(c.Request().Context(), int32(validId))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	questionNumber := 0

	for _, m := range practice.Modules {
		questionNumber += len(m.Questions)
	}

	if questionNumber != 98 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "This question isn't ready",
		})
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, practice)
}

func (s *PracticeService) Create(c echo.Context) error {
	if !isTutorOrAdmin(c) {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var payload practicePayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	id, err := s.storage.Create(c.Request().Context(), payload.Title)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"id": id,
	})
}

func (s *PracticeService) Delete(c echo.Context) error {
	if !isTutorOrAdmin(c) {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	id := c.Param("id")
	validId, err := strconv.Atoi(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	if err := s.storage.Delete(c.Request().Context(), int32(validId)); err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return c.JSON(http.StatusCreated, map[string]string{
		"Status": "Ok",
	})
}

func (s *PracticeService) AddQuestion(c echo.Context) error {
	var p questionPayload

	if err := c.Bind(&p); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	question := &Question{
		Number:        p.Number,
		PracticeId:    p.PracticeId,
		Module:        p.Module,
		Domain:        p.Domain,
		Correct:       p.Correct,
		Difficulty:    p.Difficulty,
		Paragraph:     p.Paragraph,
		Prompt:        p.Prompt,
		Explanation:   p.Explanation,
		AnswerChoices: p.AnswerChoices,
	}

	moduleId, err := s.storage.GetModuleId(c.Request().Context(), question.PracticeId, question.Module)

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "module not found")
	}

	if err := s.storage.AddQuestion(c.Request().Context(), moduleId, question); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"Status": "Ok",
	})

}

func (s *PracticeService) GetExamPreviews(c echo.Context) error {
	practices, err := s.storage.GetAllPreview(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, practices)
}

func (s *PracticeService) CreateTestSession(c echo.Context) error {
	user := c.Get("user").(*users.User)
	var p testSessionPayload
	if err := c.Bind(&p); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	practice, err := s.storage.GetById(c.Request().Context(), p.ExamID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	testSession := &TestSession{
		UserId:     user.ID,
		PracticeId: practice.ID,
		Answers:    []*TestSessionAnswers{},
	}

	answerIdx := 0

	rwCorrect := 0
	mathCorrect := 0
	for _, m := range practice.Modules {
		for _, q := range m.Questions {
			if q.Correct == p.Answers[answerIdx] {
				testSession.Answers = append(testSession.Answers, &TestSessionAnswers{
					UserAnswer:    p.Answers[answerIdx],
					CorrectAnswer: q.Correct,
					Module:        m.Name,
					Status:        "correct",
				})
				if answerIdx < 54 {
					rwCorrect++
				} else {
					mathCorrect++
				}
			} else if p.Answers[answerIdx] == "" {
				testSession.Answers = append(testSession.Answers, &TestSessionAnswers{
					UserAnswer:    p.Answers[answerIdx],
					CorrectAnswer: q.Correct,
					Module:        m.Name,
					Status:        "omitted",
				})
			} else {
				testSession.Answers = append(testSession.Answers, &TestSessionAnswers{
					UserAnswer:    p.Answers[answerIdx],
					CorrectAnswer: q.Correct,
					Module:        m.Name,
					Status:        "incorrect",
				})
			}
			answerIdx++
		}

	}
	rwScore, mathScore, totalScore := Score(rwCorrect, mathCorrect)
	testSession.MathScore = int32(mathScore)
	testSession.EnglishScore = int32(rwScore)
	testSession.TotalScore = int32(totalScore)
	if err := s.storage.CreateSession(c.Request().Context(), testSession); err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"id": testSession.ID,
	})

}

func (s *PracticeService) GetResults(c echo.Context) error {
	user := c.Get("user").(*users.User)

	resultPreviews, err := s.storage.GetResultPreviews(c.Request().Context(), user.ID)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, resultPreviews)
}

func (s *PracticeService) GetSessionById(c echo.Context) error {
	sessionId := c.Param("id")

	if sessionId == "last" {
		return s.getLastSession(c)
	}

	validSessionID, err := strconv.Atoi(sessionId)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity)
	}
	user := c.Get("user").(*users.User)

	session, err := s.storage.GetSessionById(c.Request().Context(), int32(validSessionID))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	if session.UserId != user.ID {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	answers, err := s.storage.GetSessionAnswers(c.Request().Context(), session.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	// attach question data. must do
	practice, err := s.storage.GetById(c.Request().Context(), session.PracticeId)

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	var questions []*Question

	for _, m := range practice.Modules {
		questions = append(questions, m.Questions...)
	}
	questionsIndex := 0
	for _, q := range questions {
		answers[questionsIndex].Passage = q.Paragraph
		answers[questionsIndex].Question = q.Prompt
		answers[questionsIndex].Number = q.Number
		answers[questionsIndex].ChoiceA = q.AnswerChoices[0].Text
		answers[questionsIndex].ChoiceB = q.AnswerChoices[1].Text
		answers[questionsIndex].ChoiceC = q.AnswerChoices[2].Text
		answers[questionsIndex].ChoiceD = q.AnswerChoices[3].Text
		answers[questionsIndex].Explanation = q.Explanation

		questionsIndex += 1
	}

	return c.JSON(http.StatusOK, answers)

}

func (s *PracticeService) DeleteSession(c echo.Context) error {
	user := c.Get("user").(*users.User)
	userId, _ := uuid.Parse(user.ID)
	id := c.Param("id")
	validId, err := strconv.Atoi(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity)
	}
	if err := s.storage.DeleteSessionById(c.Request().Context(), userId, int32(validId)); err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (s *PracticeService) GetSessionAIFeedback(c echo.Context) error {
	id := c.Param("id")
	validID, err := strconv.Atoi(id)
	user := c.Get("user").(*users.User)

	userId, _ := uuid.Parse(user.ID)

	session, err := s.storage.GetSessionById(c.Request().Context(), int32(validID))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity)
	}

	if session.Feedback != nil {
		var feedback AIFeedback
		if err := json.Unmarshal(session.Feedback, &feedback); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, feedback)
	}

	answers, err := s.storage.GetSessionAnswers(c.Request().Context(), int32(validID))

	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity)
	}

	practice, err := s.storage.GetById(c.Request().Context(), session.PracticeId)
	if err != nil {
		log.Println(err)
		return err
	}
	var questions []*Question

	for _, m := range practice.Modules {
		questions = append(questions, m.Questions...)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	wrongAnswers := llm.MistakeCountByDomain{}

	questionIndex := 0
	mistakeCount := 0

	for _, a := range answers {
		if a.Status != "correct" {
			q := questions[questionIndex]
			_, ok := wrongAnswers[q.Domain]
			if !ok {
				wrongAnswers[q.Domain] = 1
			} else {
				wrongAnswers[q.Domain]++
			}
			mistakeCount++
		}
		questionIndex++
	}

	feedBackParam := &llm.PracticeOverviewParams{
		CorrectAnswers: 98 - mistakeCount,
		Mistakes:       wrongAnswers,
	}

	feedback, err := s.LLM.GeneratePracticeOverview(feedBackParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	byteData, _ := json.Marshal(feedback)

	if err := s.storage.CreateFeedback(c.Request().Context(), userId, session.ID, byteData); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if err != nil {
		log.Println("LLM Error: ", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, feedback)
}

func (s *PracticeService) getLastSession(c echo.Context) error {
	user := c.Get("user").(*users.User)
	userId, _ := uuid.Parse(user.ID)
	session, err := s.storage.GetLastSession(c.Request().Context(), userId)

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, session)
}

func isTutorOrAdmin(c echo.Context) bool {
	user := c.Get("user").(*users.User)
	return (user.Role == users.ROLE_ADMIN || user.Role == users.ROLE_TUTOR)
}
