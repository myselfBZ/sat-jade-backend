package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/myselfBZ/sat-jade/internal/services/auth"
	"github.com/myselfBZ/sat-jade/internal/services/practice"
	"github.com/myselfBZ/sat-jade/internal/services/users"
)

type authConfig struct {
	secret string
}

type config struct {
	addr string
	auth authConfig
}

type api struct {
	users     *users.UserService
	auth      *auth.AuthService
	practices *practice.PracticeService
	config    config
}

func (a *api) registerRoutes() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			echo.HeaderAccessControlAllowOrigin,
		},
	}))
	v1 := e.Group("/v1")
	usersRouter := v1.Group("/users", a.AuthMiddleware)
	practices := v1.Group("/practices", a.AuthMiddleware)
	results := practices.Group("/results")
	questions := practices.Group("/questions")
	auth := v1.Group("/auth")

	usersRouter.POST("/tutor", a.users.CreateTutor)
	usersRouter.GET("/self", func(c echo.Context) error {
		user := c.Get("user").(*users.User)
		return c.JSON(http.StatusOK, user)
	})
	usersRouter.GET("/", a.users.GetMany)

	results.GET("/", a.practices.GetResults)
	results.GET("/all", a.practices.GetAllResults)
	results.GET("/:id", a.practices.GetSessionById)
	results.DELETE("/:id", a.practices.DeleteSession)
	results.POST("/:id/feedback", a.practices.GetSessionAIFeedback)

	practices.POST("/", a.practices.Create)
	practices.DELETE("/:id", a.practices.Delete)
	practices.GET("/:id", a.practices.GetById)
	practices.GET("/", a.practices.GetExamPreviews)
	practices.POST("/submit", a.practices.CreateTestSession)

	questions.POST("/", a.practices.AddQuestion)

	auth.POST("/token", a.auth.Login)
	auth.POST("/users", a.auth.SignUp)

	return e
}

func (a *api) run() error {
	e := a.registerRoutes()
	return e.Start(a.config.addr)
}
