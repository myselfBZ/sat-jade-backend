package store

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/myselfBZ/sat-jade/internal/queries/modules"
)

type Module struct {
	PracticeID int32       `json:"practice_id"`
	ID         int32       `json:"id"`
	Name       string      `json:"name"`
	Questions  []*Question `json:"questions"`
}

type ModuleStore struct {
	queries *modules.Queries
}

func NewModuleStore(db *pgxpool.Pool) *ModuleStore {
	queries := modules.New(db)
	return &ModuleStore{
		queries: queries,
	}
}

func (s *ModuleStore) GetByNameAndPracticeID(ctx context.Context, name string, practiceID int32) (*Module, error) {
	moduleRow, err := s.queries.GetByNameAndPracticeID(ctx, modules.GetByNameAndPracticeIDParams{
		PracticeID: practiceID,
		Name:       name,
	})
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &Module{
		ID:         moduleRow.ID,
		PracticeID: moduleRow.PracticeID,
		Name:       moduleRow.Name,
		Questions:  []*Question{},
	}, nil
}

func (s *ModuleStore) GetAllByPracticeID(ctx context.Context, practiceID int32) ([]*Module, error) {
	dbModules, err := s.queries.GetByPracticeId(ctx, practiceID)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	var modules []*Module
	for _, m := range dbModules {
		modules = append(modules, &Module{
			PracticeID: m.PracticeID,
			ID:         m.ID,
			Name:       m.Name,
			Questions:  []*Question{},
		})
	}
	return modules, nil
}

func (s *ModuleStore) GetByID(ctx context.Context, id int32) (*Module, error) {
	dbModule, err := s.queries.GetByID(ctx, id)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &Module{
		PracticeID: dbModule.PracticeID,
		ID:         dbModule.ID,
		Name:       dbModule.Name,
		Questions:  []*Question{},
	}, nil
}
