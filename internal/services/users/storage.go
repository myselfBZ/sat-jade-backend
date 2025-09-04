package users

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgtype"
	qusers "github.com/myselfBZ/sat-jade/internal/queries/users"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Storage interface {
	//GetByID(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, u *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
}

func NewPgStore(queries *qusers.Queries) *PostgresStorage {
	return &PostgresStorage{
		queries: *queries,
	}
}

type PostgresStorage struct {
	queries qusers.Queries
}

func (s *PostgresStorage) GetByID(ctx context.Context, id string) (*User, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	user, err := s.queries.GetUserById(ctx, pgtype.UUID{Bytes: parsedID, Valid: true})

	if err != nil {
		return nil, err
	}
	return &User{
		ID:        user.ID.String(),
		Email:     user.Email,
		Password:  user.PasswordHash,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}, nil
}

func (s *PostgresStorage) GetByEmail(ctx context.Context, email string) (*User, error) {
	user, err := s.queries.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:        user.ID.String(),
		Email:     user.Email,
		Password:  user.PasswordHash,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}, nil
}

func (s *PostgresStorage) Create(ctx context.Context, u *User) error {
	user, err := s.queries.Create(ctx, qusers.CreateParams{
		FullName:     u.FullName,
		Email:        u.Email,
		PasswordHash: u.Password,
		Role:         u.Role,
	})

	if err != nil {
		return err
	}

	u.ID = user.ID.String()
	return nil
}
