package main

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	answereval "github.com/myselfBZ/sat-jade/internal/answer_eval"
	"github.com/myselfBZ/sat-jade/internal/queries/qb_answers"
	"github.com/myselfBZ/sat-jade/internal/store"
)

type qbAnswerPayload struct {
	Answer           string    `json:"answer"`
	UserID           string `json:"user_id"`
	QuestionID       int       `json:"question_id"`
	ResponseDuration int       `json:"response_duration"` // seconds
}

func (a *api) createQBAnswerHandler(c echo.Context) error {
	user := c.Get("user").(*store.User)
	validId, _ := uuid.Parse(user.ID)
	var payload qbAnswerPayload
	if err := c.Bind(&payload); err != nil {
		a.badRequestLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid payload")
	}
	question, err := a.storage.QuestionBank.GetById(c.Request().Context(), payload.QuestionID, user.ID)
	if err != nil {
		a.notFoundLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusNotFound, "question not found")
	}

	qbAnswer := qb_answers.CreateParams{
		QuestionID: int32(payload.QuestionID),
		UserID: pgtype.UUID{Bytes: validId, Valid: true},
		Answer: payload.Answer,
		ResponseDuration: pgtype.Int4{Int32: int32(payload.ResponseDuration), Valid: true},
	}

	if question.AnswerType == "oe" {
		correctAnswers := strings.Split(question.Correct, ",")
		correct := false
		for _, answer := range correctAnswers {
			valid, _ := answereval.EvaluateAnswer(payload.Answer, answer)
			if valid {
				correct = true
				break
			}
		}

		if correct {
			qbAnswer.Status = "correct"
		} else {
			qbAnswer.Status = "incorrect"
		}

	} else {
		if question.Correct == payload.Answer {
			qbAnswer.Status = "correct"
		} else {
			qbAnswer.Status = "incorrect"
		}
	}

	if err := a.storage.QBAnswers.Create(c.Request().Context(), &qbAnswer); err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status":qbAnswer.Status,
	})
}


func (a *api) getAnswersByUserHandler(c echo.Context) error {
	user := c.Get("user").(*store.User)
	answers, err := a.storage.QBAnswers.GetByUser(c.Request().Context(), user.ID)
	if err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, answers)
}



