package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/myselfBZ/sat-jade/internal/store"
)

type practiceRequestPayload struct {
	Title string `json:"title"`
}

func (a *api) createPracticeHandler(c echo.Context) error {
	var payload practiceRequestPayload
	if err := c.Bind(&payload); err != nil {
		a.badRequestLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	id, err := a.storage.Practices.Create(c.Request().Context(), payload.Title)

	if err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"id": id,
	})
}

func (a *api) deletePracticeHandler(c echo.Context) error {
	id := c.Param("id")
	validId, err := strconv.Atoi(id)
	if err != nil {
		a.badRequestLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	if err := a.storage.Practices.Delete(c.Request().Context(), int32(validId)); err != nil {
		switch err {
		case store.ErrRecordNotFound:
			a.notFoundLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusNotFound)
		default:
			a.internalErrLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	return c.JSON(http.StatusCreated, map[string]string{
		"Status": "Ok",
	})
}

func (a *api) getPracticeByIDHandler(c echo.Context) error {
	id := c.Param("id")
	validId, err := strconv.Atoi(id)

	if err != nil {
		a.badRequestLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	practice, err := a.storage.Practices.GetFullTest(c.Request().Context(), int32(validId))

	if err != nil {
		switch err {
		case store.ErrRecordNotFound:
			a.notFoundLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusNotFound, "practice test not found")
		default:
			a.internalErrLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	return c.JSON(http.StatusOK, practice)
}

func (a *api) getPracticePreviewsHandler(c echo.Context) error {
	practices, err := a.storage.Practices.GetAllPreview(c.Request().Context())
	if err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, practices)
}

