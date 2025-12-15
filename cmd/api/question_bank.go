package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/myselfBZ/sat-jade/internal/store"
)

func (a *api) createSqbQuestion(c echo.Context) error {
	var sqbQuestion store.SQBQuestion
	if err := c.Bind(&sqbQuestion); err != nil {
		a.badRequestLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := a.storage.QuestionBank.Create(c.Request().Context(), &sqbQuestion); err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"Status": "Ok",
	})
}

func (a *api) getQuestionIdsBySkil(c echo.Context) error {
	skill := c.Param("skill")
	ids, err := a.storage.QuestionBank.GetIdBySkill(c.Request().Context(), skill)
	if err != nil {
		a.notFoundLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusNotFound, "skill not found")
	}
	return c.JSON(http.StatusOK, ids)
}

func (a *api) getQuestionByID(c echo.Context) error {
	user := c.Get("user").(*store.User)
	id := c.Param("id")
	validInt, err := strconv.Atoi(id)

	if err != nil {
		a.badRequestLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	q, err := a.storage.QuestionBank.GetById(c.Request().Context() ,validInt, user.ID)

	if err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, q)
}

func (a *api) getCollectionDetails(c echo.Context) error {
	collection, err := a.storage.QuestionBank.GetCollectionDetail(c.Request().Context())
	if err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, collection)
}
