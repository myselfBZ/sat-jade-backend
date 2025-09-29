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

func (a *api) getUsersHandler(c echo.Context) error {
	users, err := a.storage.Users.GetMany(c.Request().Context())
	if err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, users)
}
