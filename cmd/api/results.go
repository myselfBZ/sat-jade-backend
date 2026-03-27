package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/myselfBZ/sat-jade/internal/grading"
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

	correctAnswersWithChoices, err := a.storage.Practices.GetCorrectAnswersWithAnswerChoices(c.Request().Context(), p.ExamID)
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

	result := grading.Check(p.Answers, correctAnswersWithChoices)
	result.UserId = user.ID
	result.PracticeId = p.ExamID

	if err := a.storage.Results.Create(c.Request().Context(), result); err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if err := a.storage.ResultAnswers.CreateMany(c.Request().Context(), result.ID, result.Answers); err != nil {
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
	result, ok := c.Get(resultCtxKey).(*store.Result)

	if !ok {
		a.internalErrLog(c.Request().Method, c.Path(), errors.New("result not set in the context"))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	answers, err := a.storage.ResultAnswers.GetByResultID(c.Request().Context(), result.ID)
	if err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
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
