package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/myselfBZ/sat-jade/internal/store"
	"github.com/stretchr/testify/assert"
)


func TestAuthMiddleware(t *testing.T) {
	e := echo.New()
	api := newTestApi()

	t.Run("401 unauthorized", func(t *testing.T) {
		e.GET("/", func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}, api.AuthMiddleware)


		req := newTestRequest(t, http.MethodGet, "/", nil)

		rr := executeRequest(req, e)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})


	t.Run("200 success with authentication", func(t *testing.T) {
		e.GET("/", func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}, api.AuthMiddleware)

		token := newTestJWTToken(t, api, adminUser.ID)

		req := newTestRequest(t, http.MethodGet, "/", nil)

		req.Header.Add("Authorization", "Bearer "+token)

		rr := executeRequest(req, e)

		assert.Equal(t, http.StatusOK, rr.Code)

	})


	t.Run("401 unauthorized with invalid subject claim", func(t *testing.T) {
		claims := jwt.MapClaims{
			"exp": time.Now().Add(time.Hour).Unix(),
			"iat": time.Now().Unix(),
			// invalid user id
			"sub":123,
		}

		token, err := api.auth.GenerateToken(claims)
		if err != nil {
			t.Fatalf("couldn't generate a jwt token: %v", err)
		}

		req := newTestRequest(t, http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		rr := executeRequest(req, e)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})


	t.Run("401 unauthorized with malformed token", func(t *testing.T) {
		req := newTestRequest(t, http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer this.isnot.atoken")

		rr := executeRequest(req, e)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})


	t.Run("401 unauthorized with expired token", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": adminUser.ID,
			"exp": time.Now().Add(-time.Hour).Unix(), 
		}

		token, err := api.auth.GenerateToken(claims)

		if err != nil {
			t.Fatalf("couldn't generate a jwt token: %v", err)
		}

		req := newTestRequest(t, http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		rr := executeRequest(req, e)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}






func TestCheckAdminMiddleware(t *testing.T) {
	e := echo.New()

	api := newTestApi()



	t.Run("200 is admin", func(t *testing.T) {
		middleware := api.CheckAdminMiddleware(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		})

		req := newTestRequest(t, http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		c := e.NewContext(req, rr)

		c.Set(userCtxKey, &adminUser)

		err := middleware(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("401 not admin regular user", func(t *testing.T) {
		middleware := api.CheckAdminMiddleware(func(c echo.Context) error {
			return c.NoContent(http.StatusNoContent)
		})

		req := newTestRequest(t, http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		c := e.NewContext(req, rr)

		c.Set(userCtxKey, &store.User{
			Role: store.ROLE_STUDENT,
		})

		err := middleware(c)
		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, he.Code)
	})

}



