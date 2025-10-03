package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/myselfBZ/sat-jade/internal/llm"
	"github.com/myselfBZ/sat-jade/internal/store"
)

type AIFeedback struct {
	Overview     string   `json:"overview"`
	Suggesttions []string `json:"suggestions"`
	Motivation   string   `json:"motivation"`
}

func (a *api) getOrCreateAIFeedbackHandler(c echo.Context) error {
	id := c.Param("id")
	validID, err := strconv.Atoi(id)

	if err != nil {
		a.badRequestLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	user, err := a.getUserFromContext(c)

	if err != nil {
		a.unauthorizedLog(c.Request().Method, c.Path(), err)
		return err
	}

	userId, _ := uuid.Parse(user.ID)

	result, err := a.storage.Results.GetById(c.Request().Context(), int32(validID))

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

	if result.Feedback != nil {
		var feedback AIFeedback
		if err := json.Unmarshal(result.Feedback, &feedback); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, feedback)
	}

	answers, err := a.storage.ResultAnswers.GetByResultID(c.Request().Context(), int32(validID))

	if err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	practice, err := a.storage.Practices.GetById(c.Request().Context(), result.PracticeId)

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

	if err := a.constructPracticeTest(c.Request().Context(), practice); err != nil {
		return err
	}

	var questions []*store.Question

	for _, m := range practice.Modules {
		questions = append(questions, m.Questions...)
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

	feedback, err := a.llm.GeneratePracticeOverview(feedBackParam)
	if err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	byteData, _ := json.Marshal(feedback)

	if err := a.storage.Feedback.Create(c.Request().Context(), userId, result.ID, byteData); err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, feedback)
}
