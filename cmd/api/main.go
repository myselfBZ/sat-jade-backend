package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"

	"github.com/myselfBZ/sat-jade/internal/auth"
	"github.com/myselfBZ/sat-jade/internal/db"
	"github.com/myselfBZ/sat-jade/internal/llm"
	"github.com/myselfBZ/sat-jade/internal/store"
)

var GEMINI_API_KEY string
var SERVER_ADDR string
var DB string
var SECRET_KEY string
var TOKEN_EXPR_HOURS int
var FRONT_END string

func loadEnvVars() {
	GEMINI_API_KEY = os.Getenv("GEMINI_API_KEY")
	if GEMINI_API_KEY == "" {
		panic("no GEMINI_API_KEY")
	}

	SERVER_ADDR = os.Getenv("SERVER_ADDR")
	if SERVER_ADDR == "" {
		panic("no SERVER_ADDR")
	}

	DB = os.Getenv("DB")
	if DB == "" {
		panic("no DB")
	}

	SECRET_KEY = os.Getenv("SECRET_KEY")
	if SECRET_KEY == "" {
		panic("no SECRET_KEY")
	}
	tokenExpr := os.Getenv("TOKEN_EXPR_HOURS")

	var err error

	TOKEN_EXPR_HOURS, err = strconv.Atoi(tokenExpr)
	if err != nil {
		fmt.Printf("expected int for TOKEN_EXPR_HOURS GOT %s\n", tokenExpr)
	}

	FRONT_END = os.Getenv("FRONTEND_URL")

	if FRONT_END == "" {
		panic("no FRONTEND_URL")
	}
}

func main() {
	loadEnvVars()

	api := api{
		config: config{
			addr: ":" + SERVER_ADDR,
			frontEndUrl: FRONT_END,
			auth: authConfig{
				secret: SECRET_KEY,
				exp:    time.Hour * time.Duration(TOKEN_EXPR_HOURS),
				aud:    "test-aud",
			},
		},
	}
	logger := zap.Must(zap.NewProduction(zap.AddCaller())).Sugar()
	defer logger.Sync()
	api.logger = logger

	db, err := db.New(db.Config{
		Addr:        DB,
		MaxConns:    15,
		MinConns:    15,
		MaxIdleTime: "15m",
	})

	if err != nil {
		panic(err)
	}
	api.storage = store.New(db)
	api.auth = auth.NewJWTAuthenticator(
		api.config.auth.secret,
		api.config.auth.aud,
		api.config.auth.aud,
	)
	api.llm, err = llm.NewGemini(GEMINI_API_KEY)

	if err != nil {
		fmt.Println("error connecting to gemini")
		panic(err)
	}

	api.run()
}
