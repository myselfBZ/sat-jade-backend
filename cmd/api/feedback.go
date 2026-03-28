package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/myselfBZ/sat-jade/internal/services/feedback"
	"github.com/myselfBZ/sat-jade/internal/store"
)


func (a *api) generateFeedbakcHandler(c echo.Context) error {
	result := c.Get(resultCtxKey).(*store.Result)

	validUserId, err := uuid.Parse(result.UserId)

	if err != nil {
		a.badRequestLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
	}

	resultOverview, err :=  a.storage.Results.GetOverview(c.Request().Context(), result.ID)

	if err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	fb, err := a.feedbackService.Generate(c.Request().Context(), &feedback.GenerateParams{
		Overview: resultOverview,
		UserId:   validUserId,
		ResultId: result.ID,
	})

	if err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, fb)
}


