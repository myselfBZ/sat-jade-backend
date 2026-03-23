package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/myselfBZ/sat-jade/internal/answer_eval"
	"github.com/labstack/echo/v4"
	"github.com/myselfBZ/sat-jade/internal/store"
)

type testSessionPayload struct {
	ExamID  int32    `json:"exam_id"`
	Answers []string `json:"answers"`
}

func (a *api) createResultHandler(c echo.Context) error {
	user, err := a.getUserFromContext(c)

	if err != nil {
		a.unauthorizedLog(c.Request().Method, c.Path(), err)
		return err
	}

	var p testSessionPayload
	if err := c.Bind(&p); err != nil {
		a.badRequestLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	practice, err := a.storage.Practices.GetFullTest(c.Request().Context(), p.ExamID)
	if err != nil {
		switch err {
		case store.ErrRecordNotFound:
			a.notFoundLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusNotFound)
		default:
			a.internalErrLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	result := &store.Result{
		UserId:     user.ID,
		PracticeId: practice.ID,
	}
	var answers []*store.ResultAnswer

	answerIdx := 0

	rwCorrect := 0
	mathCorrect := 0
	for _, m := range practice.Modules {
		for _, q := range m.Questions {
			if len(q.AnswerChoices) == 1{
				correct, err := answereval.EvaluateAnswer(
					p.Answers[answerIdx],
					q.AnswerChoices[0].Text,
				)
				answer := &store.ResultAnswer{}

				if correct && err == nil {
					answer.UserAnswer = p.Answers[answerIdx]
					answer.CorrectAnswer = q.Correct
					answer.Status = "correct"
					answer.Module = q.Module
				} else {
					answer.UserAnswer = p.Answers[answerIdx]
					answer.CorrectAnswer = q.Correct
					answer.Status = "incorrect"
					answer.Module = q.Module
				}

				answers = append(answers, answer)
				if answerIdx < 54 {
					rwCorrect++
				} else {
					mathCorrect++
				}
				answerIdx++
				continue
			}

			if q.Correct == p.Answers[answerIdx] {
				answers = append(answers, &store.ResultAnswer{
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
				answers = append(answers, &store.ResultAnswer{
					UserAnswer:    p.Answers[answerIdx],
					CorrectAnswer: q.Correct,
					Module:        m.Name,
					Status:        "omitted",
				})
			} else {
				answers = append(answers, &store.ResultAnswer{
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
	result.MathScore = int32(mathScore)
	result.EnglishScore = int32(rwScore)
	result.TotalScore = int32(totalScore)

	if err := a.storage.Results.Create(c.Request().Context(), result); err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if err := a.storage.ResultAnswers.CreateMany(c.Request().Context(), result.ID, answers); err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"id": result.ID,
	})

}

func (a *api) getUserResultsHandler(c echo.Context) error {

	user, err := a.getUserFromContext(c)
	if err != nil {
		a.unauthorizedLog(c.Request().Method, c.Path(), err)
		return err
	}

	resultPreviews, err := a.storage.Results.GetByUserID(c.Request().Context(), user.ID)
	if err != nil {
		switch err {
		case store.ErrRecordNotFound:
			a.notFoundLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusNotFound)
		default:
			a.internalErrLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	return c.JSON(http.StatusOK, resultPreviews)
}

func (a *api) getAllResultsByUserHandler(c echo.Context) error {
	userId := c.Param("userId")
	results, err := a.storage.Results.GetByUserID(c.Request().Context(), userId)

	if err != nil {
		a.notFoundLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, results)
}

func (a *api) getAllResultsHandler(c echo.Context) error {
	results, err := a.storage.Results.GetAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, results)
}

func (a *api) getResultByIDHandler(c echo.Context) error {
	resultID := c.Param("id")
	validResultID, err := strconv.Atoi(resultID)
	if err != nil {
		a.badRequestLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity)
	}
	user, err := a.getUserFromContext(c)

	if err != nil {
		a.unauthorizedLog(c.Request().Method, c.Path(), err)
		return err
	}

	session, err := a.storage.Results.GetById(c.Request().Context(), int32(validResultID))
	if err != nil {
		switch err {
		case store.ErrRecordNotFound:
			a.notFoundLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusNotFound)
		default:
			a.internalErrLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	if session.UserId != user.ID {
		a.unauthorizedLog(c.Request().Method, c.Path(), errors.New("ownership didnt match"))
		return echo.NewHTTPError(http.StatusUnauthorized, "you are not allowed to see this result")
	}

	answers, err := a.storage.ResultAnswers.GetByResultID(c.Request().Context(), session.ID)
	if err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	// attach question data. must do
	practice, err := a.storage.Practices.GetFullTest(c.Request().Context(), session.PracticeId)

	if err != nil {
		switch err {
		case store.ErrRecordNotFound:
			a.notFoundLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusNotFound)
		default:
			a.internalErrLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	var questions []*store.Question

	for _, m := range practice.Modules {
		questions = append(questions, m.Questions...)
	}
	questionsIndex := 0
	for _, q := range questions {
		answers[questionsIndex].Passage = q.Paragraph
		answers[questionsIndex].Question = q.Prompt
		answers[questionsIndex].Number = q.Number
		answers[questionsIndex].ChoiceA = q.AnswerChoices[0].Text
		if len(q.AnswerChoices) == 4 {
			answers[questionsIndex].ChoiceB = q.AnswerChoices[1].Text
			answers[questionsIndex].ChoiceC = q.AnswerChoices[2].Text
			answers[questionsIndex].ChoiceD = q.AnswerChoices[3].Text
		}
		answers[questionsIndex].Explanation = q.Explanation

		questionsIndex += 1
	}

	return c.JSON(http.StatusOK, answers)

}

func (a *api) deleteResultByIDHandler(c echo.Context) error {
	user, err := a.getUserFromContext(c)
	if err != nil {
		return err
	}
	id := c.Param("id")
	validId, err := strconv.Atoi(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity)
	}
	if err := a.storage.Results.Delete(c.Request().Context(), user.ID, int32(validId)); err != nil {
		switch err {
		case store.ErrRecordNotFound:
			a.notFoundLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusNotFound)
		default:
			a.internalErrLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
