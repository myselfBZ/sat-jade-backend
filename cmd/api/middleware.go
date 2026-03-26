package main

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/myselfBZ/sat-jade/internal/store"
)

const userCtxKey = "user"

func (app *api) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "authorization header is missing")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "authorization header is malformed")
		}

		token := parts[1]
		jwtToken, err := app.auth.ValidateToken(token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid claims")
		}

		userID, ok := claims["sub"].(string)

		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid subject claim")
		}

		user, err := app.storage.Users.GetByID(c.Request().Context(), userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
		}

		c.Set(userCtxKey, user)

		return next(c)
	}
}

func (a *api) CheckAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := a.getUserFromContext(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		if user.Role != store.ROLE_ADMIN {
			return echo.NewHTTPError(http.StatusUnauthorized, "you cant perform this action")
		}
		return next(c)
	}
}
