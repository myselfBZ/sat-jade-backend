package main

import (
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/myselfBZ/sat-jade/internal/db"
	"github.com/myselfBZ/sat-jade/internal/services/auth"
	"github.com/myselfBZ/sat-jade/internal/services/practice"
	"github.com/myselfBZ/sat-jade/internal/services/users"
)

var GEMINI_API_KEY string
var SERVER_ADDR string
var DB string
var SECRET_KEY string

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
}

func main() {
	loadEnvVars()

	api := api{
		config: config{
			addr: ":" + SERVER_ADDR,
			auth: authConfig{
				secret: SECRET_KEY,
			},
		},
	}

	db, err := db.New(db.Config{
		Addr:        DB,
		MaxConns:    15,
		MinConns:    15,
		MaxIdleTime: "15m",
	})

	if err != nil {
		panic(err)
	}

	userService := users.New(db)
	api.users = userService
	api.auth = auth.New(db, api.config.auth.secret, api.config.auth.secret, (time.Hour * 24))
	api.practices = practice.New(db, GEMINI_API_KEY)
	api.run()
}
