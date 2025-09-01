package users

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	qusers "github.com/myselfBZ/sat-jade/internal/queries/users"
)

const (
	ROLE_STUDENT = "student"
	ROLE_ADMIN   = "admint"
	ROLE_TUTOR   = "tutor"
)

type UserService struct {
	storage Storage
}

func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
	return s.storage.GetByID(ctx, id)
}

func New(db *pgxpool.Pool) *UserService {
	queries := qusers.New(db)

	return &UserService{
		storage: &PostgresStorage{
			queries: *queries,
		},
	}
}
