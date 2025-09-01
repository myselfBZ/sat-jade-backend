package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/myselfBZ/sat-jade/internal/services/errs"
	"github.com/myselfBZ/sat-jade/internal/services/users"
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

func (s *AuthService) Login(c echo.Context) error {
	var payload loginPayload

	if err := c.Bind(&payload); err != nil {
		return err
	}

	user, err := s.users.GetByEmail(c.Request().Context(), payload.Email)

	if err != nil {
		switch err {
		case errs.ErrRecordNotFound:
			echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, errs.ErrUnauthorized.Error())
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(s.Authenticator.expr).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": s.Authenticator.aud,
		"aud": s.Authenticator.aud,
	}

	token, err := s.Authenticator.GenerateToken(claims)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, &tokenEnvelope{Token: token})
}

func (s *AuthService) SignUp(c echo.Context) error {
	var payload userPayload
	if err := c.Bind(&payload); err != nil {
		return err
	}

	user := users.User{
		Email:    payload.Email,
		FullName: payload.FullName,
		Role:     users.ROLE_STUDENT,
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	user.Password = string(hash)

	if err := s.users.Create(c.Request().Context(), &user); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(s.Authenticator.expr).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": s.Authenticator.aud,
		"aud": s.Authenticator.aud,
	}

	token, err := s.Authenticator.GenerateToken(claims)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, &tokenEnvelope{Token: token})

}
