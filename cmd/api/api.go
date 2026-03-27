package main

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/myselfBZ/sat-jade/internal/auth"
	"github.com/myselfBZ/sat-jade/internal/llm"
	"github.com/myselfBZ/sat-jade/internal/store"
	"go.uber.org/zap"
)

type authConfig struct {
	secret string
	aud    string
	exp    time.Duration
}

type config struct {
	addr string
	auth authConfig
	frontEndUrl string
}

type api struct {
	config config
	logger *zap.SugaredLogger
	llm    llm.LLM
	// New
	auth    auth.Authenticator
	storage *store.Storage
}

func (a *api) registerRoutes() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{a.config.frontEndUrl, "http://localhost:5174"},
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
	results := v1.Group("/results", a.AuthMiddleware)
	auth := v1.Group("/auth")

	usersRouter.GET("/self", a.getUserSelfHandler)

	usersRouter.GET("/", a.getUsersHandler, a.CheckAdminMiddleware)
	usersRouter.DELETE("/:id", a.deleteUserHandler, a.CheckAdminMiddleware)

	usersResults := usersRouter.Group("/results")
	usersResults.POST("/", a.createResultHandler)
	usersResults.GET("/", a.getUserResultsHandler)
	usersResults.GET("/:userId", a.getAllResultsByUserHandler)

	results.GET("/", a.getAllResultsHandler, a.CheckAdminMiddleware)
	results.GET("/:id", a.getResultByIDHandler, a.CheckResultOwnership)
	results.DELETE("/:id", a.deleteResultByIDHandler)
	// results.POST("/:id/feedback", a.getOrCreateAIFeedbackHandler)


	practices.POST("/", a.createPracticeHandler, a.CheckAdminMiddleware)
	practices.DELETE("/:id", a.deletePracticeHandler, a.CheckAdminMiddleware)
	practices.GET("/:id", a.getPracticeByIDHandler)
	practices.GET("/", a.getPracticePreviewsHandler)

	questions := practices.Group("/questions")
	// update method is a must
	questions.POST("/", a.createQuestionHandler, a.CheckAdminMiddleware)
	// needs to be updated
	auth.POST("/token", a.createTokenHandler)
	auth.POST("/users", a.createUserHandler)

	questionBank := v1.Group("/question-bank", a.AuthMiddleware)
	questionBank.POST("/", a.createSqbQuestion, a.CheckAdminMiddleware)
	questionBank.GET("/", a.getCollectionDetails)
	questionBank.GET("/:id", a.getQuestionByID)
	questionBank.GET("/ids/:skill", a.getQuestionIdsBySkil)
	questionBank.POST("/answer", a.createQBAnswerHandler)
	questionBank.GET("/myanswers", a.getAnswersByUserHandler)

	return e
}

func (a *api) run() error {
	e := a.registerRoutes()
	return e.Start(a.config.addr)
}
