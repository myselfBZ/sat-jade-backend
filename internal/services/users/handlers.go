package users

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type teacherPayload struct {
	FullName string `json:"full_name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (s *UserService) GetByEmail(c echo.Context) error {
	return nil
}

func (s *UserService) CreateTutor(c echo.Context) error {
	/*currentUser := c.Get("user").(*User)
	if currentUser.Role != "admin" {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}*/

	var payload teacherPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	user := User{
		Email:    payload.Email,
		FullName: payload.FullName,
		Role:     ROLE_TUTOR,
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	user.Password = string(hash)

	if err := s.storage.Create(c.Request().Context(), &user); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, map[string]string{
		"Success": "OK",
	})
}

func (s *UserService) GetById(c echo.Context) error {

	return nil
}

func (s *UserService) GetMany(c echo.Context) error {
	user := c.Get("user").(*User)
	if user.Role != ROLE_ADMIN {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	users, err := s.storage.GetMany(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, users)
}
