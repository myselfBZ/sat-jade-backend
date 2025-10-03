package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/myselfBZ/sat-jade/internal/store"
)

type dailyQuestionPayload struct {
	Domain      string    `json:"domain"`
	Correct     string    `json:"correct"`
	Difficulty  string    `json:"difficulty"`
	Paragraph   string    `json:"paragraph"`
	Prompt      string    `json:"prompt"`
	Explanation string    `json:"explanation"`
	CreatedAt   time.Time `json:"created_at"`
	ChoiceA     string    `json:"choice_a"`
	ChoiceB     string    `json:"choice_b"`
	ChoiceC     string    `json:"choice_c"`
	ChoiceD     string    `json:"choice_d"`
}

func (a *api) createDailyQuestionHandler(c echo.Context) error {
	var payload dailyQuestionPayload
	if err := c.Bind(&payload); err != nil {
		a.badRequestLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	question := &store.DailyQuestion{
		Domain:      payload.Domain,
		Paragraph:   payload.Paragraph,
		Correct:     payload.Correct,
		Prompt:      payload.Prompt,
		Explanation: payload.Explanation,
		Difficulty:  payload.Difficulty,
		CreatedAt:   payload.CreatedAt,
		ChoiceA:     payload.ChoiceA,
		ChoiceB:     payload.ChoiceB,
		ChoiceC:     payload.ChoiceC,
		ChoiceD:     payload.ChoiceD,
	}

	if err := a.storage.DailyQuestions.Create(c.Request().Context(), question); err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"status": "ok",
	})
}

func (a *api) getAllDailyQuestionsHandler(c echo.Context) error {
	questios, err := a.storage.DailyQuestions.GetAll(c.Request().Context())
	if err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, questios)
}

func (a *api) getDailyQuestionsHandler(c echo.Context) error {
	questions, err := a.storage.DailyQuestions.GetLatest(c.Request().Context())

	if err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, questions)
}
