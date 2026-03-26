package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/myselfBZ/sat-jade/internal/store"
)

func (a *api) getUserFromContext(c echo.Context) (*store.User, error) {
	// validation is crucial!
	user, ok := c.Get("user").(*store.User)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	return user, nil
}

func (a *api) getUserSelfHandler(c echo.Context) error {
	user, err := a.getUserFromContext(c)

	if err != nil {
		a.unauthorizedLog(c.Request().Method, c.Path(), err)
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (a *api) getUsersHandler(c echo.Context) error {
	users, err := a.storage.Users.GetMany(c.Request().Context())
	if err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, users)
}

func (a *api) deleteUserHandler(c echo.Context) error {
	userId := c.Param("id")
	if err := a.storage.Users.Delete(c.Request().Context(), userId); err != nil {
		switch err {
		case store.ErrRecordNotFound:
			a.notFoundLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		default:
			a.badRequestLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
		}
	}

	return c.JSON(http.StatusNoContent, nil)
}
