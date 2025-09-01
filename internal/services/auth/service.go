package auth

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	qusers "github.com/myselfBZ/sat-jade/internal/queries/users"
	"github.com/myselfBZ/sat-jade/internal/services/users"
)

func New(pool *pgxpool.Pool, secret, aud string, expr time.Duration) *AuthService {
	auth := Authenticator{
		secret: secret,
		aud:    aud,
		expr:   expr,
	}

	queries := qusers.New(pool)

	return &AuthService{
		users:         users.NewPgStore(queries),
		Authenticator: auth,
	}
}

type AuthService struct {
	users         users.Storage
	Authenticator Authenticator
}
