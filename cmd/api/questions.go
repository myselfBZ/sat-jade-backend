package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/myselfBZ/sat-jade/internal/store"
)

type questionPayload struct {
	Number        int32                 `json:"number"`
	PracticeId    int32                 `json:"practice_id"`
	Difficulty    string                `json:"difficulty"`
	Module        string                `json:"section"`
	Domain        string                `json:"domain"`
	Correct       string                `json:"correct"`
	Paragraph     string                `json:"paragraph"`
	Prompt        string                `json:"prompt"`
	Explanation   string                `json:"explanation"`
	AnswerChoices []*store.AnswerChoice `json:"choices"`
}

func (a *api) createQuestionHandler(c echo.Context) error {
	var p questionPayload

	if err := c.Bind(&p); err != nil {
		a.badRequestLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	question := &store.Question{
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

	module, err := a.storage.Modules.GetByNameAndPracticeID(c.Request().Context(),
		question.Module,
		question.PracticeId,
	)

	if err != nil {
		switch err {
		case store.ErrRecordNotFound:
			a.notFoundLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusNotFound, "Module for this question is not found")
		default:
			a.internalErrLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	if err := a.storage.Questions.CreateWithAnswerChoices(c.Request().Context(), module.ID, question); err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"Status": "Ok",
	})

}
