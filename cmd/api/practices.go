package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/myselfBZ/sat-jade/internal/store"
)

type practiceRequestPayload struct {
	Title string `json:"title"`
}

type practiceResponsePayload struct {
	ID      int32           `json:"id"`
	Title   string          `json:"title"`
	Modules []*store.Module `json:"sections"`
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

	practice, err := a.storage.Practices.GetById(c.Request().Context(), int32(validId))

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

	if err := a.constructPracticeTest(c.Request().Context(), practice); err != nil {
		switch err {
		case store.ErrRecordNotFound:
			a.notFoundLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusNotFound, "practice test not found")
		default:
			a.internalErrLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	return c.JSON(http.StatusOK, &practiceResponsePayload{
		ID:      practice.ID,
		Title:   practice.Title,
		Modules: practice.Modules,
	})
}

func (a *api) getPracticePreviewsHandler(c echo.Context) error {
	practices, err := a.storage.Practices.GetAllPreview(c.Request().Context())
	if err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, practices)
}

// helper
func (a *api) constructPracticeTest(ctx context.Context, p *store.Practice) error {
	modules, err := a.storage.Modules.GetAllByPracticeID(ctx, p.ID)

	if err != nil {
		return err
	}

	for _, m := range modules {
		if err := a.constructModule(ctx, m); err != nil {
			return err
		}
	}

	p.Modules = modules

	return nil
}
