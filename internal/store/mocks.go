package store

import (
	"context"
	"strings"
	"sync"

	"github.com/google/uuid"
)


func NewMockStorage() *Storage {
	return &Storage{
		Users: NewUserMockStore(),
		Practices: NewPracticeMockStore(),
		Questions: NewQuestionsMockStore(),
	}
}

type UserMockStore struct {
	mu    sync.RWMutex
	users map[string]User
}

func NewUserMockStore() *UserMockStore {
	return &UserMockStore{
		users: make(map[string]User),
	}
}

func (m *UserMockStore) Create(ctx context.Context, u *User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, existing := range m.users {
		if existing.Email == u.Email {
			return  ErrDuplicateEmail
		}
	}

	u.ID = uuid.NewString()

	m.users[u.ID] = *u
	return nil
}

func (m *UserMockStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, u := range m.users {
		if u.Email == email {
			userCopy := u
			return &userCopy, nil
		}
	}
	return nil, ErrRecordNotFound
}

func (m *UserMockStore) GetByID(ctx context.Context, id string) (*User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	u, exists := m.users[id]
	if !exists {
		return nil, ErrRecordNotFound
	}
	
	userCopy := u
	return &userCopy, nil
}

func (m *UserMockStore) GetMany(ctx context.Context) ([]User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	users := make([]User, 0, len(m.users))
	for _, u := range m.users {
		users = append(users, u)
	}
	return users, nil
}

func (m *UserMockStore) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.users[id]; !exists {
		return ErrRecordNotFound
	}
	
	delete(m.users, id)
	return nil
}


type PracticeMockStore struct {
	mu        sync.RWMutex
	lastID    int32
	practices map[int32]Practice
}

func NewPracticeMockStore() *PracticeMockStore {
	return &PracticeMockStore{
		practices: make(map[int32]Practice),
	}
}


func (m *PracticeMockStore) seedPractice(practice *Practice) {
	rw1 := Module{
		Name: "Reading And Writing 1",
		Questions: make([]*Question, 27),
	}

	rw2 := Module{
		Name: "Reading And Writing 2",
		Questions: make([]*Question, 27),
	}

	math1 := Module{
		Name: "Math 1",
		Questions: make([]*Question, 22),
	}

	math2 := Module{
		Name: "Math 2",
		Questions: make([]*Question, 22),
	}

	practice.Modules = append(practice.Modules, rw1, rw2, math1, math2)
}

func (m *PracticeMockStore) Create(ctx context.Context, title string) (int32, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.lastID++
	newID := m.lastID

	p := Practice{
		ID:    newID,
		Title: title,
	}

	m.seedPractice(&p)


	m.practices[newID] = p
	return newID, nil
}

func (m *PracticeMockStore) GetFullTest(ctx context.Context, id int32) (*Practice, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, exists := m.practices[id]
	if !exists {
		return nil, ErrRecordNotFound
	}

	practiceCopy := p
	return &practiceCopy, nil
}

func (m *PracticeMockStore) GetAllPreview(ctx context.Context) ([]PracticePreview, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	previews := make([]PracticePreview, 0, len(m.practices))
	for _, p := range m.practices {
		previews = append(previews, PracticePreview{
			ID:    p.ID,
			Title: p.Title,
			CreatedAt: p.CreatedAt,
		})
	}
	return previews, nil
}

func (m *PracticeMockStore) GetCorrectAnswers(ctx context.Context, id int32) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.practices[id]
	if !exists {
		return nil, ErrRecordNotFound
	}

	// 98 As
	return strings.Split(strings.Repeat("A ", 98), " "), nil
}

func (m *PracticeMockStore) Delete(ctx context.Context, id int32) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.practices[id]; !exists {
		return ErrRecordNotFound
	}

	delete(m.practices, id)
	return nil
}




type QuestionsMockStore struct{}

func NewQuestionsMockStore() *QuestionsMockStore {
	return &QuestionsMockStore{}
}

// Returns ErrForeignConstraint on Question.PracticeId == 0 or moduleID === 0
// any other id resolves in success
func (m *QuestionsMockStore) CreateWithAnswerChoices(ctx context.Context, moduleID int32, q *Question) error {
	if moduleID == 0 || q.PracticeId == 0{
		return ErrForeignConstraint
	}

	return nil
}
