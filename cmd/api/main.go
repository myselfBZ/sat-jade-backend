package main

import (
	"time"

	"github.com/myselfBZ/sat-jade/internal/db"
	"github.com/myselfBZ/sat-jade/internal/services/auth"
	"github.com/myselfBZ/sat-jade/internal/services/practice"
	"github.com/myselfBZ/sat-jade/internal/services/users"
)

func main() {
	api := api{
		config: config{
			addr: ":8080",
			auth: authConfig{
				secret: "secretKey",
			},
		},
	}

	db, err := db.New(db.Config{
		Addr:        "host=localhost port=32768 user=postgres password=new_password dbname=satjade sslmode=disable",
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
	api.practices = practice.New(db)
	api.run()
}
