package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/myselfBZ/sat-jade/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userPayload struct {
	FullName string `json:"full_name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type tokenEnvelope struct {
	Token string `json:"token"`
}

func (a *api) createTokenHandler(c echo.Context) error {
	var payload loginPayload

	if err := c.Bind(&payload); err != nil {
		return err
	}

	if payload.Email == "" || payload.Password == "" {
		a.badRequestLog(c.Request().Method, c.Path(), errors.New("invalid payload"))
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	user, err := a.storage.Users.GetByEmail(c.Request().Context(), payload.Email)

	if err != nil {
		switch err {
		case store.ErrRecordNotFound:
			a.notFoundLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			a.internalErrLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(a.config.auth.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": a.config.auth.aud,
		"aud": a.config.auth.aud,
	}

	token, err := a.auth.GenerateToken(claims)

	if err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, &tokenEnvelope{Token: token})
}

func (a *api) createUserHandler(c echo.Context) error {
	var payload userPayload
	if err := c.Bind(&payload); err != nil {
		return err
	}

	user := store.User{
		Email:    payload.Email,
		FullName: payload.FullName,
		Role:     store.ROLE_STUDENT,
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	if err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	user.Password = string(hash)

	if err := a.storage.Users.Create(c.Request().Context(), &user); err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			a.conflictLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusConflict, err)
		default:
			a.internalErrLog(c.Request().Method, c.Path(), err)
			return echo.NewHTTPError(http.StatusInternalServerError)

		}
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(a.config.auth.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": "test-aud",
		"aud": "test-aud",
	}

	token, err := a.auth.GenerateToken(claims)

	if err != nil {
		a.internalErrLog(c.Request().Method, c.Path(), err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, &tokenEnvelope{Token: token})

}
