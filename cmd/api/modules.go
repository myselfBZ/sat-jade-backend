package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/myselfBZ/sat-jade/internal/store"
)

func (a *api) getModuleById(c echo.Context) error {
	id := c.Param("id")
	validId, err := strconv.Atoi(id)
	if err != nil {
		a.badRequestLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	module, err := a.storage.Modules.GetByID(c.Request().Context(), int32(validId))

	if err != nil {
		switch err {
		case store.ErrRecordNotFound:
			a.notFoundLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusNotFound, "module not found")
		default:
			a.internalErrLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	if err := a.constructModule(c.Request().Context(), module); err != nil {
		switch err {
		case store.ErrRecordNotFound:
			a.notFoundLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusNotFound, "module not found")
		default:
			a.internalErrLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	return c.JSON(http.StatusOK, module)
}

func (a *api) constructModule(ctx context.Context, m *store.Module) error {
	questions, err := a.storage.Questions.GetByModuleID(ctx, m.ID)
	if err != nil {
		return err
	}

	for _, q := range questions {
		answerChoices, err := a.storage.AnswerChoiceStorage.GetByQuestionID(ctx, q.ID)
		if err != nil {
			return err
		}
		q.AnswerChoices = append(q.AnswerChoices, answerChoices...)
	}

	m.Questions = append(m.Questions, questions...)
	return nil
}
