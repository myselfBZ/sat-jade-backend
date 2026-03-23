package store

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	qusers "github.com/myselfBZ/sat-jade/internal/queries/users"
)

const (
	ROLE_STUDENT = "student"
	ROLE_ADMIN   = "admin"
	ROLE_TUTOR   = "tutor"
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

func NewUserStore(db *pgxpool.Pool) *UserStore {
	queries := qusers.New(db)
	return &UserStore{
		queries: queries,
	}
}

type UserStore struct {
	queries *qusers.Queries
}

func (s *UserStore) GetByID(ctx context.Context, id string) (*User, error) {
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
		FullName:  user.FullName,
		Email:     user.Email,
		Password:  user.PasswordHash,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}, nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	user, err := s.queries.GetByEmail(ctx, email)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
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

func (s *UserStore) Create(ctx context.Context, u *User) error {
	user, err := s.queries.Create(ctx, qusers.CreateParams{
		FullName:     u.FullName,
		Email:        u.Email,
		PasswordHash: u.Password,
		Role:         u.Role,
	})

	if err != nil {
		switch {
		case err.Error() == `ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	u.ID = user.ID.String()
	return nil
}

func (s *UserStore) GetMany(ctx context.Context) ([]User, error) {
	userRows, err := s.queries.GetMany(ctx)
	var users []User
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return users, nil
		}
		return nil, err
	}

	for _, u := range userRows {
		users = append(users, User{
			ID:        u.ID.String(),
			Email:     u.Email,
			FullName:  u.FullName,
			Role:      u.Role,
			CreatedAt: u.CreatedAt.Time,
			UpdatedAt: u.UpdatedAt.Time,
		})
	}

	return users, nil
}

func (s *UserStore) Delete(ctx context.Context, id string) error {
	validId, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidUUID
	}
	_, err = s.queries.Delete(ctx, pgtype.UUID{
		Bytes: validId,
		Valid: true,
	})

	return err
}




